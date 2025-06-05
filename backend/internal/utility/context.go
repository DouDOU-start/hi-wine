package utility

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	// ContextUserIDKey 上下文中存储用户ID的键
	ContextUserIDKey = "userId"
)

// GetUserID 从上下文中获取用户ID
// 支持从context.Context和ghttp.Request中获取用户ID
func GetUserID(ctx context.Context) (int64, error) {
	// 先尝试从请求上下文变量中获取
	req := g.RequestFromCtx(ctx)
	if req != nil {
		userId := req.GetCtxVar(ContextUserIDKey)
		if !userId.IsEmpty() {
			return userId.Int64(), nil
		}
	}

	// 再尝试从context.Value中获取
	value := ctx.Value(ContextUserIDKey)
	if value == nil {
		return 0, gerror.New("未登录或登录已过期")
	}

	// 尝试转换为int64类型
	switch v := value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return 0, gerror.New("用户ID类型错误")
	}
}
