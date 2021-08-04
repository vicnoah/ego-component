package v3

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// WxRequestPaymentResponse 微信发起支付数据响应
type WxRequestPaymentResponse struct {
	TimeStamp string `json:"timeStamp"` // 时间戳
	NonceStr  string `json:"nonceStr"`  // 随机字符串
	Package   string `json:"package"`   // 订单详情扩展字
	SignType  string `json:"signType"`  // 签名方式
	PaySign   string `json:"paySign"`   // 签名
}

// WxRequestPayment 微信发起支付数据生成
// prepayID 预支付id
func (c *Component) WxRequestPayment(ctx context.Context, prepayID string) (jsonObj string, err error) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := strings.ToUpper(GetRandomString(32))
	pack := fmt.Sprintf("prepay_id=%s", prepayID)
	signType := "RSA"
	c.mu.Lock()
	paySign, err := utils.SignSHA256WithRSA(fmt.Sprintf("%s\n%s\n%s\n%s\n", c.config.WechatMinAppID, timeStamp, nonceStr, pack), c.mchPrivateKey)
	c.mu.Unlock()
	if err != nil {
		return
	}
	resp := WxRequestPaymentResponse{
		TimeStamp: timeStamp,
		NonceStr:  nonceStr,
		Package:   pack,
		SignType:  signType,
		PaySign:   paySign,
	}
	jb, err := json.Marshal(&resp)
	if err != nil {
		return
	}
	jsonObj = string(jb)
	return
}

// JsAPIPrepay 微信JSAPI支付下单
// opt示例:jsapi.PrepayRequest{
// 	Description: core.String("Image形象店-深圳腾大-QQ公仔"),
// 	OutTradeNo:  core.String("1217752501201407033233368019"),
// 	Attach:      core.String("自定义数据说明"),
// 	NotifyUrl:   core.String("https://www.weixin.qq.com/wxpay/pay.php"),
// 	Amount: &jsapi.Amount{
// 		Total:    core.Int32(1),
// 		Currency: core.String("CNY"),
// 	},
// 	Payer: &jsapi.Payer{
// 		Openid: core.String("o_VZy5SZzHXxb4KByKQ2bnJ-Cbms"),
// 	},
// }
func (c *Component) JsAPIPrepay(ctx context.Context, opt jsapi.PrepayRequest) (prepayID string, err error) {
	c.mu.Lock()
	client, err := c.newClient(ctx)
	if err != nil {
		c.mu.Unlock()
		return
	}
	opt.Appid = core.String(c.config.WechatMinAppID)
	opt.Mchid = core.String(c.config.MchID)
	opt.NotifyUrl = core.String(c.config.NotifyURL)
	c.mu.Unlock()
	svc := jsapi.JsapiApiService{Client: client}
	resp, _, err := svc.Prepay(ctx, opt)
	if err != nil {
		return
	}
	prepayID = *resp.PrepayId
	return
}

// JsAPICloseOrder 关闭订单
//
// 以下情况需要调用关单接口：
// 1. 商户订单支付失败需要生成新单号重新发起支付，要对原订单号调用关单，避免重复支付；
// 2. 系统下单后，用户支付超时，系统退出不再受理，避免用户继续，请调用关单接口。
func (c *Component) JsAPICloseOrder(ctx context.Context, opt jsapi.CloseOrderRequest) (err error) {
	c.mu.Lock()
	client, err := c.newClient(ctx)
	if err != nil {
		c.mu.Unlock()
		return
	}
	opt.Mchid = core.String(c.config.MchID)
	c.mu.Unlock()
	svc := jsapi.JsapiApiService{Client: client}
	_, err = svc.CloseOrder(ctx, opt)
	return
}

// JsAPIQueryOrderById 微信支付订单号查询订单
// 商户可以通过查询订单接口主动查询订单状态
func (c *Component) JsAPIQueryOrderById(ctx context.Context, opt jsapi.QueryOrderByIdRequest) (resp *payments.Transaction, err error) {
	c.mu.Lock()
	client, err := c.newClient(ctx)
	if err != nil {
		c.mu.Unlock()
		return
	}
	opt.Mchid = core.String(c.config.MchID)
	c.mu.Unlock()
	svc := jsapi.JsapiApiService{Client: client}
	resp, _, err = svc.QueryOrderById(ctx, opt)
	return
}

// JsAPIQueryOrderByOutTradeNo 商户订单号查询订单
// 商户可以通过查询订单接口主动查询订单状态
func (c *Component) JsAPIQueryOrderByOutTradeNo(ctx context.Context, opt jsapi.QueryOrderByOutTradeNoRequest) (resp *payments.Transaction, err error) {
	c.mu.Lock()
	client, err := c.newClient(ctx)
	if err != nil {
		c.mu.Unlock()
		return
	}
	opt.Mchid = core.String(c.config.MchID)
	c.mu.Unlock()
	svc := jsapi.JsapiApiService{Client: client}
	resp, _, err = svc.QueryOrderByOutTradeNo(ctx, opt)
	return
}

