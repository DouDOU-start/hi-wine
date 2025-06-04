package product

import (
	"context"

	"backend/api/common"
	v1 "backend/api/product/v1"
	"backend/internal/service"
)

func (c *ControllerV1) UserProductDetail(ctx context.Context, req *v1.UserProductDetailReq) (res *v1.UserProductDetailRes, err error) {
	// 创建响应对象
	res = &v1.UserProductDetailRes{}

	// 调用商品服务获取商品详情
	productService := service.Product()
	productDetail, err := productService.GetProductDetail(ctx, req.ProductID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取商品详情成功"
	res.Data = *productDetail

	return res, nil
}
