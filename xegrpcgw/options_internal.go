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
	"github.com/gotomicro/ego/core/eerrors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	inruntime "github.com/vicnoah/ego-component/xegrpcgw/runtime"
)

const (
	urlMetadataName = "url" // url信息
)

// withCorsIntercepter 跨域注入
func withCorsIntercepter(c *Container) {
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

// withCustomerEcodeOption 自定义错误码
func withCustomerEcodeOption(c *Container) {
	// grpc ecode转http error
	c.muxOptions = append(c.muxOptions, runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		inruntime.DefaultHTTPErrorHandler(ctx, c.config.MinECode, mux, marshaler, w, r, err)
	}))
	// grpc stream error处理
	c.muxOptions = append(c.muxOptions, runtime.WithStreamErrorHandler(func(ctx context.Context, err error) *status.Status {
		return eerrors.FromError(err).GRPCStatus()
	}))
}

// withIncomintHeaderMatcherOption 传输自定义http头到grpc server
func withIncomingHeaderMatcherOption(c *Container) {
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

// withUrlPathTransOption 传输url到grpc server
func withUrlPathTransOption(c *Container) {
	c.muxOptions = append(c.muxOptions, runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		md := make(metadata.MD)
		md[urlMetadataName] = []string{req.URL.Path}
		return md
	}))
}

// withHttpResponseModifier 自定义响应头
func withHttpResponseModifier(c *Container) {
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
