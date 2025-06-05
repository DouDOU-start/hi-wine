package middleware

import (
	"backend/api/common"
	"reflect"

	"github.com/gogf/gf/v2/frame/g"
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

	// 如果响应为 nil，不处理
	if res == nil {
		g.Log().Debug(r.Context(), "响应为nil，路径:", r.URL.Path)
		return
	}

	// 记录响应类型，用于调试
	g.Log().Debug(r.Context(), "响应类型:", g.Map{
		"path": r.URL.Path,
		"type": reflect.TypeOf(res).String(),
	})

	// 如果响应已经是我们期望的格式，不再包装，但确保写入响应
	// 检查是否已经是标准响应格式
	if _, ok := res.(common.Response[any]); ok {
		g.Log().Debug(r.Context(), "已是标准响应格式(any)，直接写入:", r.URL.Path)
		r.Response.WriteJson(res)
		return
	}

	r.Response.WriteJson(res)
}
