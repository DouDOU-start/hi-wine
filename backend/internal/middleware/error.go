package middleware

import (
	"backend/api/common"
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()

	if err := r.GetError(); err != nil {
		// 已经处理过的错误直接返回
		if r.Response.Status != http.StatusOK {
			return
		}

		// 默认错误码
		code := common.CodeServerError
		msg := "服务器内部错误"

		// 获取自定义错误码
		if e := gerror.Unwrap(err); e != nil {
			if err1, ok := e.(gcode.Code); ok {
				// 如果是自定义错误码
				switch err1 {
				case gcode.CodeValidationFailed:
					code = common.CodeParamError
					msg = err.Error()
				case gcode.CodeNotAuthorized:
					code = common.CodeUnauthorized
					msg = "未授权或登录已过期"
				case gcode.CodeNotFound:
					code = common.CodeNotFound
					msg = "资源不存在"
				default:
					msg = err.Error()
				}
			} else {
				msg = err.Error()
			}
		}

		// 输出JSON响应
		r.Response.ClearBuffer()
		r.Response.WriteJson(common.Response[interface{}]{
			Code:    code,
			Message: msg,
			Data:    nil,
		})
	}
}
