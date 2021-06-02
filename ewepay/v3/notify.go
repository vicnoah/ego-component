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

// ParseNotify 支付通知解析函数
// 需将参数封装为*http.Request,并且需携带[]byte body数据
// 返回resource为接口,当为支付通知时返回PayNotifyRequestResource类型,当为退款通知时返回RefundNotifyRequestResource
// 请进行类型断言使用
func (c *Component) ParseNotify(ctx context.Context, request *http.Request) (nt NotifyRequest, resource interface{}, err error) {
	var (
		alg = "AEAD_AES_256_GCM" // 加密算法
	)
	// 读取并拷贝request body
	bs, _ := io.ReadAll(request.Body)
	request.Body = io.NopCloser(bytes.NewBuffer(bs))
	c.mu.Lock()
	err = c.requestValidator(ctx, request)
	if err != nil {
		return
	}
	c.mu.Unlock()
	err = json.Unmarshal(bs, &nt)
	if err != nil {
		return
	}
	// 判断通知类型
	switch nt.EventType {
	case PaySuccessNotifyType: // 支付成功
		resource = PayNotifyRequestResource{}
	case RefundSuccessNotifyType: // 退款成功
		resource = RefundNotifyRequestResource{}
	case RefundAbnormalNotifyType: // 退款异常
		resource = RefundNotifyRequestResource{}
	case RefundCloseNotifyType: // 退款关闭
		resource = RefundNotifyRequestResource{}
	default:
		err = fmt.Errorf("pay err, event: %s", nt.EventType)
		return
	}
	c.mu.Lock()
	if nt.Resource.Algorithm != alg {
		err = fmt.Errorf("unsupported encryption algorithms: %s", nt.Resource.Algorithm)
		return
	}
	plaintext, err := utils.DecryptAES256GCM(c.config.AesKeyPasswd, nt.Resource.AssociatedData, nt.Resource.Nonce, nt.Resource.Ciphertext)
	c.mu.Unlock()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(plaintext), &resource)
	return
}
