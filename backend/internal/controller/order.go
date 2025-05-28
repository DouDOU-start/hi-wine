package controller

import (
	"backend/api"
	"backend/internal/service"
	"context"
)

type OrderController struct{}

func NewOrder() *OrderController {
	return &OrderController{}
}

// 创建订单
func (c *OrderController) Create(ctx context.Context, req *api.OrderCreateReq) (res *api.OrderCreateRes, err error) {
	return service.Order().Create(ctx, req)
}

// 获取订单列表
func (c *OrderController) List(ctx context.Context, req *api.OrderListReq) (res *api.OrderListRes, err error) {
	return service.Order().List(ctx, req)
}

// 获取订单详情
func (c *OrderController) Detail(ctx context.Context, req *api.OrderDetailReq) (res *api.OrderDetailRes, err error) {
	return service.Order().Detail(ctx, req)
}

// 更新订单状态
func (c *OrderController) UpdateStatus(ctx context.Context, req *api.OrderUpdateStatusReq) (res *api.OrderUpdateStatusRes, err error) {
	return service.Order().UpdateStatus(ctx, req)
}
