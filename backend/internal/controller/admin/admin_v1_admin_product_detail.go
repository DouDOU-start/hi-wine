package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

func (c *ControllerV1) AdminProductDetail(ctx context.Context, req *v1.AdminProductDetailReq) (res *v1.AdminProductDetailRes, err error) {
	// 创建响应对象
	res = &v1.AdminProductDetailRes{}

	// 调用商品服务获取商品详情
	product, err := service.Product().GetByID(ctx, req.ProductID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	if product == nil {
		res.Code = common.CodeNotFound
		res.Message = "商品不存在"
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取商品详情成功"
	res.Data = *product

	return res, nil
}
