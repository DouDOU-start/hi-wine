package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/common"
	v1 "backend/api/user/v1"
	"backend/internal/service"
)

// AuthController 认证控制器
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

// 实现 IUserV1 接口的其他方法，但返回未实现错误

func (c *AuthController) UserOrderList(ctx context.Context, req *v1.UserOrderListReq) (res *v1.UserOrderListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "认证控制器不支持此方法")
}

func (c *AuthController) UserOrderDetail(ctx context.Context, req *v1.UserOrderDetailReq) (res *v1.UserOrderDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "认证控制器不支持此方法")
}

func (c *AuthController) UserPackageList(ctx context.Context, req *v1.UserPackageListReq) (res *v1.UserPackageListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "认证控制器不支持此方法")
}

func (c *AuthController) UserPackageDetail(ctx context.Context, req *v1.UserPackageDetailReq) (res *v1.UserPackageDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "认证控制器不支持此方法")
}

func (c *AuthController) UserBuyPackage(ctx context.Context, req *v1.UserBuyPackageReq) (res *v1.UserBuyPackageRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "认证控制器不支持此方法")
}

func (c *AuthController) UserProfile(ctx context.Context, req *v1.UserProfileReq) (res *v1.UserProfileRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "认证控制器不支持此方法")
}

func (c *AuthController) UpdateUserProfile(ctx context.Context, req *v1.UpdateUserProfileReq) (res *v1.UpdateUserProfileRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "认证控制器不支持此方法")
}
