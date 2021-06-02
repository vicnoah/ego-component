package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	v3 "github.com/vicnoah/ego-component/ewepay/v3"
)

var ewepay *v3.Component

// export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	err := ego.New().
		Invoker(invokeWepay).
		Cron().
		Run()
	if err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

func invokeWepay() error {
	ewepay = v3.Load("ewepay.one").Build()
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// defer cancel()
	// prepayID, err := ewepay.JsAPIPrepay(ctx, jsapi.PrepayRequest{
	// 	Description: core.String("通威旗舰店-罗非鱼饲料"),
	// 	OutTradeNo:  core.String("1217752501201407033233368029"),
	// 	Attach:      core.String("自定义数据说明"),
	// 	Amount: &jsapi.Amount{
	// 		Total:    core.Int32(1),
	// 		Currency: core.String("CNY"),
	// 	},
	// 	Payer: &jsapi.Payer{
	// 		Openid: core.String("o_VZy5SZzHXxb4KByKQ2bnJ-Cbms"),
	// 	},
	// })
	// if err != nil {
	// 	return err
	// }
	// jsObj, err := ewepay.WxRequestPayment(ctx, prepayID)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(jsObj)
	err := parseNt(ewepay)
	if err != nil {
		return err
	}
	return nil
}

func parseNt(wepay *v3.Component) (err error) {
	bs, err := base64.StdEncoding.DecodeString("eyJpZCI6IjdjZGUzZTA4LWMxYmYtNTAzYy05NGM4LTM0OGIzY2Y1N2FhMyIsImNyZWF0ZV90aW1lIjoiMjAyMS0wNi0wMlQyMzoyNjowMSswODowMCIsInJlc291cmNlX3R5cGUiOiJlbmNyeXB0LXJlc291cmNlIiwiZXZlbnRfdHlwZSI6IlRSQU5TQUNUSU9OLlNVQ0NFU1MiLCJzdW1tYXJ5Ijoi5pSv5LuY5oiQ5YqfIiwicmVzb3VyY2UiOnsib3JpZ2luYWxfdHlwZSI6InRyYW5zYWN0aW9uIiwiYWxnb3JpdGhtIjoiQUVBRF9BRVNfMjU2X0dDTSIsImNpcGhlcnRleHQiOiJkWnFtbVV6eUZPSFdoK2lWR3FJRWhQbnFsTWJlbHREcm1lMEVjUTNwa0ZxVS9YOXJKUnJyak5WNUxlRTVHV05ZS1NsaEZpTTJjcEljWXZkTFdXVElrbGRMZitMd2MrSmFMYjhzZlp2c3VrT0gzeVUwMERlWkhlem1uclA2bDlJRlFNNURERmNFdjd2VHRIdkFqQ0RvN1BFcFNqb0YzWEZ3SlFPeE5GbW1JUVdyaEtZclc1MUd6ZjdrOEVuYmhCays1VEFxWmpJeVNKQ1ZrdlBkSFdia2ZTdGhjekE4SGlGc0NVOHdMRENCOEoxOWVHTmE0N3dIaFlFaWY5Tk1aeDR1d0x3VTEybWx1blZVTkU0ZTYwVUt0OEFZejBOUDdtcXVYK1JSKytxclZjR2NYbENXUHkyREN4VUhpdy9Na2dKYk9xYWpVeWlvbzRCKzRtMUxTSzZYTjZuZU5XcGd6dUVNNVZUL04xZXE4eHVES1lYcWhDVXdPS2hUTmZHRElObW83a2J5ZW42akQyQWdzL2ZaZ0V3ZjcrWThFeGtUZkJZeU14MUR3RC9jUm8rU1JIenhsMllJWnhJMGUvdy9vUi9sN3NEVzNydVdTWEUxVkxSWnZvL3lrclp1M3VFZGV6emIrMzJ4TDJoa3VXTTh1K2p3Z1A5S00xVHJlcytnaFRRTC8xYjJIbloxVEkrNFFNY2lJbVJxeHEvMERJZ3djV1M0d29KM0dJZks3eDJ3Ykk5VXd6QWZCanFGOS93a0h3MGMrb0NPY1hoaURrejZBWE9wRy9lU25Obmc1dGxud1ZSZG13b2FLY3ozMWUzSE0zMS8iLCJhc3NvY2lhdGVkX2RhdGEiOiJ0cmFuc2FjdGlvbiIsIm5vbmNlIjoiUXVHTVZwbDZOeVd3In19")
	if err != nil {
		return
	}
	request := &http.Request{}
	request.Body = io.NopCloser(bytes.NewBuffer(bs))
	request.Header = make(http.Header)
	request.Header.Set("Wechatpay-Signature", "aH0ecdTH7FMeowfuHvs59dUKtZubZoiztgIxl6Osh6PuT167y2UDvLv9zW/M0gxm1w5dNdXMW1/Q0bSqXdezuXFaoDgSos2jVyFtHl2j5lkzH/zgjF/B/T9i8FR3Pr3l6/AhFjOp8ecB/I1AmjaSGOKm4qmXXaYW0tm/NIyuMjAB2XtC5CBrgllzkjCqbP4kQfoSXWdNTgAkvujx2DgVT4dtFoXTys0zzeI8My5Vg97T3rnx1LVCf+U7MMk850RCEOY4V92FEDV6ajFw/rkMay3qRnGLBKCDXrlMAxCZR5PAQ7IeOdkCf0x56+4VowKxKRgFiZ0ZYxsENzItxL/tFQ==")
	request.Header.Set("Wechatpay-Serial", "136F0674C080462AC4E91069D181A36")
	request.Header.Set("Wechatpay-Timestamp", "1622647567")
	request.Header.Set("Wechatpay-Nonce", "k3bTOg2GgDp1Z9TDXH4wMFapT5jTNsot")
	nt, resource, err := wepay.ParseNotify(context.TODO(), request)
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", nt)
	fmt.Printf("%+v\n", resource)
	return
}
