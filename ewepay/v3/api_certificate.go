package v3

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"io"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// 证书相关操作接口

// GetCertificatesResponse 微信支付下载证书结果
type GetCertificatesResponse struct {
	Data []struct {
		EffectiveTime      time.Time `json:"effective_time"`
		EncryptCertificate struct {
			Algorithm      string `json:"algorithm"`
			AssociatedData string `json:"associated_data"`
			Ciphertext     string `json:"ciphertext"`
			Nonce          string `json:"nonce"`
		} `json:"encrypt_certificate"`
		ExpireTime time.Time `json:"expire_time"`
		SerialNo   string    `json:"serial_no"`
	} `json:"data"`
}

// GetCertificates 获取微信支付证书
func (c *Component) GetCertificates(ctx context.Context) (certs []*x509.Certificate, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	client, err := c.newClientWithoutValidator(ctx)
	if err != nil {
		return
	}
	result, err := client.Get(ctx, "https://api.mch.weixin.qq.com/v3/certificates")
	if err != nil {
		return
	}
	defer result.Response.Body.Close()
	cb, err := io.ReadAll(result.Response.Body)
	if err != nil {
		return
	}
	var resp GetCertificatesResponse
	err = json.Unmarshal(cb, &resp)
	if err != nil {
		return
	}
	for _, data := range resp.Data {
		plaintext, er := utils.DecryptAES256GCM(c.config.AesKeyPasswd, data.EncryptCertificate.AssociatedData, data.EncryptCertificate.Nonce, data.EncryptCertificate.Ciphertext)
		if er != nil {
			err = er
			return
		}
		cf, er := utils.LoadCertificate(plaintext)
		if er != nil {
			err = er
			return
		}
		certs = append(certs, cf)
	}
	return
}
