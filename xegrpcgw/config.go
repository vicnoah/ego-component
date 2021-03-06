package xegrpcgw

import (
	"fmt"
	"time"
)

// Config GRPC Gateway config
type Config struct {
	// egin部分配置
	Host                    string        // IP地址，默认0.0.0.0
	Port                    int           // PORT端口，默认8080
	Mode                    string        // gin的模式，默认是release模式
	EnableMetricInterceptor bool          // 是否开启监控，默认开启
	EnableTraceInterceptor  bool          // 是否开启链路追踪，默认开启
	EnableLocalMainIP       bool          // 自动获取ip地址
	SlowLogThreshold        time.Duration // 服务慢日志，默认500ms
	// grpc部分配置
	EnableCors                    bool     // 打开跨域，默认开启
	EnableURLPathTrans            bool     // 是否开启url传递，开启后会将http请求url传递到grpc url metadata中，默认开启
	AccessControlAllowOrigin      []string // 允许访问域名,默认:localhost
	AccessControlAllowHeaders     []string // 允许的header头,默认:"Content-Type", "AccessToken", "X-CSRF-Token", "Authorization"
	AccessControlAllowCredentials bool     // 设置为true,允许ajax异步请求带cookie信息,默认true
	AccessControlAllowMethods     []string // 允许请求方法,默认:"GET", "POST", "PUT", "DELETE", "OPTIONS"
	AccessControlExposeHeaders    []string // 允许跨域访问的header,默认:"Content-Length"
	IncomingHeaders               []string // 允许传递给grpc的http请求头,默认:空
	GrpcEndpoint                  string   // grpc服务endpoint
	MinECode                      int32    // 业务层最小grpc status code,默认: 10000
	GinRelativePath               string   // gin的相对路径(必须使用泛匹配),默认:/api/*action
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Host:                          "0.0.0.0",
		Port:                          8080,
		EnableTraceInterceptor:        false,
		EnableMetricInterceptor:       false,
		EnableCors:                    true,
		AccessControlAllowOrigin:      []string{"localhost"},
		AccessControlAllowHeaders:     []string{"Content-Type", "AccessToken", "X-CSRF-Token", "Authorization"},
		AccessControlAllowCredentials: true,
		AccessControlAllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AccessControlExposeHeaders:    []string{"Content-Length"},
		EnableURLPathTrans:            true,
		MinECode:                      10000,
		GinRelativePath:               "/api/*action",
	}
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
