package v3

import (
	"context"
	"time"

	"github.com/gotomicro/ego/core/elog"
)

// AutoUpdateCert 自动更新证书
// 应每天凌晨调用接口进行证书更新
func (c *Component) AutoUpdateCert() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	certs, err := c.GetCertificates(ctx)
	if err != nil {
		c.logger.Error("upgrade pay certificates error", elog.FieldErr(err), elog.FieldKey(PackageName))
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.wechatPayCertList = certs
}
