package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/common"
	v1 "backend/api/user/v1"
	"backend/internal/service"
)

func (c *ControllerV1) UserProfile(ctx context.Context, req *v1.UserProfileReq) (res *v1.UserProfileRes, err error) {
	// 从上下文中获取当前登录用户ID
	userId := ctx.Value("userId")
	if userId == nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 调用服务层获取用户信息
	user, err := service.User().GetUserInfo(ctx, userId.(int64))
	if err != nil {
		return nil, gerror.New("获取用户信息失败：" + err.Error())
	}

	// 构建响应
	res = &v1.UserProfileRes{}
	res.Code = common.CodeSuccess
	res.Message = "获取成功"
	res.Data = *user

	return res, nil
}
