package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

// 获取分类列表
func (c *ControllerV1) CategoryList(ctx context.Context, req *v1.AdminCategoryListReq) (res *v1.AdminCategoryListRes, err error) {
	categoryService := service.Category{}
	list, total, err := categoryService.List(ctx, req)
	if err != nil {
		return nil, err
	}

	res = &v1.AdminCategoryListRes{}
	res.Code = common.CodeSuccess
	res.Message = "获取分类列表成功"
	res.Data.List = list
	res.Data.Total = total
	return
}

// 创建分类
func (c *ControllerV1) CreateCategory(ctx context.Context, req *v1.AdminCreateCategoryReq) (res *v1.AdminCreateCategoryRes, err error) {
	categoryService := service.Category{}
	category, err := categoryService.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	res = &v1.AdminCreateCategoryRes{}
	res.Code = common.CodeSuccess
	res.Message = "创建分类成功"
	res.Data = *category
	return
}

// 更新分类
func (c *ControllerV1) UpdateCategory(ctx context.Context, req *v1.AdminUpdateCategoryReq) (res *v1.AdminUpdateCategoryRes, err error) {
	categoryService := service.Category{}
	category, err := categoryService.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	res = &v1.AdminUpdateCategoryRes{}
	res.Code = common.CodeSuccess
	res.Message = "更新分类成功"
	res.Data = *category
	return
}

// 删除分类
func (c *ControllerV1) DeleteCategory(ctx context.Context, req *v1.AdminDeleteCategoryReq) (res *v1.AdminDeleteCategoryRes, err error) {
	categoryService := service.Category{}
	err = categoryService.Delete(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	res = &v1.AdminDeleteCategoryRes{}
	res.Code = common.CodeSuccess
	res.Message = "删除分类成功"
	return
}

// 获取分类详情
func (c *ControllerV1) CategoryDetail(ctx context.Context, req *v1.AdminCategoryDetailReq) (res *v1.AdminCategoryDetailRes, err error) {
	categoryService := service.Category{}
	category, err := categoryService.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	res = &v1.AdminCategoryDetailRes{}
	res.Code = common.CodeSuccess
	res.Message = "获取分类详情成功"
	res.Data = *category
	return
}
