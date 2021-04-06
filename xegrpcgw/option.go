package xegrpcgw

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	urlMetadataName = "url"
)

// Option 可选项
type Option func(c *Container)

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

func urlPathTransOption(c *Container) {
	c.muxOptions = append(c.muxOptions, runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		md := make(metadata.MD)
		md[urlMetadataName] = []string{req.URL.Path}
		return md
	}))
}

func customerEcodeOption(c *Container) {
	c.muxOptions = append(c.muxOptions, runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		s := status.Convert(err)
		pb := s.Proto()
		if pb.Code >= 10000 {
			w.WriteHeader(http.StatusOK)
		}
		runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
	}))
	c.muxOptions = append(c.muxOptions, runtime.WithStreamErrorHandler(func(ctx context.Context, err error) *status.Status {
		return status.Convert(err)
	}))
}
