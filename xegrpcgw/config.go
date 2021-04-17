package xegrpcgw

import (
	"fmt"

	"github.com/gotomicro/ego/core/eflag"
)

// Config GRPC Gateway config
type Config struct {
	Host                    string   // IP地址，默认0.0.0.0
	Port                    int      // PORT端口，默认8080
	EnableMetricInterceptor bool     // 是否开启监控，默认开启
	EnableTraceInterceptor  bool     // 是否开启链路追踪，默认开启
	EnableLocalMainIP       bool     // 自动获取ip地址
	EnableURLPathTrans      bool     // 是否开启url传递，开启后会将http请求url传递到grpc url metadata中，默认开启
	IncomingHeaders         []string // 允许传递给grpc的http请求头
	GrpcEndpoint            string   // grpc服务endpoint
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Host:                    eflag.String("host"),
		Port:                    8080,
		EnableTraceInterceptor:  false,
		EnableMetricInterceptor: false,
		EnableURLPathTrans:      true,
	}
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
