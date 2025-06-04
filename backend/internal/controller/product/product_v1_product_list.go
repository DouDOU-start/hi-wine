package product

import (
	"context"

	v1 "backend/api/product/v1"
	"backend/internal/service"
)

func (c *ControllerV1) ProductList(ctx context.Context, req *v1.ProductListReq) (res *v1.ProductListRes, err error) {
	// 创建响应对象
	res = &v1.ProductListRes{}

	// 调用商品服务获取分类下的商品列表
	productService := service.Product()
	products, _, err := productService.GetProductsByCategory(ctx, req.CategoryID, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	list := make([]v1.Product, len(products))
	for i, p := range products {
		list[i] = v1.Product{
			ID:          p.ID,
			Name:        p.Name,
			Price:       p.Price,
			ImageURL:    p.ImageURL,
			Stock:       p.Stock,
			Description: p.Description,
		}
	}

	// 设置响应数据
	res.List = list

	return res, nil
}
