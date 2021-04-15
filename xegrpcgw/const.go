package xegrpcgw

// 原始请求头都将被自动转为小写
const (
	// TokenKey 登录tokenkey
	TokenKey = "x-app-token"
	// XForWardedFor 客户端ip
	XForWardedFor = "x-forwarded-for"
	// RequestURL 请求url
	RequestURL = "url"
	// XAppAction 操作名称
	XAppAction = "x-app-action"
	// XAppGrpcMethod Grpc方法,客户端访问grpc方法
	XAppGrpcMethod = "x-app-grpc-method"
	// XAppRequestParams 请求参数
	XAppRequestParams = "x-app-request-params"
	// XAppLogWrite 是否写入app访问日志
	XAppLogWrite = "x-app-log-write"
	// XAppRequestUserName 请求用户名
	XAppRequestUserName = "x-app-request-username"
	// XAppRequestUserID 请求用户
	XAppRequestUserID = "x-app-request-userid"
)
