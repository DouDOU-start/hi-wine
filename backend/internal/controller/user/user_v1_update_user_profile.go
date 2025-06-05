package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/common"
	v1 "backend/api/user/v1"
	"backend/internal/service"
)

func (c *ControllerV1) UpdateUserProfile(ctx context.Context, req *v1.UpdateUserProfileReq) (res *v1.UpdateUserProfileRes, err error) {
	// 从上下文中获取当前登录用户ID
	userId := ctx.Value("userId")
	if userId == nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 调用服务层更新用户信息
	user, err := service.User().UpdateUserInfo(ctx, userId.(int64), req)
	if err != nil {
		return nil, gerror.New("更新用户信息失败：" + err.Error())
	}

	// 构建响应
	res = &v1.UpdateUserProfileRes{}
	res.Code = common.CodeSuccess
	res.Message = "更新成功"
	res.Data = *user

	return res, nil
}
