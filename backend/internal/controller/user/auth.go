package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/common"
	v1 "backend/api/user/v1"
	"backend/internal/service"
)

// AuthController 认证控制器 - 仅处理身份认证相关功能
type AuthController struct{}

// NewAuth 创建认证控制器
func NewAuth() *AuthController {
	return &AuthController{}
}

// WechatLogin 微信登录
func (c *AuthController) WechatLogin(ctx context.Context, req *v1.WechatLoginReq) (res *v1.WechatLoginRes, err error) {
	// 调用服务层进行微信登录
	user, token, err := service.User().LoginByWechat(ctx, req.Code, req.Nickname, req.AvatarURL)
	if err != nil {
		return nil, gerror.New("微信登录失败：" + err.Error())
	}

	// 构建响应
	res = &v1.WechatLoginRes{}
	res.Code = common.CodeSuccess
	res.Message = "登录成功"
	res.Data = struct {
		Token string         `json:"token"`
		User  v1.UserProfile `json:"user"`
	}{
		Token: token,
		User:  *user,
	}

	return res, nil
}
