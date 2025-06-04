package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

// AdminPackageProductList 获取套餐包含的商品列表
func (c *ControllerV1) AdminPackageProductList(ctx context.Context, req *v1.AdminPackageProductListReq) (res *v1.AdminPackageProductListRes, err error) {
	// 创建响应对象
	res = &v1.AdminPackageProductListRes{}

	// 调用服务获取套餐包含的商品列表
	list, err := service.Package().GetPackageProducts(ctx, req.PackageID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取套餐商品列表成功"
	res.Data.List = list

	return res, nil
}

// AdminPackageAddProducts 为套餐添加商品
func (c *ControllerV1) AdminPackageAddProducts(ctx context.Context, req *v1.AdminPackageAddProductsReq) (res *v1.AdminPackageAddProductsRes, err error) {
	// 创建响应对象
	res = &v1.AdminPackageAddProductsRes{}

	// 调用服务为套餐添加商品
	err = service.Package().AddPackageProducts(ctx, req.PackageID, req.ProductIDs)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "添加套餐商品成功"

	return res, nil
}

// AdminPackageRemoveProduct 从套餐中移除单个商品
func (c *ControllerV1) AdminPackageRemoveProduct(ctx context.Context, req *v1.AdminPackageRemoveProductReq) (res *v1.AdminPackageRemoveProductRes, err error) {
	// 创建响应对象
	res = &v1.AdminPackageRemoveProductRes{}

	// 调用服务从套餐中移除商品
	err = service.Package().RemovePackageProduct(ctx, req.PackageID, req.ProductID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "移除套餐商品成功"

	return res, nil
}

// AdminPackageBatchRemoveProducts 批量从套餐中移除商品
func (c *ControllerV1) AdminPackageBatchRemoveProducts(ctx context.Context, req *v1.AdminPackageBatchRemoveProductsReq) (res *v1.AdminPackageBatchRemoveProductsRes, err error) {
	// 创建响应对象
	res = &v1.AdminPackageBatchRemoveProductsRes{}

	// 调用服务批量从套餐中移除商品
	err = service.Package().BatchRemovePackageProducts(ctx, req.PackageID, req.ProductIDs)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "批量移除套餐商品成功"

	return res, nil
}

// AdminPackageAvailableProducts 获取可添加到套餐的商品列表
func (c *ControllerV1) AdminPackageAvailableProducts(ctx context.Context, req *v1.AdminPackageAvailableProductsReq) (res *v1.AdminPackageAvailableProductsRes, err error) {
	// 创建响应对象
	res = &v1.AdminPackageAvailableProductsRes{}

	// 调用服务获取可添加到套餐的商品列表
	list, total, err := service.Package().GetAvailableProducts(ctx, req.PackageID, req.Keyword, req.Page, req.Limit)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取可添加商品列表成功"
	res.Data.List = list
	res.Data.Total = total

	return res, nil
}
