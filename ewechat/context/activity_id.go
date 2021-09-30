package context

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/vicnoah/ego-component/ewechat/util"
)

const (
	// ActivityIdURL 获取activity_id的接口
	ActivityIdURL = "https://api.weixin.qq.com/cgi-bin/message/wxopen/activityid/create"
)

// ResActivityId struct
type ResActivityId struct {
	util.CommonError

	// ActivityId 动态消息的 ID
	ActivityId string `json:"activity_id"`
	// ExpirationTime activity_id 的过期时间戳。默认24小时后过期
	ExpirationTime int64 `json:"expiration_time"`
}

// GetActivityIdFunc 获取 activityId 的函数签名
type GetActivityIdFunc func(ctx *Context) (activityId string, err error)

// SetActivityIdLock 设置读写锁（一个appID一个读写锁）
func (ctx *Context) SetActivityIdLock(l *sync.RWMutex) {
	ctx.activityIdLock = l
}

// SetGetActivityIdFunc 设置自定义获取activityId的方式, 需要自己实现缓存
func (ctx *Context) SetGetActivityIdFunc(f GetActivityIdFunc) {
	ctx.activityIdFunc = f
}

// GetActivityId 获取activity_id
func (ctx *Context) GetActivityId(accessToken, openid string) (activityId string, err error) {
	ctx.activityIdLock.Lock()
	defer ctx.activityIdLock.Unlock()

	if ctx.activityIdFunc != nil {
		return ctx.activityIdFunc(ctx)
	}
	activityIdCacheKey := fmt.Sprintf("activity_id_%s", openid)
	activityId, err = ctx.Cache.Get(context.Background(), activityIdCacheKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", nil
	}
	if activityId != "" {
		return activityId, nil
	}

	// 从微信服务器获取
	var resActivityId ResActivityId
	resActivityId, err = ctx.GetAcitvityIdFromServer(accessToken, openid)
	if err != nil {
		err = fmt.Errorf("get activityId err %w", err)
		return
	}

	activityId = resActivityId.ActivityId
	return
}

// GetAcitvityIdFromServer 强制从微信服务器获取ID
// 存入redis忽略掉过期时间，永久有效
func (ctx *Context) GetAcitvityIdFromServer(accessToken, openid string) (resActivityId ResActivityId, err error) {
	url := fmt.Sprintf("%s?access_token=%s&unionid=%s", ActivityIdURL, accessToken, openid)
	var body []byte
	body, err = ctx.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resActivityId)
	if err != nil {
		err = fmt.Errorf("activity_id from server parse json err %w", err)
		return
	}
	if resActivityId.ErrCode != 0 {
		err = fmt.Errorf("get activity_id err : errcode=%v, errormsg=%v", resActivityId.ErrCode, resActivityId.ErrMsg)
		return
	}

	activityIdCacheKey := fmt.Sprintf("activity_id_%s", openid)
	// expirationTime := resActivityId.ExpirationTime
	err = ctx.Cache.Set(context.Background(), activityIdCacheKey, resActivityId.ActivityId, 0)
	if err != nil {
		err = fmt.Errorf("set activity_id error %w", err)
		return
	}
	return
}
