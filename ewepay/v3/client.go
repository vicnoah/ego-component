package v3

import (
	"context"
	"net/http"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
)

// newClient 新建支付http客户端
// 注意创建客户端涉及到数据操作可能发生竞争,故应该加锁使用
func (c *Component) newClient(ctx context.Context) (client *core.Client, err error) {
	client, err = core.NewClient(ctx, c.normalOptions()...)
	return
}

// newClientWithoutValidator 非验签http客户端
// 注意创建客户端涉及到数据操作可能发生竞争,故应该加锁使用
func (c *Component) newClientWithoutValidator(ctx context.Context) (client *core.Client, err error) {
	client, err = core.NewClient(ctx, c.withOutValidatorOptions()...)
	return
}

// withOutValidator 不验签的options生成
func (c *Component) withOutValidatorOptions() (opts []core.ClientOption) {
	var (
		customHTTPClient *http.Client
	)
	opts = append(opts, option.WithMerchantCredential(c.config.MchID, c.config.MchCertSN, c.mchPrivateKey)) // 必要，使用商户信息生成默认 WechatPayCredential
	opts = append(opts, option.WithoutValidator())                                                          // 必要，使用微信支付平台证书列表生成默认,此处忽略签名验证
	opts = append(opts, option.WithHTTPClient(customHTTPClient))                                            // 可选，设置自定义 HTTPClient 实例，不设置时使用默认 http.Client{}
	// opts = append(opts, core.WithTimeout(2*time.Second))                                                  // 可选，设置自定义超时时间，不设置时使用 http.Client{} 默认超时
	return
}

// normalOptions 普通支付options生成,此option有验签操作
func (c *Component) normalOptions() (opts []core.ClientOption) {
	var (
		customHTTPClient *http.Client
	)
	opts = append(opts, option.WithMerchantCredential(c.config.MchID, c.config.MchCertSN, c.mchPrivateKey)) // 必要，使用商户信息生成默认 WechatPayCredential
	opts = append(opts, option.WithWechatPayCertificate(c.wechatPayCertList))                               // 必要，使用微信支付平台证书列表生成默认,此处忽略签名验证
	opts = append(opts, option.WithHTTPClient(customHTTPClient))                                            // 可选，设置自定义 HTTPClient 实例，不设置时使用默认 http.Client{}
	// opts = append(opts, core.WithTimeout(2*time.Second))                                                  // 可选，设置自定义超时时间，不设置时使用 http.Client{} 默认超时
	return
}
