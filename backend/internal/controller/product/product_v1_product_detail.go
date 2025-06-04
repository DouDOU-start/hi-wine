package product

import (
	"context"

	v1 "backend/api/product/v1"
	"backend/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) ProductDetail(ctx context.Context, req *v1.ProductDetailReq) (res *v1.ProductDetailRes, err error) {
	// 创建响应对象
	res = &v1.ProductDetailRes{}

	// 调用商品服务获取商品详情
	productService := service.Product()
	product, err := productService.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, gerror.New("商品不存在")
	}

	// 设置响应数据
	res.Product = *product

	return res, nil
}
