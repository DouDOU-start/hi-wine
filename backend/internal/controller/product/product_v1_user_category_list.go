package product

import (
	"context"

	"backend/api/common"
	v1 "backend/api/product/v1"
	"backend/internal/service"
)

func (c *ControllerV1) UserCategoryList(ctx context.Context, req *v1.UserCategoryListReq) (res *v1.UserCategoryListRes, err error) {
	// 创建响应对象
	res = &v1.UserCategoryListRes{}

	// 调用分类服务获取激活的分类列表
	categoryService := service.Category{}
	categories, err := categoryService.GetActiveCategories(ctx)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取分类列表成功"

	// 转换为用户分类列表
	userCategories := make([]v1.UserCategory, len(categories))
	for i, category := range categories {
		userCategories[i] = v1.UserCategory{
			ID:        category.ID,
			Name:      category.Name,
			SortOrder: category.SortOrder,
		}
	}

	res.Data.List = userCategories

	return res, nil
}
