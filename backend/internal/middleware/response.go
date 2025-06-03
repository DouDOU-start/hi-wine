package middleware

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/net/ghttp"
)

// ResponseWrapper 统一响应格式的中间件
func ResponseWrapper(r *ghttp.Request) {
	// 先执行请求
	r.Middleware.Next()

	// 如果已经有错误响应或者已经被错误处理中间件处理过，不再处理
	if r.Response.Status >= 400 || r.GetError() != nil {
		return
	}

	// 获取原始响应数据
	res := r.GetHandlerResponse()

	// 如果响应已经是我们期望的格式，不再处理
	if _, ok := res.(common.Response[any]); ok {
		return
	}

	// 包装响应
	r.Response.ClearBuffer()
	r.Response.WriteJson(common.Response[any]{
		Code:    common.CodeSuccess,
		Message: "操作成功",
		Data:    res,
	})
}
