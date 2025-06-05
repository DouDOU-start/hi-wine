package order

import (
	"context"

	v1 "backend/api/order/v1"
	"backend/internal/service"
	"backend/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) CreateOrder(ctx context.Context, req *v1.CreateOrderReq) (res *v1.CreateOrderRes, err error) {
	// 1. 从上下文中获取用户ID
	userId, err := utility.GetUserID(ctx)
	if err != nil {
		g.Log().Error(ctx, "创建订单失败: 未找到用户ID", err)
		return nil, gerror.New("未登录或登录已过期")
	}

	g.Log().Debug(ctx, "创建订单，用户ID:", userId, "桌号二维码ID:", req.TableQrcodeID)

	// 2. 调用订单服务创建订单
	orderService := service.Order()
	orderRes, err := orderService.CreateOrder(ctx, req, userId)
	if err != nil {
		g.Log().Error(ctx, "创建订单失败:", err, "用户ID:", userId)
		return nil, err
	}

	g.Log().Info(ctx, "创建订单成功，用户ID:", userId, "订单ID:", orderRes.OrderID)

	return orderRes, nil
}
