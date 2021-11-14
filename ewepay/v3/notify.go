package v3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// 微信小程序支付通知解析

const (
	// PaySuccessNotifyType 支付成功通知
	PaySuccessNotifyType = "TRANSACTION.SUCCESS"
	// RefundSuccessNotifyType 退款成功通知
	RefundSuccessNotifyType = "REFUND.SUCCESS"
	// RefundAbnormalNotifyType 退款异常通知
	RefundAbnormalNotifyType = "REFUND.ABNORMAL"
	// RefundCloseNotifyType 退款异常通知
	RefundCloseNotifyType = "REFUND.CLOSED"
)

// NotifyRequest 支付请求数据,来源为微信发送的支付结果
type NotifyRequest struct {
	ID           string    `json:"id"`
	CreateTime   time.Time `json:"create_time"`
	ResourceType string    `json:"resource_type"`
	EventType    string    `json:"event_type"`
	Resource     struct {
		Algorithm      string `json:"algorithm"`
		Ciphertext     string `json:"ciphertext"`
		Nonce          string `json:"nonce"`
		OriginalType   string `json:"original_type"`
		AssociatedData string `json:"associated_data"`
	} `json:"resource"`
	Summary string `json:"summary"`
}

// PayNotifyRequestResource 支付通知内容
type PayNotifyRequestResource struct {
	TransactionID string `json:"transaction_id"`
	Amount        struct {
		PayerTotal    int    `json:"payer_total"`
		Total         int    `json:"total"`
		Currency      string `json:"currency"`
		PayerCurrency string `json:"payer_currency"`
	} `json:"amount"`
	Mchid           string `json:"mchid"`
	TradeState      string `json:"trade_state"`
	BankType        string `json:"bank_type"`
	PromotionDetail []struct {
		Amount              int    `json:"amount"`
		WechatpayContribute int    `json:"wechatpay_contribute"`
		CouponID            string `json:"coupon_id"`
		Scope               string `json:"scope"`
		MerchantContribute  int    `json:"merchant_contribute"`
		Name                string `json:"name"`
		OtherContribute     int    `json:"other_contribute"`
		Currency            string `json:"currency"`
		StockID             string `json:"stock_id"`
		GoodsDetail         []struct {
			GoodsRemark    string `json:"goods_remark"`
			Quantity       int    `json:"quantity"`
			DiscountAmount int    `json:"discount_amount"`
			GoodsID        string `json:"goods_id"`
			UnitPrice      int    `json:"unit_price"`
		} `json:"goods_detail"`
	} `json:"promotion_detail"`
	SuccessTime time.Time `json:"success_time"`
	Payer       struct {
		Openid string `json:"openid"`
	} `json:"payer"`
	OutTradeNo     string `json:"out_trade_no"`
	Appid          string `json:"appid"`
	TradeStateDesc string `json:"trade_state_desc"`
	TradeType      string `json:"trade_type"`
	Attach         string `json:"attach"`
	SceneInfo      struct {
		DeviceID string `json:"device_id"`
	} `json:"scene_info"`
}

// RefundNotifyRequestResource 退款通知内容
type RefundNotifyRequestResource struct {
	Mchid               string    `json:"mchid"`
	TransactionID       string    `json:"transaction_id"`
	OutTradeNo          string    `json:"out_trade_no"`
	RefundID            string    `json:"refund_id"`
	OutRefundNo         string    `json:"out_refund_no"`
	RefundStatus        string    `json:"refund_status"`
	SuccessTime         time.Time `json:"success_time"`
	UserReceivedAccount string    `json:"user_received_account"`
	Amount              struct {
		Total       int `json:"total"`
		Refund      int `json:"refund"`
		PayerTotal  int `json:"payer_total"`
		PayerRefund int `json:"payer_refund"`
	} `json:"amount"`
}

// PayNotifyCallbackFunc 支付通知回调函数
type PayNotifyCallbackFunc func(ntr NotifyRequest, resource PayNotifyRequestResource)

// RefundNotifyCallbackFunc 退款通知回调函数
type RefundNotifyCallbackFunc func(ntr NotifyRequest, resource RefundNotifyRequestResource)

// ParseNotify 支付通知解析函数
// 需将参数封装为*http.Request,并且需携带[]byte body数据
func (c *Component) ParseNotify(ctx context.Context, request *http.Request, payCallback PayNotifyCallbackFunc, refundCallback RefundNotifyCallbackFunc) (err error) {
	var (
		alg   = "AEAD_AES_256_GCM" // 加密算法
		isPay = false              // 是否是支付通知,否则是退款通知
	)
	c.mu.RLock()
	defer c.mu.RUnlock()
	// 读取并拷贝request body
	bs, _ := io.ReadAll(request.Body)
	request.Body = io.NopCloser(bytes.NewBuffer(bs))
	err = c.requestValidator(ctx, request)
	if err != nil {
		return
	}
	var ntr NotifyRequest
	err = json.Unmarshal(bs, &ntr)
	if err != nil {
		return
	}
	// 判断通知类型
	switch ntr.EventType {
	case PaySuccessNotifyType: // 支付成功
		isPay = true
	case RefundSuccessNotifyType: // 退款成功
	case RefundAbnormalNotifyType: // 退款异常
	case RefundCloseNotifyType: // 退款关闭
	default:
		err = fmt.Errorf("pay err, event: %s", ntr.EventType)
		return
	}
	c.mu.Lock()
	if ntr.Resource.Algorithm != alg {
		c.mu.Unlock()
		err = fmt.Errorf("unsupported encryption algorithms: %s", ntr.Resource.Algorithm)
		return
	}
	plaintext, err := utils.DecryptAES256GCM(c.config.AesKeyPasswd, ntr.Resource.AssociatedData, ntr.Resource.Nonce, ntr.Resource.Ciphertext)
	c.mu.Unlock()
	if err != nil {
		return
	}
	if isPay {
		var resource PayNotifyRequestResource
		err = json.Unmarshal([]byte(plaintext), &resource)
		if err != nil {
			return
		}
		// 支付回调
		payCallback(ntr, resource)
		return
	}
	var resource RefundNotifyRequestResource
	err = json.Unmarshal([]byte(plaintext), &resource)
	if err != nil {
		return
	}
	// 退款回调
	refundCallback(ntr, resource)
	return
}
