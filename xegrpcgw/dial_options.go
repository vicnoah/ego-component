package xegrpcgw

import (
	"net/http"
	"time"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/etrace"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	"github.com/uber/jaeger-client-go"
)

// handler 中间件
type handler func(http.Handler) http.Handler

// withTracer tracing服务
func withTracer(c *Container) grpc.DialOption {
	return grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor(
		grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
	))
}

// metricServerInterceptor 度量服务
func metricServerInterceptor(c *Container) {
	c.muxWrappers = append(c.muxWrappers, metricWrapper)
}

func metricWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beg := time.Now()
		emetric.ServerHandleHistogram.Observe(float64(time.Since(beg).Seconds()), emetric.TypeHTTP, r.Method+"."+r.URL.Path, r.Header.Get("app"))
		emetric.ServerHandleCounter.Inc(emetric.TypeHTTP, r.Method+"."+r.URL.Path, r.Header.Get("app"), http.StatusText(r.Response.StatusCode))
	})
}

// traceServerIntercepter 链路追踪服务
func traceServerIntercepter(c *Container) {
	c.muxWrappers = append(c.muxWrappers, tracingWrapper)
}

func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		// 判断了全局jaeger的设置，所以这里一定能够断言为jaeger
		r.Header.Set(eapp.EgoTraceIDName(), span.(*jaeger.Span).Context().(jaeger.SpanContext).TraceID().String())
	})
}

// func tracingWrapper(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		parentSpanContext, err := opentracing.GlobalTracer().Extract(
// 			opentracing.HTTPHeaders,
// 			opentracing.HTTPHeadersCarrier(r.Header))
// 		if err == nil || err == opentracing.ErrSpanContextNotFound {
// 			serverSpan := opentracing.GlobalTracer().StartSpan(
// 				"ServeHTTP",
// 				// this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
// 				ext.RPCServerOption(parentSpanContext),
// 				grpcGatewayTag,
// 			)
// 			r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))
// 			defer serverSpan.Finish()
// 		}
// 		h.ServeHTTP(w, r)
// 	})
// }
