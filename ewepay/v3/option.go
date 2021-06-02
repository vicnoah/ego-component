package v3

import "github.com/gotomicro/ego/core/elog"

type Option func(c *Container)

// WithName 设置name
func WithName(name string) Option {
	return func(c *Container) {
		c.name = name
	}
}

// WithLogger 设置 log
func WithLogger(logger *elog.Component) Option {
	return func(c *Container) {
		c.logger = logger
	}
}
