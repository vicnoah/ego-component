package xegrpcgw

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Container ...
type Container struct {
	config              *Config
	name                string
	logger              *elog.Component
	muxOptions          []runtime.ServeMuxOption
	ctx                 context.Context
	mux                 *runtime.ServeMux
	grpcDialOptions     []grpc.DialOption
	serviceRegisterFunc RegisterFunc // 服务注册方法
	configKey           string
	eginOptions         []egin.Option
	eginCorsOriginFunc  func(origin string) bool // egin跨域处理函数
	eg                  *egin.Component
}

// DefaultContainer ...
func DefaultContainer() *Container {
	return &Container{
		config:          DefaultConfig(),
		ctx:             context.Background(),
		grpcDialOptions: make([]grpc.DialOption, 0),
		eginOptions:     make([]egin.Option, 0),
		logger:          elog.DefaultLogger,
	}
}

// Load ...
func Load(key string) *Container {
	c := DefaultContainer()
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}
	c.configKey = key

	c.logger = c.logger.With(elog.FieldComponentName(key))
	c.name = key
	return c
}

// Build 构建组件
// dopt 参数一为日志记录特殊options
func (c *Container) Build(options ...Option) *egin.Component {
	// 初始化gin framework
	c.eg = egin.Load(c.configKey).Build(c.eginOptions...)
	// 跨域
	if c.config.EnableCors {
		corsIntercepter(c)
	}
	// 处理muxOptions
	// 处理http转grpc的header及ecode
	incomingHeaderMatcherOption(c)
	customerEcodeOption(c)
	httpResponseModifier(c)
	if c.config.EnableURLPathTrans {
		urlPathTransOption(c)
	}

	// 处理grpc-gateway及参数注入
	mux := runtime.NewServeMux(c.muxOptions...)
	c.mux = mux
	for _, option := range options {
		option(c)
	}
	// 注册grpc-gateway服务
	if c.serviceRegisterFunc != nil {
		c.serviceRegisterFunc(c.ctx, mux, c.config.GrpcEndpoint, c.grpcDialOptions)
	}
	// 注册http服务
	c.eg.Any(c.config.GinRelativePath, gin.WrapF(mux.ServeHTTP))

	// 注入handler
	server := newComponent(c.eg)
	return server
}
