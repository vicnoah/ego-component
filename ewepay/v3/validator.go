package v3

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.vectec.io/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"git.vectec.io/wechatpay-apiv3/wechatpay-go/core/consts"
	"git.vectec.io/wechatpay-apiv3/wechatpay-go/utils"
)

// 数据验证

// requestValidator 对来自微信支付的请求进行数据验证
// 需将参数封装为*http.Request,并且需携带[]byte body数据
func (c *Component) requestValidator(ctx context.Context, request *http.Request) (err error) {
	certificates := map[string]*x509.Certificate{}
	// 处理证书
	for _, certificate := range c.wechatPayCertList {
		serialNo := utils.GetCertificateSerialNumber(*certificate)
		certificates[serialNo] = certificate
	}
	ver := &verifiers.SHA256WithRSAVerifier{
		Certificates: certificates,
	}
	if ver == nil {
		err = fmt.Errorf("you must init WechatPayValidator with auth.Verifier")
		return
	}
	err = validateParameters(request)
	if err != nil {
		return
	}
	message, err := buildMessage(request)
	if err != nil {
		return
	}
	serialNumber := strings.TrimSpace(request.Header.Get(consts.WechatPaySerial))
	signature, err := base64.StdEncoding.DecodeString(strings.TrimSpace(request.Header.Get(consts.WechatPaySignature)))
	if err != nil {
		return fmt.Errorf("base64 decode string wechat pay signature err:%s", err.Error())
	}
	err = ver.Verify(ctx, serialNumber, message, string(signature))
	if err != nil {
		err = fmt.Errorf("validate verify fail serial=%s request-id=%s err=%s", serialNumber,
			strings.TrimSpace(request.Header.Get(consts.RequestID)), err)
		return
	}
	return
}

func validateParameters(request *http.Request) (err error) {
	if strings.TrimSpace(request.Header.Get(consts.WechatPaySerial)) == "" {
		return fmt.Errorf("empty %s", consts.WechatPaySerial)
	}
	if strings.TrimSpace(request.Header.Get(consts.WechatPaySignature)) == "" {
		return fmt.Errorf("empty %s", consts.WechatPaySignature)
	}
	if strings.TrimSpace(request.Header.Get(consts.WechatPayTimestamp)) == "" {
		return fmt.Errorf("empty %s", consts.WechatPayTimestamp)
	}
	if strings.TrimSpace(request.Header.Get(consts.WechatPayNonce)) == "" {
		return fmt.Errorf("empty %s", consts.WechatPayNonce)
	}
	timeStampStr := strings.TrimSpace(request.Header.Get(consts.WechatPayTimestamp))
	timeStamp, err := strconv.Atoi(timeStampStr)
	if err != nil {
		return fmt.Errorf("invalid timestamp:[%s] err:[%v]", timeStampStr, err)
	}
	if math.Abs(float64(timeStamp)-float64(time.Now().Unix())) >= consts.FiveMinute {
		return fmt.Errorf("timestamp=[%d] expires", timeStamp)
	}
	return nil
}

func buildMessage(request *http.Request) (message string, err error) {
	timeStamp := strings.TrimSpace(request.Header.Get(consts.WechatPayTimestamp))
	nonce := strings.TrimSpace(request.Header.Get(consts.WechatPayNonce))
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return "", fmt.Errorf("read request body err:[%s]", err.Error())
	}
	body := string(bodyBytes)
	request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	message = fmt.Sprintf("%s\n%s\n%s\n", timeStamp, nonce, body)
	return message, nil
}
