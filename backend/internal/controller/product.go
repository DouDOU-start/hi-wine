package controller

import (
	"backend/api"
	"backend/internal/service"
	"context"
)

type ProductController struct{}

func NewProduct() *ProductController {
	return &ProductController{}
}

// 获取商品列表
func (c *ProductController) List(ctx context.Context, req *api.ProductListReq) (res *api.ProductListRes, err error) {
	return service.Product().List(ctx, req)
}

// 获取商品详情
func (c *ProductController) Detail(ctx context.Context, req *api.ProductDetailReq) (res *api.ProductDetailRes, err error) {
	return service.Product().Detail(ctx, req)
}

// 获取商品分类列表
func (c *ProductController) CategoryList(ctx context.Context, req *api.CategoryListReq) (res *api.CategoryListRes, err error) {
	return service.Product().CategoryList(ctx, req)
}
