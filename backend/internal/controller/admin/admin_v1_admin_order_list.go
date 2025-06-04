package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

func (c *ControllerV1) AdminOrderList(ctx context.Context, req *v1.AdminOrderListReq) (res *v1.AdminOrderListRes, err error) {
	// 创建响应对象
	res = &v1.AdminOrderListRes{}

	// 调用订单服务获取订单列表
	orderService := service.Order()
	list, total, err := orderService.GetOrderList(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取订单列表成功"
	res.Data.List = list
	res.Data.Total = total

	return res, nil
}
