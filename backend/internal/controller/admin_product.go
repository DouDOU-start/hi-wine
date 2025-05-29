package controller

import (
	"backend/api"
	"backend/internal/service"
	"context"
)

type AdminProductController struct{}

func NewAdminProduct() *AdminProductController {
	return &AdminProductController{}
}

// 获取商品列表
func (c *AdminProductController) List(ctx context.Context, req *api.AdminProductListReq) (res *api.AdminProductListRes, err error) {
	return service.AdminProduct().List(ctx, req)
}

// 获取商品详情
func (c *AdminProductController) Detail(ctx context.Context, req *api.AdminProductDetailReq) (res *api.AdminProductDetailRes, err error) {
	return service.AdminProduct().Detail(ctx, req)
}

// 添加商品
func (c *AdminProductController) Add(ctx context.Context, req *api.AdminProductAddReq) (res *api.AdminProductAddRes, err error) {
	return service.AdminProduct().Add(ctx, req)
}

// 更新商品
func (c *AdminProductController) Update(ctx context.Context, req *api.AdminProductUpdateReq) (res *api.AdminProductUpdateRes, err error) {
	return service.AdminProduct().Update(ctx, req)
}

// 删除商品
func (c *AdminProductController) Delete(ctx context.Context, req *api.AdminProductDeleteReq) (res *api.AdminProductDeleteRes, err error) {
	return service.AdminProduct().Delete(ctx, req)
}

// 更新商品状态
func (c *AdminProductController) UpdateStatus(ctx context.Context, req *api.AdminProductStatusUpdateReq) (res *api.AdminProductStatusUpdateRes, err error) {
	return service.AdminProduct().UpdateStatus(ctx, req)
}
