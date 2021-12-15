package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vicnoah/ego-component/xegrpcgw"
	v1 "github.com/vicnoah/ego-component/xegrpcgw/examples/space/api/pb/v1"
	"google.golang.org/grpc"
)

// export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New().Serve(func() *xegrpcgw.Component {
		server := xegrpcgw.Load("server.test").Build(xegrpcgw.WithGrpcDialOptions(
			grpc.WithInsecure(),
		), xegrpcgw.WithService(func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
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
		server.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{})
		})
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}
