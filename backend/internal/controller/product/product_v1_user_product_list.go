package product

import (
	"context"

	"backend/api/common"
	v1 "backend/api/product/v1"
	"backend/internal/service"
)

func (c *ControllerV1) UserProductList(ctx context.Context, req *v1.UserProductListReq) (res *v1.UserProductListRes, err error) {
	// 创建响应对象
	res = &v1.UserProductListRes{}

	// 调用商品服务获取分类下的商品列表
	productService := service.Product()
	list, total, err := productService.GetProductsByCategory(ctx, req.CategoryID, req.Page, req.Limit)
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
