package main

import (
	"context"

	"git.sabertrain.com/vector-dev/ego-component/xegrpcgw"
	v1 "git.sabertrain.com/vector-dev/ego-component/xegrpcgw/examples/space/api/pb/v1"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New().Serve(func() *xegrpcgw.Component {
		server := xegrpcgw.Load("server.test").Build(xegrpcgw.WithServiceHandler(func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
			// 用户服务网关
			err = v1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
			if err != nil {
				elog.Panic("start-gateway", elog.FieldErr(err))
			}
			// 视频服务网关
			err = v1.RegisterVideoServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
			if err != nil {
				elog.Panic("start-gateway", elog.FieldErr(err))
			}
			return
		}))
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}
