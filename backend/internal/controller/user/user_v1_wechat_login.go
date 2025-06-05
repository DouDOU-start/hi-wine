package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "backend/api/user/v1"
)

// WechatLogin 微信登录（重定向到 /auth/wechat-login）
func (c *ControllerV1) WechatLogin(ctx context.Context, req *v1.WechatLoginReq) (res *v1.WechatLoginRes, err error) {
	// 这个方法不应该被调用，因为我们已经在 AuthController 中实现了该方法
	// 如果被调用，返回未实现错误
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "请使用 /auth/wechat-login 接口")
}
