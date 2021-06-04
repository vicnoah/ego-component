package v3

// Config GRPC Gateway config
type Config struct {
	MchID             string `toml:"mchId"`             // 微信支付商户号
	MchCertSN         string `toml:"mchCertSN"`         // 商户API证书序列号
	MchPrivateKeyPath string `toml:"mchPrivateKeyPath"` // 商户私钥路径,pem格式
	WechatMinAppID    string `toml:"wechatMinAppId"`    // 微信小程序appid
	AesKeyPasswd      string `toml:"aesKeyPasswd"`      // 微信aeskey密码
	NotifyURL         string `toml:"notifyUrl"`         // 微信支付通知url
	IsLoadCert        bool   `toml:"isLoadCert"`        // 是否初始化加载微信支付证书,默认加载
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		IsLoadCert: true,
	}
}
