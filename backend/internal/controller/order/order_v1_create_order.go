package order

import (
	"context"

	v1 "backend/api/order/v1"
	"backend/internal/middleware"
	"backend/internal/service"
)

func (c *ControllerV1) CreateOrder(ctx context.Context, req *v1.CreateOrderReq) (res *v1.CreateOrderRes, err error) {
	// 1. 从上下文中获取用户ID
	userId := middleware.GetUserId(ctx)
	// 暂时注释掉用户ID检查逻辑，后续可以根据实际需求决定是否需要
	// if userId <= 0 {
	// 	return nil, common.NewError(common.CodeUnauthorized, "未登录或登录已过期")
	// }

	// 2. 调用订单服务创建订单
	orderService := service.Order()
	orderRes, err := orderService.CreateOrder(ctx, req, userId)
	if err != nil {
		return nil, err
	}

	return orderRes, nil
}
