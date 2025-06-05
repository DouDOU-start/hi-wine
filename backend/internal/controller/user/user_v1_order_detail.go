package user

import (
	"context"

	userv1 "backend/api/user/v1"
	"backend/internal/service"
	"backend/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// UserOrderDetail 获取指定订单的详细信息
func (c *ControllerV1) UserOrderDetail(ctx context.Context, req *userv1.UserOrderDetailReq) (res *userv1.UserOrderDetailRes, err error) {
	// 1. 从上下文中获取用户ID
	userId, err := utility.GetUserID(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取订单详情失败: 未找到用户ID", err)
		return nil, gerror.New("未登录或登录已过期")
	}

	g.Log().Debug(ctx, "获取订单详情，用户ID:", userId, "订单ID:", req.OrderID)

	// 2. 调用订单服务获取订单详情
	orderService := service.Order()
	order, err := orderService.GetOrderDetail(ctx, req.OrderID, userId)
	if err != nil {
		g.Log().Error(ctx, "获取订单详情失败:", err, "用户ID:", userId, "订单ID:", req.OrderID)
		return nil, err
	}

	// 3. 转换为API响应格式
	items := make([]userv1.OrderItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = userv1.OrderItem{
			ProductID:     item.ProductID,
			Name:          item.Name,
			Quantity:      item.Quantity,
			ItemPrice:     item.ItemPrice,
			IsPackageItem: item.IsPackageItem,
			UserPackageID: item.UserPackageID,
			Notes:         item.Notes,
		}
	}

	userOrder := userv1.UserOrder{
		ID:            order.ID,
		OrderSN:       order.OrderSN,
		Status:        order.Status,
		TotalAmount:   order.TotalAmount,
		Items:         items,
		CreatedAt:     order.CreatedAt,
		TableQrcodeID: order.TableQrcodeID,
		TotalNotes:    order.TotalNotes,
	}

	// 4. 返回结果
	res = &userv1.UserOrderDetailRes{}
	res.Code = 200
	res.Message = "获取订单详情成功"
	res.Data = userOrder

	g.Log().Debug(ctx, "获取订单详情成功，用户ID:", userId, "订单ID:", req.OrderID)

	return res, nil
}
