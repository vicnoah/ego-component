package main

import (
	"context"
	"fmt"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego-component/eredis"
	"github.com/gotomicro/ego/core/elog"
	"github.com/vicnoah/ego-component/ewechat"
)

func main() {
	err := ego.New().Invoker(
		invokerRedis,
		invokerEwechat,
		getActivityId,
	).Run()
	if err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

var ewechatClient *ewechat.Component

func invokerEwechat() error {
	ewechatClient = ewechat.Load("config.wechat").Build(ewechat.WithRedis(eredisClient))
	return nil
}

var eredisClient *eredis.Component

func invokerRedis() error {
	eredisClient = eredis.Load("redis.test").Build()
	str, err := eredisClient.Get(context.Background(), "min")
	fmt.Println(str, err)
	return nil
}

func getActivityId() error {
	min := ewechatClient.GetMiniProgram()
	fmt.Printf("MiniProgram----------%#v\n", min)
	accessToken, err := min.Context.GetAccessToken()
	if err != nil {
		return err
	}
	fmt.Printf("accessToken----------:\t%s\n", accessToken)
	activityId, err := min.GetActivityId(accessToken, "oEKJN4x_mTtdb3RAygTpY7P8_Ob4")
	if err != nil {
		return err
	}
	fmt.Printf("activityId-----------:\t%s\n", activityId)
	return nil
}
