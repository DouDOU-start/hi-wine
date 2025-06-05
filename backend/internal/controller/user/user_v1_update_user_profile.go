package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"backend/api/common"
	v1 "backend/api/user/v1"
	"backend/internal/service"
	"backend/internal/utility"
)

func (c *ControllerV1) UpdateUserProfile(ctx context.Context, req *v1.UpdateUserProfileReq) (res *v1.UpdateUserProfileRes, err error) {
	// 从上下文中获取当前登录用户ID
	userId, err := utility.GetUserID(ctx)
	if err != nil {
		g.Log().Error(ctx, "更新用户信息失败: 未找到用户ID", err)
		return nil, gerror.New("未登录或登录已过期")
	}

	g.Log().Debug(ctx, "更新用户信息，用户ID:", userId)

	// 调用服务层更新用户信息
	user, err := service.User().UpdateUserInfo(ctx, userId, req)
	if err != nil {
		g.Log().Error(ctx, "更新用户信息失败:", err)
		return nil, gerror.New("更新用户信息失败：" + err.Error())
	}

	// 构建响应
	res = &v1.UpdateUserProfileRes{}
	res.Code = common.CodeSuccess
	res.Message = "更新成功"
	res.Data = *user

	g.Log().Debug(ctx, "更新用户信息成功，用户ID:", userId)

	return res, nil
}
