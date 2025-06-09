package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

// AdminUserPackageCreate 创建用户套餐
func (c *ControllerV1) AdminUserPackageCreate(ctx context.Context, req *v1.AdminUserPackageCreateReq) (res *v1.AdminUserPackageCreateRes, err error) {
	// 创建响应对象
	res = &v1.AdminUserPackageCreateRes{}

	// 调用服务创建用户套餐
	detail, err := service.UserPackage().CreateUserPackage(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "创建用户套餐成功"
	res.Data = *detail

	return res, nil
}

// AdminUserPackageUpdateStatus 更新用户套餐状态
func (c *ControllerV1) AdminUserPackageUpdateStatus(ctx context.Context, req *v1.AdminUserPackageUpdateStatusReq) (res *v1.AdminUserPackageUpdateStatusRes, err error) {
	// 创建响应对象
	res = &v1.AdminUserPackageUpdateStatusRes{}

	// 调用服务更新用户套餐状态
	detail, err := service.UserPackage().UpdateUserPackageStatus(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "更新用户套餐状态成功"
	res.Data = *detail

	return res, nil
}

// AdminUserPackageList 获取用户套餐列表
func (c *ControllerV1) AdminUserPackageList(ctx context.Context, req *v1.AdminUserPackageListReq) (res *v1.AdminUserPackageListRes, err error) {
	// 创建响应对象
	res = &v1.AdminUserPackageListRes{}

	// 调用服务获取用户套餐列表
	list, total, err := service.UserPackage().GetUserPackageList(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取用户套餐列表成功"
	res.Data.List = list
	res.Data.Total = total

	return res, nil
}

// AdminUserPackageDetail 获取用户套餐详情
func (c *ControllerV1) AdminUserPackageDetail(ctx context.Context, req *v1.AdminUserPackageDetailReq) (res *v1.AdminUserPackageDetailRes, err error) {
	// 创建响应对象
	res = &v1.AdminUserPackageDetailRes{}

	// 调用服务获取用户套餐详情
	detail, err := service.UserPackage().GetUserPackageDetail(ctx, req.UserPackageID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取用户套餐详情成功"
	res.Data = *detail

	return res, nil
}

// AdminUserActivePackages 获取用户有效套餐
func (c *ControllerV1) AdminUserActivePackages(ctx context.Context, req *v1.AdminUserActivePackagesReq) (res *v1.AdminUserActivePackagesRes, err error) {
	// 创建响应对象
	res = &v1.AdminUserActivePackagesRes{}

	// 调用服务获取用户有效套餐
	list, err := service.UserPackage().GetUserActivePackages(ctx, req.UserID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取用户有效套餐成功"
	res.Data.List = list

	return res, nil
}

// AdminUserPackageFullDetail 获取用户套餐详细信息（包含用户信息、套餐信息和使用情况）
func (c *ControllerV1) AdminUserPackageFullDetail(ctx context.Context, req *v1.AdminUserPackageFullDetailReq) (res *v1.AdminUserPackageFullDetailRes, err error) {
	// 创建响应对象
	res = &v1.AdminUserPackageFullDetailRes{}

	// 调用服务获取用户套餐详细信息
	detail, err := service.UserPackage().GetUserPackageFullDetail(ctx, req.UserPackageID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取用户套餐详细信息成功"
	res.Data = *detail

	return res, nil
}
