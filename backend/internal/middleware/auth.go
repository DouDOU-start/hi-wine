package middleware

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ContextKey 上下文键类型
type ContextKey string

const (
	// UserIdKey 用户ID上下文键
	UserIdKey ContextKey = "userId"
)

// GetUserId 从上下文中获取用户ID
func GetUserId(ctx context.Context) int64 {
	value := ctx.Value(UserIdKey)
	if value == nil {
		return 0
	}

	userId, ok := value.(int64)
	if !ok {
		return 0
	}

	return userId
}

// SetUserId 设置用户ID到上下文
func SetUserId(ctx context.Context, userId int64) context.Context {
	return context.WithValue(ctx, UserIdKey, userId)
}

// Auth 认证中间件
func Auth(r *ghttp.Request) {
	// 从请求头中获取Token
	token := r.Header.Get("Authorization")
	if token == "" {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": "未登录或登录已过期",
		})
		r.Exit()
		return
	}

	// TODO: 验证Token并获取用户ID
	// 这里应该调用认证服务验证Token并获取用户ID
	// 为了示例，我们假设验证通过，用户ID为1
	userId := int64(1)

	// 将用户ID设置到上下文
	ctx := SetUserId(r.Context(), userId)
	r.SetCtx(ctx)

	r.Middleware.Next()
}
