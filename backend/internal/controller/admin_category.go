package controller

import (
	"backend/api"
	"backend/internal/service"
	"context"
)

type AdminCategoryController struct{}

func NewAdminCategory() *AdminCategoryController {
	return &AdminCategoryController{}
}

// 获取分类列表
func (c *AdminCategoryController) List(ctx context.Context, req *api.AdminCategoryListReq) (res *api.AdminCategoryListRes, err error) {
	return service.AdminCategory().List(ctx, req)
}

// 获取分类详情
func (c *AdminCategoryController) Detail(ctx context.Context, req *api.AdminCategoryDetailReq) (res *api.AdminCategoryDetailRes, err error) {
	return service.AdminCategory().Detail(ctx, req)
}

// 添加分类
func (c *AdminCategoryController) Add(ctx context.Context, req *api.AdminCategoryAddReq) (res *api.AdminCategoryAddRes, err error) {
	return service.AdminCategory().Add(ctx, req)
}

// 更新分类
func (c *AdminCategoryController) Update(ctx context.Context, req *api.AdminCategoryUpdateReq) (res *api.AdminCategoryUpdateRes, err error) {
	return service.AdminCategory().Update(ctx, req)
}

// 删除分类
func (c *AdminCategoryController) Delete(ctx context.Context, req *api.AdminCategoryDeleteReq) (res *api.AdminCategoryDeleteRes, err error) {
	return service.AdminCategory().Delete(ctx, req)
}
