package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

// AdminProductList 获取商品列表
func (c *ControllerV1) AdminProductList(ctx context.Context, req *v1.AdminProductListReq) (res *v1.AdminProductListRes, err error) {
	// 创建响应对象
	res = &v1.AdminProductListRes{}

	// 调用商品服务获取列表
	list, total, err := service.Product().List(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取商品列表成功"
	res.Data.List = list
	res.Data.Total = total

	return res, nil
}

// AdminProductCreate 创建商品
func (c *ControllerV1) AdminProductCreate(ctx context.Context, req *v1.AdminProductCreateReq) (res *v1.AdminProductCreateRes, err error) {
	// 创建响应对象
	res = &v1.AdminProductCreateRes{}

	// 调用商品服务创建商品
	product, err := service.Product().Create(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "创建商品成功"
	res.Data = *product

	return res, nil
}

// AdminProductUpdate 更新商品
func (c *ControllerV1) AdminProductUpdate(ctx context.Context, req *v1.AdminProductUpdateReq) (res *v1.AdminProductUpdateRes, err error) {
	// 创建响应对象
	res = &v1.AdminProductUpdateRes{}

	// 调用商品服务更新商品
	product, err := service.Product().Update(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "更新商品成功"
	res.Data = *product

	return res, nil
}

// AdminProductDelete 删除商品
func (c *ControllerV1) AdminProductDelete(ctx context.Context, req *v1.AdminProductDeleteReq) (res *v1.AdminProductDeleteRes, err error) {
	// 创建响应对象
	res = &v1.AdminProductDeleteRes{}

	// 调用商品服务删除商品
	err = service.Product().Delete(ctx, req.ProductID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "删除商品成功"

	return res, nil
}

// AdminProductDetail 获取商品详情
func (c *ControllerV1) AdminProductDetail(ctx context.Context, req *v1.AdminProductDetailReq) (res *v1.AdminProductDetailRes, err error) {
	// 创建响应对象
	res = &v1.AdminProductDetailRes{}

	// 调用商品服务获取详情
	productDetail, err := service.Product().GetProductDetail(ctx, req.ProductID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取商品详情成功"
	res.Data = productDetail.UserProduct

	return res, nil
}
