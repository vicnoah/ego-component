package main

import (
	"context"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego-component/eetcd"
	"github.com/gotomicro/ego-component/eetcd/examples/helloworld"
	"github.com/gotomicro/ego-component/eetcd/registry"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
	"github.com/gotomicro/ego/server/egrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New().
		Invoker(invoker).
		Registry(EtcdRegistry).
		Serve(func() server.Server {
			server := egrpc.Load("server.grpc").Build()
			helloworld.RegisterGreeterServer(server.Server, &Greeter{server: server})
			return server
		}()).Run(); err != nil {
		elog.Panic("startup", elog.Any("err", err))
	}
}

var (
	EtcdClient   *eetcd.Component
	EtcdRegistry *registry.Component
)

func invoker() error {
	EtcdClient = eetcd.Load("etcd").Build()
	EtcdRegistry = registry.Load("registry").Build(registry.WithClientEtcd(EtcdClient))
	return nil
}

type Greeter struct {
	server *egrpc.Component
}

func (g Greeter) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if request.Name == "error" {
		return nil, status.Error(codes.Unavailable, "error")
	}

	return &helloworld.HelloReply{
		Message: "Hello EGO, I'm " + g.server.Address(),
	}, nil
}
