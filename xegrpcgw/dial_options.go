package xegrpcgw

import (
	"fmt"
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

// metricServerInterceptor 度量服务
func metricServerInterceptor(c *Container) {
	c.handlerFuncs = append(c.handlerFuncs, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beg := time.Now()
		emetric.ServerHandleHistogram.Observe(float64(time.Since(beg).Seconds()), emetric.TypeHTTP, r.Method+"."+r.URL.Path, r.Header.Get("app"))
		emetric.ServerHandleCounter.Inc(emetric.TypeHTTP, r.Method+"."+r.URL.Path, r.Header.Get("app"), http.StatusText(r.Response.StatusCode))
	}))
}

// traceServerIntercepter 链路追踪服务
func traceServerIntercepter(c *Container) {
	c.handlerFuncs = append(c.handlerFuncs, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("开始度量")
		span, ctx := etrace.StartSpanFromContext(
			r.Context(),
			r.Method+"."+r.URL.Path,
			etrace.TagComponent("http"),
			etrace.TagSpanKind("server"),
			etrace.HeaderExtractor(r.Response.Header),
			etrace.CustomTag("http.url", r.URL.Path),
			etrace.CustomTag("http.method", r.Method),
			etrace.CustomTag("peer.ipv4", r.RemoteAddr),
		)
		r = r.WithContext(ctx)
		defer span.Finish()
		elog.Info("http", elog.FieldType("http"), elog.FieldMethod(r.URL.Path), elog.FieldPeerIP(r.RemoteAddr))
		// 判断了全局jaeger的设置，所以这里一定能够断言为jaeger
		r.Header.Set(eapp.EgoTraceIDName(), span.(*jaeger.Span).Context().(jaeger.SpanContext).TraceID().String())
	}))
}
