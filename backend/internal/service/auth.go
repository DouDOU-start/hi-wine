package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type IAuthService interface {
	GetUserId(ctx context.Context) int64
	GetOpenid(ctx context.Context) string
}

var authService = localAuthService{}

type localAuthService struct{}

func Auth() IAuthService {
	return &authService
}

// 获取当前用户ID
func (s *localAuthService) GetUserId(ctx context.Context) int64 {
	value := g.RequestFromCtx(ctx).GetCtxVar("userId")
	if value.IsNil() {
		return 0
	}
	return value.Int64()
}

// 获取当前用户Openid
func (s *localAuthService) GetOpenid(ctx context.Context) string {
	value := g.RequestFromCtx(ctx).GetCtxVar("openid")
	if value.IsNil() {
		return ""
	}
	return value.String()
}
