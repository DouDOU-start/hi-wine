package product

import (
	"context"

	"backend/api/common"
	v1 "backend/api/product/v1"
	"backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) UserCategoryList(ctx context.Context, req *v1.UserCategoryListReq) (res *v1.UserCategoryListRes, err error) {
	// 创建响应对象
	res = &v1.UserCategoryListRes{
		Response: common.Response[struct {
			List  []v1.UserCategory `json:"list"`
			Total int               `json:"total"`
		}]{
			Code:    common.CodeSuccess,
			Message: "获取分类列表成功",
		},
	}

	// 调用分类服务获取激活的分类列表
	categoryService := service.Category{}
	categories, err := categoryService.GetActiveCategories(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取分类列表失败", err)
		res.Code = common.CodeServerError
		res.Message = "获取分类列表失败"
		return res, nil
	}

	// 转换为用户分类列表
	userCategories := make([]v1.UserCategory, len(categories))
	for i, category := range categories {
		userCategories[i] = v1.UserCategory{
			ID:        category.ID,
			Name:      category.Name,
			SortOrder: category.SortOrder,
			ImageURL:  category.ImageURL,
		}

		// 获取每个分类下的商品数量
		productService := service.Product()
		productCount, _, _ := productService.GetProductsByCategory(ctx, category.ID, 1, 1, "", "")
		userCategories[i].ProductCount = len(productCount)
	}

	// 设置响应数据
	res.Data.List = userCategories
	res.Data.Total = len(userCategories)

	return res, nil
}
