package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

func (c *ControllerV1) AdminProductCreate(ctx context.Context, req *v1.AdminProductCreateReq) (res *v1.AdminProductCreateRes, err error) {
	// 创建响应对象
	res = &v1.AdminProductCreateRes{}

	// 调用商品服务创建商品
	productService := &service.Product{}
	product, err := productService.Create(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "创建商品成功"
	res.Data = *product

	return res, nil
}
