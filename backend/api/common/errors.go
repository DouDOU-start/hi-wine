package common

import (
	"github.com/gogf/gf/v2/errors/gerror"
)

// 常用错误定义
var (
	ErrNotAuthenticated = gerror.New("用户未认证")
	ErrNotAuthorized    = gerror.New("用户无权限")
	ErrNotFound         = gerror.New("资源不存在")
	ErrParamInvalid     = gerror.New("参数无效")
	ErrServerInternal   = gerror.New("服务器内部错误")
)
