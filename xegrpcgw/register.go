package xegrpcgw

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// RegisterFunc grpc-gateway注册器
type RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

// WithService 注册服务
func WithService(f func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) Option {
	return func(c *Container) {
		c.serviceRegisterFunc = f
	}
}
