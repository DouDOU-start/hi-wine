package product

import (
	"context"

	"backend/api/common"
	v1 "backend/api/product/v1"
	"backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) UserProductList(ctx context.Context, req *v1.UserProductListReq) (res *v1.UserProductListRes, err error) {
	// 创建响应对象
	res = &v1.UserProductListRes{
		Response: common.Response[struct {
			List  []v1.UserProduct `json:"list"`
			Total int              `json:"total"`
		}]{
			Code:    common.CodeSuccess,
			Message: "获取商品列表成功",
		},
	}

	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// 验证排序参数
	sortBy := req.SortBy
	if sortBy != "" && sortBy != "price" && sortBy != "sales_count" {
		sortBy = "" // 如果不是有效的排序字段，则使用默认排序
	}

	sortOrder := req.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc" // 默认降序
	}

	// 调用商品服务获取分类下的商品列表
	productService := service.Product()
	list, total, err := productService.GetProductsByCategory(ctx, req.CategoryID, page, limit, sortBy, sortOrder)
	if err != nil {
		g.Log().Error(ctx, "获取分类商品列表失败", err, map[string]interface{}{
			"category_id": req.CategoryID,
			"page":        page,
			"limit":       limit,
		})
		res.Code = common.CodeServerError
		res.Message = "获取商品列表失败"
		return res, nil
	}

	// 设置响应数据
	res.Data.List = list
	res.Data.Total = total

	return res, nil
}
