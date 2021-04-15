package xegrpcgw

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// Register grpc-gateway注册器
type Register func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

// WithServiceHandler 注册服务
func WithServiceHandler(f func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) Option {
	return func(c *Container) {
		f(c.ctx, c.mux, c.config.GrpcEndpoint, c.grpcDialOptions)
	}
}

// LogDetail 日志信息
type LogDetail struct {
	Action     string // 请求名称
	URL        string // 请求url
	OriginIP   string // 来源ip
	GrpcMethod string // 请求grpc方法
	Params     string // 请求参数
	Username   string // 用户名
	UserID     string // 用户id
}

// WithLogRecord 日志记录
func WithLogRecord(f func(lg LogDetail)) Option {
	return func(c *Container) {
		c.muxOptions = append(c.muxOptions, runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
			var (
				action     string
				url        string
				originIP   string
				grpcMethod string
				params     string
				username   string
				userid     string
			)
			sm, ok := runtime.ServerMetadataFromContext(ctx)
			if !ok {
				return nil
			}
			md := sm.TrailerMD
			for k, v := range md {
				fmt.Printf("key: %s, value: %s\n", k, v)
			}
			isLogWrite, ok := md[XAppLogWrite]
			if !ok {
				return nil
			}
			if len(isLogWrite) == 0 {
				return nil
			}
			if isLogWrite[0] != "true" {
				return nil
			}
			// 访问操作
			actions, ok := md[XAppAction]
			if !ok {
				return nil
			}
			if len(actions) == 0 {
				return nil
			}
			action = actions[0]
			// 请求url
			urls, ok := md[RequestURL]
			if !ok {
				return nil
			}
			if len(urls) == 0 {
				return nil
			}
			url = urls[0]
			// 来源ip
			originIPs, ok := md[XForWardedFor]
			if !ok {
				return nil
			}
			if len(originIPs) == 0 {
				return nil
			}
			originIP = originIPs[0]
			// grpc方法
			grpcMethods, ok := md[XAppGrpcMethod]
			if !ok {
				return nil
			}
			if len(grpcMethods) == 0 {
				return nil
			}
			grpcMethod = grpcMethods[0]
			// 参数信息
			paramses, ok := md[XAppRequestParams]
			if !ok {
				return nil
			}
			if len(paramses) == 0 {
				return nil
			} else {
				params = paramses[0]
			}
			// 用户名
			usernames, ok := md[XAppRequestUserName]
			if !ok {
				return nil
			}
			if len(usernames) == 0 {
				return nil
			} else {
				username = usernames[0]
			}
			// 用户id
			userids, ok := md[XAppRequestUserID]
			if !ok {
				return nil
			}
			if len(userids) == 0 {
				return nil
			} else {
				userid = userids[0]
			}
			f(LogDetail{
				Action:     action,
				URL:        url,
				OriginIP:   originIP,
				GrpcMethod: grpcMethod,
				Params:     params,
				Username:   username,
				UserID:     userid,
			})
			return nil
		}))
	}
}
