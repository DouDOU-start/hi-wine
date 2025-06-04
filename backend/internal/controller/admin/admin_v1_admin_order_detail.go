package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

func (c *ControllerV1) AdminOrderDetail(ctx context.Context, req *v1.AdminOrderDetailReq) (res *v1.AdminOrderDetailRes, err error) {
	// 创建响应对象
	res = &v1.AdminOrderDetailRes{}

	// 调用订单服务获取订单详情
	orderService := service.Order()
	order, err := orderService.GetOrderDetailAdmin(ctx, req.OrderID)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取订单详情成功"
	res.Data = *order

	return res, nil
}
