package service

import (
	"backend/api"
	"backend/internal/dao"
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type IAdminOrderService interface {
	List(ctx context.Context, req *api.AdminOrderListReq) (res *api.AdminOrderListRes, err error)
	Detail(ctx context.Context, req *api.AdminOrderDetailReq) (res *api.AdminOrderDetailRes, err error)
	UpdateStatus(ctx context.Context, req *api.AdminOrderUpdateStatusReq) (res *api.AdminOrderUpdateStatusRes, err error)
	Stats(ctx context.Context, req *api.AdminOrderStatsReq) (res *api.AdminOrderStatsRes, err error)
	Export(ctx context.Context, req *api.AdminOrderExportReq) (res interface{}, err error)
}

var adminOrderService = localAdminOrderService{}

type localAdminOrderService struct{}

func AdminOrder() IAdminOrderService {
	return &adminOrderService
}

// 获取订单列表
func (s *localAdminOrderService) List(ctx context.Context, req *api.AdminOrderListReq) (res *api.AdminOrderListRes, err error) {
	res = &api.AdminOrderListRes{}

	// 构建查询条件
	model := dao.Order.Ctx(ctx)

	if req.OrderId != "" {
		model = model.Where("id=?", req.OrderId)
	}

	if req.Username != "" {
		// 联表查询
		model = model.Where("user_id IN (SELECT id FROM user WHERE nickname LIKE ?)", "%"+req.Username+"%")
	}

	if req.Status != "" {
		model = model.Where("status=?", req.Status)
	}

	if req.StartDate != "" {
		model = model.Where("create_time>=?", req.StartDate+" 00:00:00")
	}

	if req.EndDate != "" {
		model = model.Where("create_time<=?", req.EndDate+" 23:59:59")
	}

	// 分页查询
	count, err := model.Count()
	if err != nil {
		return nil, err
	}

	list, err := model.Page(req.Page, req.Size).Order("id DESC").All()
	if err != nil {
		return nil, err
	}

	// 转换为普通的map数组，以便添加用户信息
	orderList := make([]map[string]interface{}, len(list))
	for i, v := range list {
		orderList[i] = gconv.Map(v)
	}

	// 查询用户信息
	for i, order := range orderList {
		userId := gconv.Int64(order["user_id"])
		user, err := dao.User.Ctx(ctx).Where("id=?", userId).One()
		if err != nil {
			return nil, err
		}

		// 添加用户信息
		orderList[i]["user"] = gconv.Map(user)
	}

	res.Total = count
	res.List = orderList

	return res, nil
}

// 获取订单详情
func (s *localAdminOrderService) Detail(ctx context.Context, req *api.AdminOrderDetailReq) (res *api.AdminOrderDetailRes, err error) {
	res = &api.AdminOrderDetailRes{}

	// 查询订单信息
	order, err := dao.Order.Ctx(ctx).Where("id=?", req.Id).One()
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, gerror.New("订单不存在")
	}

	// 转换为普通map
	orderMap := gconv.Map(order)

	// 查询用户信息
	userId := gconv.Int64(orderMap["user_id"])
	user, err := dao.User.Ctx(ctx).Where("id=?", userId).One()
	if err != nil {
		return nil, err
	}

	// 添加用户信息
	orderMap["user"] = gconv.Map(user)

	// 查询订单商品
	items, err := dao.OrderItem.Ctx(ctx).Where("order_id=?", req.Id).All()
	if err != nil {
		return nil, err
	}

	res.Order = orderMap
	res.Items = items

	return res, nil
}

// 更新订单状态
func (s *localAdminOrderService) UpdateStatus(ctx context.Context, req *api.AdminOrderUpdateStatusReq) (res *api.AdminOrderUpdateStatusRes, err error) {
	res = &api.AdminOrderUpdateStatusRes{}

	// 查询订单是否存在
	order, err := dao.Order.Ctx(ctx).Where("id=?", req.Id).One()
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, gerror.New("订单不存在")
	}

	// 更新订单状态
	data := g.Map{
		"status": req.Status,
	}

	// 如果是支付状态，添加支付时间
	if req.Status == 1 {
		data["pay_time"] = gtime.Now()
	}

	_, err = dao.Order.Ctx(ctx).Data(data).Where("id=?", req.Id).Update()
	if err != nil {
		return nil, err
	}

	res.Success = true

	return res, nil
}

