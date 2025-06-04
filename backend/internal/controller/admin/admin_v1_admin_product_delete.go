package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

func (c *ControllerV1) AdminProductDelete(ctx context.Context, req *v1.AdminProductDeleteReq) (res *v1.AdminProductDeleteRes, err error) {
	// 创建响应对象
	res = &v1.AdminProductDeleteRes{}

	// 调用商品服务删除商品
	err = service.Product().Delete(ctx, req.ProductID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "删除商品成功"

	return res, nil
}
