package utility

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	// ContextUserIDKey 上下文中存储用户ID的键
	ContextUserIDKey = "userId"
)

// GetUserID 从上下文中获取用户ID
func GetUserID(ctx context.Context) (int64, error) {
	// 从上下文中获取用户ID
	value := ctx.Value(ContextUserIDKey)
	if value == nil {
		return 0, gerror.New("未登录或登录已过期")
	}

	// 转换为int64类型
	userID, ok := value.(int64)
	if !ok {
		return 0, gerror.New("用户ID类型错误")
	}

	return userID, nil
}
