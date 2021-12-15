package xegrpcgw

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/gin-contrib/cors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	inruntime "github.com/vicnoah/ego-component/xegrpcgw/runtime"
)

// corsIntercepter 跨域注入
func corsIntercepter(c *Container) {
	corsConf := cors.Config{
		AllowOrigins:     c.config.AccessControlAllowOrigin,
		AllowMethods:     c.config.AccessControlAllowMethods,
		AllowHeaders:     c.config.AccessControlAllowHeaders,
		ExposeHeaders:    c.config.AccessControlExposeHeaders,
		AllowCredentials: c.config.AccessControlAllowCredentials,
		MaxAge:           12 * time.Hour,
	}
	if c.eginCorsOriginFunc != nil {
		corsConf.AllowOriginFunc = c.eginCorsOriginFunc
	}
	c.eg.Use(cors.New(corsConf))
}

// customerEcodeOption 自定义错误码
func customerEcodeOption(c *Container) {
	c.muxOptions = append(c.muxOptions, runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		inruntime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
	}))
	c.muxOptions = append(c.muxOptions, runtime.WithStreamErrorHandler(func(ctx context.Context, err error) *status.Status {
		return status.Convert(err)
	}))
}

// incomintHeaderMatcherOption 传输自定义http头到grpc server
func incomingHeaderMatcherOption(c *Container) {
	c.muxOptions = append(c.muxOptions, runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		// 匹配http请求头到grpc
		for _, header := range c.config.IncomingHeaders {
			if key == header {
				return key, true
			}
		}
		return runtime.DefaultHeaderMatcher(key)
	}))
}

// urlPathTransOption 传输url到grpc server
func urlPathTransOption(c *Container) {
	c.muxOptions = append(c.muxOptions, runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		md := make(metadata.MD)
		md[urlMetadataName] = []string{req.URL.Path}
		return md
	}))
}

// httpResponseModifier 自定义响应头
func httpResponseModifier(c *Container) {
	c.muxOptions = append(c.muxOptions, runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
		md, ok := runtime.ServerMetadataFromContext(ctx)
		if !ok {
			return nil
		}
		// set http status code
		if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
			code, err := strconv.Atoi(vals[0])
			if err != nil {
				return err
			}
			// delete the headers to not expose any grpc-metadata in http response
			delete(md.HeaderMD, "x-http-code")
			delete(w.Header(), "Grpc-Metadata-X-Http-Code")
			w.WriteHeader(code)
		}

		return nil
	}))
}
