package xegrpcgw

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"

	"github.com/gotomicro/ego/server"
)

// PackageName ...
const PackageName = "server.xegrpcgw"

// Component ...
type Component struct {
	mu     sync.Mutex
	name   string
	config *Config
	logger *elog.Component
	http.Handler
	Server   *http.Server
	listener net.Listener
}

// WithContext ...
func WithContext(ctx context.Context, mux *Component) *Component {
	return mux
}

// newComponent ...
func newComponent(name string, handler http.Handler, config *Config, logger *elog.Component) *Component {
	return &Component{
		name:     name,
		config:   config,
		logger:   logger,
		Handler:  handler,
		listener: nil,
	}
}

// Name 配置名称
func (c *Component) Name() string {
	return c.name
}

// PackageName 包名
func (c *Component) PackageName() string {
	return PackageName
}

// Init 初始化
func (c *Component) Init() error {
	listener, err := net.Listen("tcp", c.config.Address())
	if err != nil {
		c.logger.Panic("new egin server err", elog.FieldErrKind("listen err"), elog.FieldErr(err))
	}
	c.config.Port = listener.Addr().(*net.TCPAddr).Port
	c.listener = listener
	return nil
}

// Start implements server.Component interface.
func (c *Component) Start() error {
	// 因为start和stop在多个goroutine里，需要对Server上写锁
	c.mu.Lock()
	c.Server = &http.Server{
		Addr:    c.config.Address(),
		Handler: c,
	}
	c.mu.Unlock()
	err := c.Server.Serve(c.listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Stop implements server.Component interface
// it will terminate gin server immediately
func (c *Component) Stop() error {
	c.mu.Lock()
	err := c.Server.Close()
	c.mu.Unlock()
	return err
}

// GracefulStop implements server.Component interface
// it will stop gin server gracefully
func (c *Component) GracefulStop(ctx context.Context) error {
	c.mu.Lock()
	err := c.Server.Shutdown(ctx)
	c.mu.Unlock()
	return err
}

// Info returns server info, used by governor and consumer balancer
func (c *Component) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(c.listener.Addr().String()),
		server.WithKind(constant.ServiceConsumer),
	)
	return &info
}
