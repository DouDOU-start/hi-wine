package service

import (
	"backend/api"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sOrder struct{}

var (
	orderService = sOrder{}
)

// Order 订单服务
func Order() *sOrder {
	return &orderService
}

// Create 创建订单
func (s *sOrder) Create(ctx context.Context, req *api.OrderCreateReq) (res *api.OrderCreateRes, err error) {
	res = &api.OrderCreateRes{}

	// 获取当前用户ID
	userId := ctx.Value("userId")
	if userId == nil {
		return nil, gerror.New("用户未登录")
	}

	// 开始事务
	err = dao.Order.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 创建订单
		order := &entity.Order{
			UserId:      gconv.Int64(userId),
			TableId:     req.TableId,
			TotalAmount: 0,           // 初始金额为0，后续计算
			Status:      0,           // 0待支付
			CreateTime:  gtime.Now(), // 设置创建时间
			UpdateTime:  gtime.Now(), // 设置更新时间
		}

		// 插入订单并获取ID
		orderId, err := dao.Order.Ctx(ctx).TX(tx).InsertAndGetId(order)
		if err != nil {
			return err
		}

		// 2. 处理订单商品
		var totalAmount float64 = 0
		for _, item := range req.Items {
			// 查询商品信息
			product, err := dao.Product.Ctx(ctx).TX(tx).Where(dao.Product.Columns().Id, item.ProductId).One()
			if err != nil {
				return err
			}
			if product.IsEmpty() {
				return gerror.Newf("商品ID=%d不存在", item.ProductId)
			}

			// 从记录中获取商品信息
			var productEntity *entity.Product
			if err = product.Struct(&productEntity); err != nil {
				return err
			}

			// 检查库存
			if productEntity.Stock < item.Quantity {
				return gerror.Newf("商品[%s]库存不足", productEntity.Name)
			}

			// 创建订单商品
			orderItem := &entity.OrderItem{
				OrderId:     orderId,
				ProductId:   item.ProductId,
				ProductName: productEntity.Name,
				Price:       productEntity.Price,
				Quantity:    item.Quantity,
			}

			// 插入订单商品
			_, err = dao.OrderItem.Ctx(ctx).TX(tx).Insert(orderItem)
			if err != nil {
				return err
			}

			// 累加总金额
			totalAmount += orderItem.Price * float64(orderItem.Quantity)

			// 减少库存
			_, err = dao.Product.Ctx(ctx).TX(tx).
				Where(dao.Product.Columns().Id, item.ProductId).
				Decrement(dao.Product.Columns().Stock, item.Quantity)
			if err != nil {
				return err
			}
		}

		// 3. 更新订单总金额
		_, err = dao.Order.Ctx(ctx).TX(tx).
			Where(dao.Order.Columns().Id, orderId).
			Data(g.Map{dao.Order.Columns().TotalAmount: totalAmount}).
			Update()
		if err != nil {
			return err
		}

		// 设置返回信息
		res.Id = orderId

		// 查询完整订单信息
		orderInfo, err := dao.Order.Ctx(ctx).TX(tx).Where(dao.Order.Columns().Id, orderId).One()
		if err != nil {
			return err
		}
		var orderEntity *entity.Order
		if err = orderInfo.Struct(&orderEntity); err != nil {
			return err
		}
		res.Order = orderEntity

		return nil
	})

	return res, err
}

// List 获取订单列表
func (s *sOrder) List(ctx context.Context, req *api.OrderListReq) (res *api.OrderListRes, err error) {
	res = &api.OrderListRes{}

	// 获取当前用户ID
	userId := ctx.Value("userId")
	if userId == nil {
		return nil, gerror.New("用户未登录")
	}

	// 构建查询条件
	m := dao.Order.Ctx(ctx).Where(dao.Order.Columns().UserId, userId)

	// 按状态筛选
	if req.Status >= 0 {
		m = m.Where(dao.Order.Columns().Status, req.Status)
	}

	// 查询总数
	count, err := m.Count()
	if err != nil {
		return nil, err
	}
	res.Total = count

	// 分页查询
	var orders []*entity.Order
	err = m.Page(req.Page, req.Size).Order(dao.Order.Columns().Id + " DESC").Scan(&orders)
	if err != nil {
		return nil, err
	}

	res.List = orders
	return res, nil
}

// Detail 获取订单详情
func (s *sOrder) Detail(ctx context.Context, req *api.OrderDetailReq) (res *api.OrderDetailRes, err error) {
	res = &api.OrderDetailRes{}

	// 获取当前用户ID
	userId := ctx.Value("userId")
	if userId == nil {
		return nil, gerror.New("用户未登录")
	}

	// 查询订单
	order, err := dao.Order.Ctx(ctx).
		Where(dao.Order.Columns().Id, req.Id).
		Where(dao.Order.Columns().UserId, userId).
		One()
	if err != nil {
		return nil, err
	}
	if order.IsEmpty() {
		return nil, gerror.New("订单不存在")
	}

	var orderEntity *entity.Order
	if err = order.Struct(&orderEntity); err != nil {
		return nil, err
	}
	res.Order = orderEntity

	// 查询订单商品
	var orderItems []*entity.OrderItem
	err = dao.OrderItem.Ctx(ctx).
		Where(dao.OrderItem.Columns().OrderId, req.Id).
		Scan(&orderItems)
	if err != nil {
		return nil, err
	}
	res.OrderItems = orderItems

	return res, nil
}

// UpdateStatus 更新订单状态
func (s *sOrder) UpdateStatus(ctx context.Context, req *api.OrderUpdateStatusReq) (res *api.OrderUpdateStatusRes, err error) {
	res = &api.OrderUpdateStatusRes{}

	// 获取当前用户信息
	userId := ctx.Value("userId")
	if userId == nil {
		return nil, gerror.New("用户未登录")
	}

	// 查询订单是否存在
	order, err := dao.Order.Ctx(ctx).
		Where(dao.Order.Columns().Id, req.Id).
		Where(dao.Order.Columns().UserId, userId).
		One()
	if err != nil {
		return nil, err
	}
	if order.IsEmpty() {
		return nil, gerror.New("订单不存在")
	}

	// 更新订单状态
	updateData := g.Map{
		dao.Order.Columns().Status: req.Status,
	}

	// 如果是支付状态，添加支付时间
	if req.Status == 1 {
		updateData[dao.Order.Columns().PayTime] = gtime.Now()
	}

	_, err = dao.Order.Ctx(ctx).
		Where(dao.Order.Columns().Id, req.Id).
		Data(updateData).
		Update()
	if err != nil {
		return nil, err
	}

	res.Status = true
	return res, nil
}
