package user

import (
	"context"

	userv1 "backend/api/user/v1"
	"backend/internal/service"
	"backend/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// UserOrderList 获取当前登录用户的所有订单列表
func (c *ControllerV1) UserOrderList(ctx context.Context, req *userv1.UserOrderListReq) (res *userv1.UserOrderListRes, err error) {
	// 1. 从上下文中获取用户ID
	userId, err := utility.GetUserID(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取订单列表失败: 未找到用户ID", err)
		return nil, gerror.New("未登录或登录已过期")
	}

	g.Log().Debug(ctx, "获取订单列表，用户ID:", userId, "状态:", req.Status)

	// 2. 调用订单服务获取用户订单列表
	orderService := service.Order()
	orders, total, err := orderService.GetUserOrderList(ctx, &service.UserOrderListReq{
		Status: req.Status,
		Page:   req.Page,
		Limit:  req.Limit,
	}, userId)
	if err != nil {
		g.Log().Error(ctx, "获取订单列表失败:", err, "用户ID:", userId)
		return nil, err
	}

	// 3. 转换为API响应格式
	userOrders := make([]userv1.UserOrder, len(orders))
	for i, order := range orders {
		items := make([]userv1.OrderItem, len(order.Items))
		for j, item := range order.Items {
			items[j] = userv1.OrderItem{
				ProductID:     item.ProductID,
				Name:          item.Name,
				Quantity:      item.Quantity,
				ItemPrice:     item.ItemPrice,
				IsPackageItem: item.IsPackageItem,
				UserPackageID: item.UserPackageID,
				Notes:         item.Notes,
			}
		}

		userOrders[i] = userv1.UserOrder{
			ID:            order.ID,
			OrderSN:       order.OrderSN,
			Status:        order.Status,
			TotalAmount:   order.TotalAmount,
			Items:         items,
			CreatedAt:     order.CreatedAt,
			TableQrcodeID: order.TableQrcodeID,
			TotalNotes:    order.TotalNotes,
		}
	}

	// 4. 返回结果
	res = &userv1.UserOrderListRes{}
	res.Code = 200
	res.Message = "获取订单列表成功"
	res.Data.List = userOrders
	res.Data.Total = total

	g.Log().Debug(ctx, "获取订单列表成功，用户ID:", userId, "订单数量:", len(userOrders))

	return res, nil
}
