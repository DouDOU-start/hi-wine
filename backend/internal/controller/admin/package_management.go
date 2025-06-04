package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

// AdminPackageCreate 创建套餐
func (c *ControllerV1) AdminPackageCreate(ctx context.Context, req *v1.AdminPackageCreateReq) (res *v1.AdminPackageCreateRes, err error) {
	// 创建响应对象
	res = &v1.AdminPackageCreateRes{}

	// 调用服务创建套餐
	detail, err := service.Package().CreatePackage(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "创建套餐成功"
	res.Data = *detail

	return res, nil
}

// AdminPackageUpdate 更新套餐
func (c *ControllerV1) AdminPackageUpdate(ctx context.Context, req *v1.AdminPackageUpdateReq) (res *v1.AdminPackageUpdateRes, err error) {
	// 创建响应对象
	res = &v1.AdminPackageUpdateRes{}

	// 调用服务更新套餐
	detail, err := service.Package().UpdatePackage(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "更新套餐成功"
	res.Data = *detail

	return res, nil
}

// AdminPackageDelete 删除套餐
func (c *ControllerV1) AdminPackageDelete(ctx context.Context, req *v1.AdminPackageDeleteReq) (res *v1.AdminPackageDeleteRes, err error) {
	// 创建响应对象
	res = &v1.AdminPackageDeleteRes{}

	// 调用服务删除套餐
	err = service.Package().DeletePackage(ctx, req.PackageID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "删除套餐成功"

	return res, nil
}

// AdminPackageList 获取套餐列表
func (c *ControllerV1) AdminPackageList(ctx context.Context, req *v1.AdminPackageListReq) (res *v1.AdminPackageListRes, err error) {
	// 创建响应对象
	res = &v1.AdminPackageListRes{}

	// 调用服务获取套餐列表
	list, total, err := service.Package().GetPackageList(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取套餐列表成功"
	res.Data.List = list
	res.Data.Total = total

	return res, nil
}

// AdminPackageDetail 获取套餐详情
func (c *ControllerV1) AdminPackageDetail(ctx context.Context, req *v1.AdminPackageDetailReq) (res *v1.AdminPackageDetailRes, err error) {
	// 创建响应对象
	res = &v1.AdminPackageDetailRes{}

	// 调用服务获取套餐详情
	detail, err := service.Package().GetPackageDetail(ctx, req.PackageID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取套餐详情成功"
	res.Data = *detail

	return res, nil
}