// JsAPIRefundRequest jsapi发起退款请求
type JsAPIRefundRequest struct {
	TransactionID string `json:"transaction_id,omitempty"`
	OutTradeNo    string `json:"out_trade_no,omitempty"`
	OutRefundNo   string `json:"out_refund_no,omitempty"`
	Reason        string `json:"reason,omitempty"`
	NotifyURL     string `json:"notify_url,omitempty"`
	FundsAccount  string `json:"funds_account,omitempty"`
	Amount        struct {
		Refund int `json:"refund"`
		From   []struct {
			Account string `json:"account"`
			Amount  int    `json:"amount"`
		} `json:"from,omitempty"`
		Total    int    `json:"total"`
		Currency string `json:"currency"`
	} `json:"amount"`
	GoodsDetail []struct {
		MerchantGoodsID  string `json:"merchant_goods_id"`
		WechatpayGoodsID string `json:"wechatpay_goods_id,omitempty"`
		GoodsName        string `json:"goods_name,omitempty"`
		UnitPrice        int    `json:"unit_price"`
		RefundAmount     int    `json:"refund_amount"`
		RefundQuantity   int    `json:"refund_quantity"`
	} `json:"goods_detail,omitempty"`
}

// JsAPIRefundResponse jsapi发起退款响应
type JsAPIRefundResponse struct {
	RefundID            string    `json:"refund_id"`
	OutRefundNo         string    `json:"out_refund_no"`
	TransactionID       string    `json:"transaction_id"`
	OutTradeNo          string    `json:"out_trade_no"`
	Channel             string    `json:"channel"`
	UserReceivedAccount string    `json:"user_received_account"`
	SuccessTime         time.Time `json:"success_time"`
	CreateTime          time.Time `json:"create_time"`
	Status              string    `json:"status"`
	FundsAccount        string    `json:"funds_account"`
	Amount              struct {
		Total  int `json:"total"`
		Refund int `json:"refund"`
		From   []struct {
			Account string `json:"account"`
			Amount  int    `json:"amount"`
		} `json:"from"`
		PayerTotal       int    `json:"payer_total"`
		PayerRefund      int    `json:"payer_refund"`
		SettlementRefund int    `json:"settlement_refund"`
		SettlementTotal  int    `json:"settlement_total"`
		DiscountRefund   int    `json:"discount_refund"`
		Currency         string `json:"currency"`
	} `json:"amount"`
	PromotionDetail []struct {
		PromotionID  string `json:"promotion_id"`
		Scope        string `json:"scope"`
		Type         string `json:"type"`
		Amount       int    `json:"amount"`
		RefundAmount int    `json:"refund_amount"`
		GoodsDetail  struct {
			MerchantGoodsID  string `json:"merchant_goods_id"`
			WechatpayGoodsID string `json:"wechatpay_goods_id"`
			GoodsName        string `json:"goods_name"`
			UnitPrice        int    `json:"unit_price"`
			RefundAmount     int    `json:"refund_amount"`
			RefundQuantity   int    `json:"refund_quantity"`
		} `json:"goods_detail"`
	} `json:"promotion_detail"`
}

