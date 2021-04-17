package xegrpcgw

import (
	"net/http"
	"time"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/etrace"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	"github.com/uber/jaeger-client-go"
)

// withTracer tracing服务
func withTracer(c *Container) grpc.DialOption {
	return grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor(
		grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
	))
}

// handlerInterceptor 中间件注入
func handlerInterceptor(c *Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range c.handlerFuncs {
			handler(w, r)
		}
		c.mux.ServeHTTP(w, r)
	})
}

// corsIntercepter 跨域注入
func corsIntercepter(c *Container) {
	c.handlerFuncs = append(c.handlerFuncs, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", c.config.AccessControlAllowOrigin)           // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", c.config.AccessControlAllowHeaders)         //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", c.config.AccessControlAllowCredentials) //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", c.config.AccessControlAllowMethods)         //允许请求方法
		w.Header().Set("content-type", c.config.ContentType)                                       //返回数据格式是json
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}))
}

// metricServerInterceptor 度量服务
func metricServerInterceptor(c *Container) {
	c.handlerFuncs = append(c.handlerFuncs, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beg := time.Now()
		emetric.ServerHandleHistogram.Observe(float64(time.Since(beg).Seconds()), emetric.TypeHTTP, r.Method+"."+r.URL.Path, r.Header.Get("app"))
		emetric.ServerHandleCounter.Inc(emetric.TypeHTTP, r.Method+"."+r.URL.Path, r.Header.Get("app"), http.StatusText(200))
	}))
}

// traceServerIntercepter 链路追踪服务
func traceServerIntercepter(c *Container) {
	c.handlerFuncs = append(c.handlerFuncs, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 为了性能考虑，如果要加日志字段，需要改变slice大小
		var fields = make([]elog.Field, 0, 15)
		span, ctx := etrace.StartSpanFromContext(
			r.Context(),
			r.Method+"."+r.URL.Path,
			etrace.TagComponent("http"),
			etrace.TagSpanKind("server"),
			etrace.HeaderExtractor(r.Header),
			etrace.CustomTag("http.url", r.URL.Path),
			etrace.CustomTag("http.method", r.Method),
			etrace.CustomTag("peer.ipv4", r.RemoteAddr),
		)
		r = r.WithContext(ctx)
		defer span.Finish()
		fields = append(fields,
			elog.FieldType("http"),
			elog.FieldMethod(r.URL.Path),
			elog.FieldPeerIP(r.RemoteAddr),
		)
		if c.config.EnableTraceInterceptor && opentracing.IsGlobalTracerRegistered() {
			fields = append(fields, elog.FieldTid(etrace.ExtractTraceID(ctx)))
		}
		c.logger.Info("grpc-gateway", fields...)
		// 判断了全局jaeger的设置，所以这里一定能够断言为jaeger
		r.Header.Set(eapp.EgoTraceIDName(), span.(*jaeger.Span).Context().(jaeger.SpanContext).TraceID().String())
	}))
}
