package middleware

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ResponseWrapper 统一响应格式的中间件
func ResponseWrapper(r *ghttp.Request) {
	// 先执行请求
	r.Middleware.Next()

	// 如果已经有错误响应，不处理
	if r.Response.Status >= 400 {
		return
	}

	// 获取原始响应数据
	var (
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = 0
		msg  = "操作成功"
	)

	if err != nil {
		code = 500
		msg = gerror.Current(err).Error()
	}

	// 包装响应
	r.Response.ClearBuffer()
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    code,
		Message: msg,
		Data:    res,
	})
}
