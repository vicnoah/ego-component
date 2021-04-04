package xegrpcgw

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Register grpc-gateway注册器
type Register func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

func WithServiceHandler(f func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) Option {
	return func(c *Container) {
		f(c.ctx, c.mux, c.config.GrpcEndpoint, c.grpcDialOptions)
	}
}
