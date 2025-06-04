package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/service"
)

// AdminOrderList 获取订单列表
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

// AdminOrderDetail 获取订单详情
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

// AdminOrderUpdateStatus 更新订单状态
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
