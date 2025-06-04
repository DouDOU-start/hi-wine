package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

func (c *ControllerV1) AdminProductList(ctx context.Context, req *v1.AdminProductListReq) (res *v1.AdminProductListRes, err error) {
	// 创建响应对象
	res = &v1.AdminProductListRes{}

	// 调用商品服务获取列表
	list, total, err := service.Product().List(ctx, req)
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
