package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"backend/api/common"
	v1 "backend/api/user/v1"
	"backend/internal/service"
)

func (c *ControllerV1) UserProfile(ctx context.Context, req *v1.UserProfileReq) (res *v1.UserProfileRes, err error) {
	// 从上下文中获取当前登录用户ID
	userId := g.RequestFromCtx(ctx).GetCtxVar("userId")
	if userId.IsEmpty() {
		g.Log().Error(ctx, "获取用户信息失败: 未找到用户ID")
		return nil, gerror.New("未登录或登录已过期")
	}

	g.Log().Debug(ctx, "获取用户信息，用户ID:", userId.Int64())

	// 调用服务层获取用户信息
	user, err := service.User().GetUserInfo(ctx, userId.Int64())
	if err != nil {
		g.Log().Error(ctx, "获取用户信息失败:", err.Error())
		return nil, gerror.New("获取用户信息失败：" + err.Error())
	}

	g.Log().Debug(ctx, "获取用户信息成功:", g.Map{
		"userId":    user.ID,
		"nickname":  user.Nickname,
		"hasAvatar": user.AvatarURL != "",
	})

	// 构建响应
	res = &v1.UserProfileRes{}
	res.Code = common.CodeSuccess
	res.Message = "获取成功"
	res.Data = *user

	return res, nil
}
