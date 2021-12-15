package xegrpcgw

import (
	"github.com/gotomicro/ego/server/egin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// 外部可用自定义配置

// WithEGinOptions 处理egin options
func WithEGinOptions(options ...egin.Option) Option {
	return func(c *Container) {
		c.eginOptions = options
	}
}

// WithGinCorsOriginFunc 处理egin跨域函数
func WithGinCorsOriginFunc(f func(origin string) bool) Option {
	return func(c *Container) {
		c.eginCorsOriginFunc = f
	}
}

// WithGrpcDialOptions 处理grpc dial options
func WithGrpcDialOptions(options ...grpc.DialOption) Option {
	return func(c *Container) {
		c.grpcDialOptions = options
	}
}

// WithMuxOptions 处理mux options
func WithMuxOptions(options ...runtime.ServeMuxOption) Option {
	return func(c *Container) {
		c.muxOptions = options
	}
}
