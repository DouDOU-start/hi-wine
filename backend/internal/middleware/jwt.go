package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	jwtutil "backend/internal/utility/jwt"
)

// JwtAuth JWT认证中间件
func JwtAuth(r *ghttp.Request) {
	// 1. 获取 Authorization 头部
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    401,
			Message: "未登录或token缺失",
		})
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 2. 使用 utility/jwt 包解析 token
	claims, err := jwtutil.ParseToken(tokenString)
	if err != nil {
		g.Log().Error(r.Context(), "JWT解析失败:", err.Error())
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    401,
			Message: "token无效或已过期",
		})
		return
	}

	// 3. 将用户信息写入 context
	if claims != nil && claims.Data != nil {
		// 处理用户ID
		if userId, ok := claims.Data["userId"]; ok {
			// 将 userId 转换为 int64
			var userIdInt64 int64
			switch v := userId.(type) {
			case float64:
				userIdInt64 = int64(v)
			case int:
				userIdInt64 = int64(v)
			case int64:
				userIdInt64 = v
			default:
				// 尝试转换为float64
				if f, ok := userId.(float64); ok {
					userIdInt64 = int64(f)
				}
			}

			// 记录日志
			g.Log().Debug(r.Context(), "JWT认证成功，用户ID:", userIdInt64)

			// 设置到上下文
			r.SetCtxVar("userId", userIdInt64)
		}

		// 处理管理员ID
		if adminId, ok := claims.Data["id"]; ok {
			// 将 adminId 转换为 int
			var adminIdInt int
			switch v := adminId.(type) {
			case float64:
				adminIdInt = int(v)
			case int:
				adminIdInt = v
			case int64:
				adminIdInt = int(v)
			default:
				// 尝试转换为float64
				if f, ok := adminId.(float64); ok {
					adminIdInt = int(f)
				}
			}

			// 记录日志
			g.Log().Debug(r.Context(), "JWT认证成功，管理员ID:", adminIdInt)

			// 设置到上下文
			r.SetCtxVar("adminId", adminIdInt)
		}

		if openid, ok := claims.Data["openid"]; ok {
			r.SetCtxVar("openid", openid)
		}
	}

	r.Middleware.Next()
}
