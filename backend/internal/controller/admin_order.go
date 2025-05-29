package controller

import (
	"backend/api"
	"backend/internal/service"
	"context"
)

type AdminOrderController struct{}

func NewAdminOrder() *AdminOrderController {
	return &AdminOrderController{}
}

// 获取订单列表
func (c *AdminOrderController) List(ctx context.Context, req *api.AdminOrderListReq) (res *api.AdminOrderListRes, err error) {
	return service.AdminOrder().List(ctx, req)
}

// 获取订单详情
func (c *AdminOrderController) Detail(ctx context.Context, req *api.AdminOrderDetailReq) (res *api.AdminOrderDetailRes, err error) {
	return service.AdminOrder().Detail(ctx, req)
}

// 更新订单状态
func (c *AdminOrderController) UpdateStatus(ctx context.Context, req *api.AdminOrderUpdateStatusReq) (res *api.AdminOrderUpdateStatusRes, err error) {
	return service.AdminOrder().UpdateStatus(ctx, req)
}

// 获取订单统计
func (c *AdminOrderController) Stats(ctx context.Context, req *api.AdminOrderStatsReq) (res *api.AdminOrderStatsRes, err error) {
	return service.AdminOrder().Stats(ctx, req)
}

// 导出订单数据
func (c *AdminOrderController) Export(ctx context.Context, req *api.AdminOrderExportReq) (res interface{}, err error) {
	return service.AdminOrder().Export(ctx, req)
}
