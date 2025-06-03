package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

func (c *ControllerV1) AdminProductUpdate(ctx context.Context, req *v1.AdminProductUpdateReq) (res *v1.AdminProductUpdateRes, err error) {
	// 创建响应对象
	res = &v1.AdminProductUpdateRes{}

	// 调用商品服务更新商品
	productService := &service.Product{}
	product, err := productService.Update(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "更新商品成功"
	res.Data = *product

	return res, nil
}
