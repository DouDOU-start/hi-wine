package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

// AdminUserList 获取用户列表
func (c *ControllerV1) AdminUserList(ctx context.Context, req *v1.AdminUserListReq) (res *v1.AdminUserListRes, err error) {
	// 初始化响应
	res = &v1.AdminUserListRes{}

	// 调用服务层获取用户列表
	list, total, err := service.AdminUser().GetUserList(ctx, req)
	if err != nil {
		res.Response = common.Response[struct {
			List  []v1.UserInfo `json:"list"`
			Total int           `json:"total"`
		}]{
			Code:    common.CodeServerError,
			Message: err.Error(),
			Data: struct {
				List  []v1.UserInfo `json:"list"`
				Total int           `json:"total"`
			}{
				List:  []v1.UserInfo{},
				Total: 0,
			},
		}
		return
	}

	// 设置成功响应
	res.Response = common.Response[struct {
		List  []v1.UserInfo `json:"list"`
		Total int           `json:"total"`
	}]{
		Code:    common.CodeSuccess,
		Message: "获取成功",
		Data: struct {
			List  []v1.UserInfo `json:"list"`
			Total int           `json:"total"`
		}{
			List:  list,
			Total: total,
		},
	}

	return
}

// AdminUserDetail 获取用户详情
func (c *ControllerV1) AdminUserDetail(ctx context.Context, req *v1.AdminUserDetailReq) (res *v1.AdminUserDetailRes, err error) {
	// 初始化响应
	res = &v1.AdminUserDetailRes{}

	// 调用服务层获取用户详情
	userInfo, err := service.AdminUser().GetUserDetail(ctx, req.UserID)
	if err != nil {
		res.Response = common.Response[v1.UserInfo]{
			Code:    common.CodeServerError,
			Message: err.Error(),
			Data:    v1.UserInfo{},
		}
		return
	}

	// 设置成功响应
	res.Response = common.Response[v1.UserInfo]{
		Code:    common.CodeSuccess,
		Message: "获取成功",
		Data:    *userInfo,
	}

	return
}
