package xegrpcgw

import (
	"context"
	"net/http"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Container ...
type Container struct {
	config          *Config
	name            string
	logger          *elog.Component
	muxOptions      []runtime.ServeMuxOption
	ctx             context.Context
	mux             *runtime.ServeMux
	handlerFuncs    []http.HandlerFunc
	grpcDialOptions []grpc.DialOption
}

// DefaultContainer ...
func DefaultContainer() *Container {
	return &Container{
		config:          DefaultConfig(),
		logger:          elog.EgoLogger.With(elog.FieldComponent(PackageName)),
		ctx:             context.Background(),
		grpcDialOptions: make([]grpc.DialOption, 0),
	}
}

// Load ...
func Load(key string) *Container {
	c := DefaultContainer()
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}

	c.logger = c.logger.With(elog.FieldComponentName(key))
	c.name = key
	return c
}

func (c *Container) setGrpcOptions() {
	// 设置options
	c.grpcDialOptions = append(c.grpcDialOptions, grpc.WithInsecure(), withTracer(c))
}

// Build 构建组件
// dopt 参数一为日志记录特殊options
func (c *Container) Build(dopt Option, options ...Option) *Component {
	// 初始化选项
	// tracing
	if true {
		traceServerIntercepter(c)
	}
	// 度量
	if true {
		metricServerInterceptor(c)
	}
	c.setGrpcOptions()
	incomingHeaderMatcherOption(c)
	customerEcodeOption(c)
	httpResponseModifier(c)
	if c.config.EnableURLPathTrans {
		urlPathTransOption(c)
	}
	// 初始化特殊依赖option
	dopt(c)
	mux := runtime.NewServeMux(c.muxOptions...)
	c.mux = mux
	// 注入handler
	handlerInterceptor(c)
	for _, option := range options {
		option(c)
	}
	server := newComponent(c.name, c.mux, c.config, c.logger)
	return server
}