// 获取订单统计
func (s *localAdminOrderService) Stats(ctx context.Context, req *api.AdminOrderStatsReq) (res *api.AdminOrderStatsRes, err error) {
	res = &api.AdminOrderStatsRes{}

	// 查询订单总数
	totalCount, err := dao.Order.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}
	res.TotalCount = totalCount

	// 查询今日订单数
	today := time.Now().Format("2006-01-02")
	todayCount, err := dao.Order.Ctx(ctx).
		Where("create_time>=?", today+" 00:00:00").
		Where("create_time<=?", today+" 23:59:59").
		Count()
	if err != nil {
		return nil, err
	}
	res.TodayCount = todayCount

	// 查询总金额
	totalAmountValue, err := dao.Order.Ctx(ctx).Fields("SUM(total_amount) as total").Value()
	if err != nil {
		return nil, err
	}
	res.TotalAmount = totalAmountValue.Float64()

	// 查询今日金额
	todayAmountValue, err := dao.Order.Ctx(ctx).
		Where("create_time>=?", today+" 00:00:00").
		Where("create_time<=?", today+" 23:59:59").
		Fields("SUM(total_amount) as total").
		Value()
	if err != nil {
		return nil, err
	}
	res.TodayAmount = todayAmountValue.Float64()

	// 查询待支付订单数
	pendingCount, err := dao.Order.Ctx(ctx).Where("status=0").Count()
	if err != nil {
		return nil, err
	}
	res.PendingCount = pendingCount

	// 查询已完成订单数
	completedCount, err := dao.Order.Ctx(ctx).Where("status=2").Count()
	if err != nil {
		return nil, err
	}
	res.CompletedCount = completedCount

	return res, nil
}

// 导出订单数据
func (s *localAdminOrderService) Export(ctx context.Context, req *api.AdminOrderExportReq) (res interface{}, err error) {
	// 构建查询条件
	model := dao.Order.Ctx(ctx)

	if req.OrderId != "" {
		model = model.Where("id=?", req.OrderId)
	}

	if req.Username != "" {
		// 联表查询
		model = model.Where("user_id IN (SELECT id FROM user WHERE nickname LIKE ?)", "%"+req.Username+"%")
	}

	if req.Status != "" {
		model = model.Where("status=?", req.Status)
	}

	if req.StartDate != "" {
		model = model.Where("create_time>=?", req.StartDate+" 00:00:00")
	}

	if req.EndDate != "" {
		model = model.Where("create_time<=?", req.EndDate+" 23:59:59")
	}

	// 查询所有符合条件的订单
	list, err := model.Order("id DESC").All()
	if err != nil {
		return nil, err
	}

	// 转换为普通的map数组
	orderList := make([]map[string]interface{}, len(list))
	for i, v := range list {
		orderList[i] = gconv.Map(v)
	}

	// 查询用户信息和订单商品
	for i, order := range orderList {
		// 查询用户信息
		userId := gconv.Int64(order["user_id"])
		user, err := dao.User.Ctx(ctx).Where("id=?", userId).One()
		if err != nil {
			return nil, err
		}
		orderList[i]["user"] = gconv.Map(user)

		// 查询订单商品
		orderId := order["id"]
		items, err := dao.OrderItem.Ctx(ctx).Where("order_id=?", orderId).All()
		if err != nil {
			return nil, err
		}

		// 转换订单商品为map数组
		itemsList := make([]map[string]interface{}, len(items))
		for j, item := range items {
			itemsList[j] = gconv.Map(item)
		}

		orderList[i]["items"] = itemsList
	}

	return orderList, nil
}
