// 此文件已废弃，管理端分类列表接口请使用 /api/v1/admin/categories 路由及相关实现。

package product

import (
	"context"

	v1 "backend/api/product/v1"
	"backend/internal/service"
)

func (c *ControllerV1) CategoryList(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	// 创建响应对象
	res = &v1.CategoryListRes{}

	// 调用分类服务获取分类列表
	categoryService := service.Category{}
	categories, err := categoryService.GetActiveCategories(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	list := make([]v1.Category, len(categories))
	for i, category := range categories {
		list[i] = v1.Category{
			ID:        category.ID,
			Name:      category.Name,
			SortOrder: category.SortOrder,
			ImageURL:  category.ImageURL,
		}
	}

	// 设置响应数据
	res.List = list
	return res, nil
}
