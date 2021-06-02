package v3

import (
	"crypto/rsa"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type Container struct {
	name          string
	config        *Config
	mchPrivateKey *rsa.PrivateKey // 商户API私钥
	logger        *elog.Component
}

// DefaultContainer 构造默认容器
func DefaultContainer() *Container {
	return &Container{
		name:   PackageName,
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

// Load ...
func Load(key string) *Container {
	c := DefaultContainer()
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
	}
	return c
}

// Build 构建组件
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	// 加载证书
	// 加载商户私钥
	var err error
	c.mchPrivateKey, err = utils.LoadPrivateKeyWithPath(c.config.MchPrivateKeyPath)
	if err != nil {
		c.logger.Panic("build error", elog.FieldErr(err), elog.FieldKey(PackageName))
	}
	return newComponent(c)
}