// JsAPIRefund 发起退款请求
func (c *Component) JsAPIRefund(ctx context.Context, opt JsAPIRefundRequest) (resp JsAPIRefundResponse, err error) {
	c.mu.Lock()
	client, err := c.newClient(ctx)
	if err != nil {
		c.mu.Unlock()
		return
	}
	opt.NotifyURL = c.config.NotifyURL
	c.mu.Unlock()
	result, err := client.Post(ctx, "https://api.mch.weixin.qq.com/v3/refund/domestic/refunds", opt)
	if err != nil {
		return
	}
	defer result.Response.Body.Close()
	cb, err := io.ReadAll(result.Response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(cb, &resp)
	return
}

// JsAPIGetRefundResponse jsapi查询退款响应
type JsAPIGetRefundResponse struct {
	RefundID            string    `json:"refund_id"`
	OutRefundNo         string    `json:"out_refund_no"`
	TransactionID       string    `json:"transaction_id"`
	OutTradeNo          string    `json:"out_trade_no"`
	Channel             string    `json:"channel"`
	UserReceivedAccount string    `json:"user_received_account"`
	SuccessTime         time.Time `json:"success_time"`
	CreateTime          time.Time `json:"create_time"`
	Status              string    `json:"status"`
	FundsAccount        string    `json:"funds_account"`
	Amount              struct {
		Total  int `json:"total"`
		Refund int `json:"refund"`
		From   []struct {
			Account string `json:"account"`
			Amount  int    `json:"amount"`
		} `json:"from"`
		PayerTotal       int    `json:"payer_total"`
		PayerRefund      int    `json:"payer_refund"`
		SettlementRefund int    `json:"settlement_refund"`
		SettlementTotal  int    `json:"settlement_total"`
		DiscountRefund   int    `json:"discount_refund"`
		Currency         string `json:"currency"`
	} `json:"amount"`
	PromotionDetail []struct {
		PromotionID  string `json:"promotion_id"`
		Scope        string `json:"scope"`
		Type         string `json:"type"`
		Amount       int    `json:"amount"`
		RefundAmount int    `json:"refund_amount"`
		GoodsDetail  struct {
			MerchantGoodsID  string `json:"merchant_goods_id"`
			WechatpayGoodsID string `json:"wechatpay_goods_id"`
			GoodsName        string `json:"goods_name"`
			UnitPrice        int    `json:"unit_price"`
			RefundAmount     int    `json:"refund_amount"`
			RefundQuantity   int    `json:"refund_quantity"`
		} `json:"goods_detail"`
	} `json:"promotion_detail"`
}

// JsAPIGetRefund 查询退款
// outRefundNo 退款单号
func (c *Component) JsAPIGetRefund(ctx context.Context, outRefundNo string) (resp JsAPIGetRefundResponse, err error) {
	c.mu.Lock()
	client, err := c.newClient(ctx)
	if err != nil {
		c.mu.Unlock()
		return
	}
	c.mu.Unlock()
	result, err := client.Get(ctx, fmt.Sprintf("https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/%s", outRefundNo))
	if err != nil {
		return
	}
	defer result.Response.Body.Close()
	cb, err := io.ReadAll(result.Response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(cb, &resp)
	return
}

// JsAPITradeBillRequest 申请交易账单请求
type JsAPITradeBillRequest struct {
	BillDate string `json:"bill_date"`           // 账单日期
	BillType string `json:"bill_type,omitempty"` // 账单类型
	TarType  string `json:"tar_type,omitempty"`  // 压缩包类型
}

// JsAPITradeBillResponse 申请交易账单响应
type JsAPITradeBillResponse struct {
	DownloadURL string `json:"download_url"`
	HashType    string `json:"hash_type"`
	HashValue   string `json:"hash_value"`
}

// JsAPITradeBill 申请交易账单
func (c *Component) JsAPITradeBill(ctx context.Context, opt JsAPITradeBillRequest) (resp JsAPITradeBillResponse, err error) {
	c.mu.Lock()
	client, err := c.newClient(ctx)
	if err != nil {
		c.mu.Unlock()
		return
	}
	c.mu.Unlock()
	// Setup Query Params
	queryParams := url.Values{}
	queryParams.Add("bill_date", opt.BillDate)
	queryParams.Add("bill_type", opt.BillType)
	queryParams.Add("tar_type", opt.TarType)

	result, err := client.Get(ctx, fmt.Sprintf("https://api.mch.weixin.qq.com/v3/bill/tradebill?%s", queryParams.Encode()))
	if err != nil {
		return
	}
	defer result.Response.Body.Close()
	cb, err := io.ReadAll(result.Response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(cb, &resp)
	return
}

// JsAPIFundFlowBillRequest 申请资金账单请求
type JsAPIFundFlowBillRequest struct {
	BillDate string `json:"bill_date"`           // 账单日期
	BillType string `json:"bill_type,omitempty"` // 账单类型
	TarType  string `json:"tar_type,omitempty"`  // 压缩包类型
}

// JsAPIFundFlowBillResponse 申请资金账单响应
type JsAPIFundFlowBillResponse struct {
	DownloadURL string `json:"download_url"`
	HashType    string `json:"hash_type"`
	HashValue   string `json:"hash_value"`
}

// JsAPIFundFlowBill 申请资金账单
func (c *Component) JsAPIFundFlowBill(ctx context.Context, opt JsAPIFundFlowBillRequest) (resp JsAPIFundFlowBillResponse, err error) {
	c.mu.Lock()
	client, err := c.newClient(ctx)
	if err != nil {
		c.mu.Unlock()
		return
	}
	c.mu.Unlock()
	// Setup Query Params
	queryParams := url.Values{}
	queryParams.Add("bill_date", opt.BillDate)
	queryParams.Add("bill_type", opt.BillType)
	queryParams.Add("tar_type", opt.TarType)

	result, err := client.Get(ctx, fmt.Sprintf("https://api.mch.weixin.qq.com/v3/bill/fundflowbill?%s", queryParams.Encode()))
	if err != nil {
		return
	}
	defer result.Response.Body.Close()
	cb, err := io.ReadAll(result.Response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(cb, &resp)
	return
}

// GetRandomString 生成随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
