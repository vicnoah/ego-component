package xegrpcgw

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"

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
