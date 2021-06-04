package v3

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"sync"
	"time"

	"github.com/gotomicro/ego/core/elog"
)

// 微信支付v3协议

const PackageName = "component.ewepay"

type (
	// Component ewepay 组件
	Component struct {
		mu                sync.Mutex // 锁
		name              string
		config            *Config
		mchPrivateKey     *rsa.PrivateKey     // 商户API私钥
		wechatPayCertList []*x509.Certificate // 微信支付平台证书
		logger            *elog.Component
	}
)

func newComponent(c *Container) *Component {
	// 初始化定时任务
	com := &Component{
		name:          c.name,
		logger:        c.logger,
		config:        c.config,
		mchPrivateKey: c.mchPrivateKey,
	}
	if !c.config.IsLoadCert {
		return com
	}
	// 初始化调用获取微信支付证书,随后证书应该使用定时任务自动更新
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	certs, err := com.GetCertificates(ctx)
	if err != nil {
		c.logger.Panic("download pay certificates error", elog.FieldErr(err), elog.FieldKey(PackageName))
	}
	com.wechatPayCertList = certs
	return com
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) PackageName() string {
	return PackageName
}

func (c *Component) Start() error {
	var err error
	return err
}

func (c *Component) Stop() error {
	return nil
}
