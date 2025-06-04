package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

func (c *ControllerV1) AdminOrderUpdateStatus(ctx context.Context, req *v1.AdminOrderUpdateStatusReq) (res *v1.AdminOrderUpdateStatusRes, err error) {
	// 创建响应对象
	res = &v1.AdminOrderUpdateStatusRes{}

	// 调用订单服务更新订单状态
	orderService := service.Order()
	order, err := orderService.UpdateOrderStatus(ctx, req.OrderID, req.Status, req.Reason)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "更新订单状态成功"
	res.Data = *order

	return res, nil
}
