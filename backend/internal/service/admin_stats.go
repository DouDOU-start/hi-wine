package service

import (
	"backend/api"
	"backend/internal/dao"
	"context"
	"time"
)

type IAdminStatsService interface {
	Dashboard(ctx context.Context, req *api.AdminDashboardStatsReq) (res *api.AdminDashboardStatsRes, err error)
	Sales(ctx context.Context, req *api.AdminSalesStatsReq) (res *api.AdminSalesStatsRes, err error)
	ProductRanking(ctx context.Context, req *api.AdminProductRankingReq) (res *api.AdminProductRankingRes, err error)
	UserRanking(ctx context.Context, req *api.AdminUserRankingReq) (res *api.AdminUserRankingRes, err error)
	CategorySales(ctx context.Context, req *api.AdminCategorySalesReq) (res *api.AdminCategorySalesRes, err error)
}

var adminStatsService = localAdminStatsService{}

type localAdminStatsService struct{}

func AdminStats() IAdminStatsService {
	return &adminStatsService
}

// 获取仪表盘统计数据
func (s *localAdminStatsService) Dashboard(ctx context.Context, req *api.AdminDashboardStatsReq) (res *api.AdminDashboardStatsRes, err error) {
	res = &api.AdminDashboardStatsRes{}

	// 查询用户总数
	totalUsers, err := dao.User.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}
	res.TotalUsers = totalUsers

	// 查询商品总数
	totalProducts, err := dao.Product.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}
	res.TotalProducts = totalProducts

	// 查询订单总数
	totalOrders, err := dao.Order.Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}
	res.TotalOrders = totalOrders

	// 查询总销售额
	totalSalesValue, err := dao.Order.Ctx(ctx).Fields("SUM(total_amount) as total").Value()
	if err != nil {
		return nil, err
	}
	res.TotalSales = totalSalesValue.Float64()

	// 查询今日订单数
	today := time.Now().Format("2006-01-02")
	todayOrders, err := dao.Order.Ctx(ctx).
		Where("create_time>=?", today+" 00:00:00").
		Where("create_time<=?", today+" 23:59:59").
		Count()
	if err != nil {
		return nil, err
	}
	res.TodayOrders = todayOrders

	// 查询今日销售额
	todaySalesValue, err := dao.Order.Ctx(ctx).
		Where("create_time>=?", today+" 00:00:00").
		Where("create_time<=?", today+" 23:59:59").
		Fields("SUM(total_amount) as total").
		Value()
	if err != nil {
		return nil, err
	}
	res.TodaySales = todaySalesValue.Float64()

	// 查询待处理订单数
	pendingOrders, err := dao.Order.Ctx(ctx).Where("status=0").Count()
	if err != nil {
		return nil, err
	}
	res.PendingOrders = pendingOrders

	// 查询已完成订单数
	completedOrders, err := dao.Order.Ctx(ctx).Where("status=2").Count()
	if err != nil {
		return nil, err
	}
	res.CompletedOrders = completedOrders

	// 查询缺货商品数
	productsOutStock, err := dao.Product.Ctx(ctx).Where("stock=0").Count()
	if err != nil {
		return nil, err
	}
	res.ProductsOutStock = productsOutStock

	return res, nil
}

// 获取销售统计数据
func (s *localAdminStatsService) Sales(ctx context.Context, req *api.AdminSalesStatsReq) (res *api.AdminSalesStatsRes, err error) {
	res = &api.AdminSalesStatsRes{}

	// 处理日期范围
	startDate := req.StartDate
	endDate := req.EndDate

	if startDate == "" {
		// 默认查询最近30天
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// 生成日期列表
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var dates []string
	var sales []float64
	var orders []int

	totalSales := 0.0
	totalOrders := 0

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dates = append(dates, dateStr)

		// 查询当天销售额
		salesValue, err := dao.Order.Ctx(ctx).
			Where("create_time>=?", dateStr+" 00:00:00").
			Where("create_time<=?", dateStr+" 23:59:59").
			Fields("SUM(total_amount) as total").
			Value()
		if err != nil {
			return nil, err
		}

		saleAmount := salesValue.Float64()
		sales = append(sales, saleAmount)
		totalSales += saleAmount

		// 查询当天订单数
		orderCount, err := dao.Order.Ctx(ctx).
			Where("create_time>=?", dateStr+" 00:00:00").
			Where("create_time<=?", dateStr+" 23:59:59").
			Count()
		if err != nil {
			return nil, err
		}

		orders = append(orders, orderCount)
		totalOrders += orderCount
	}

	res.Dates = dates
	res.Sales = sales
	res.Orders = orders

	// 计算汇总数据
	days := len(dates)
	res.Summary.TotalSales = totalSales
	res.Summary.TotalOrders = totalOrders

	if days > 0 {
		res.Summary.AvgSales = totalSales / float64(days)
		res.Summary.AvgOrders = float64(totalOrders) / float64(days)
	}

	return res, nil
}

// 获取商品销量排行
func (s *localAdminStatsService) ProductRanking(ctx context.Context, req *api.AdminProductRankingReq) (res *api.AdminProductRankingRes, err error) {
	res = &api.AdminProductRankingRes{}

	// 处理日期范围
	startDate := req.StartDate
	endDate := req.EndDate

	if startDate == "" {
		// 默认查询最近30天
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// 构建查询条件
	model := dao.OrderItem.Ctx(ctx).
		Fields("product_id, product_name, SUM(quantity) as sales, SUM(price * quantity) as amount").
		InnerJoin("order", "order_item.order_id = order.id")

	if startDate != "" {
		model = model.Where("order.create_time>=?", startDate+" 00:00:00")
	}

	if endDate != "" {
		model = model.Where("order.create_time<=?", endDate+" 23:59:59")
	}

	// 分组查询
	records, err := model.Group("product_id").Order("sales DESC").Limit(req.Limit).All()
	if err != nil {
		return nil, err
	}

	// 构建结果
	list := make([]struct {
		ProductId   int64   `json:"productId" description:"商品ID"`
		ProductName string  `json:"productName" description:"商品名称"`
		Sales       int     `json:"sales" description:"销量"`
		Amount      float64 `json:"amount" description:"销售额"`
	}, 0)

	for _, record := range records {
		item := struct {
			ProductId   int64   `json:"productId" description:"商品ID"`
			ProductName string  `json:"productName" description:"商品名称"`
			Sales       int     `json:"sales" description:"销量"`
			Amount      float64 `json:"amount" description:"销售额"`
		}{
			ProductId:   record["product_id"].Int64(),
			ProductName: record["product_name"].String(),
			Sales:       record["sales"].Int(),
			Amount:      record["amount"].Float64(),
		}
		list = append(list, item)
	}

	res.List = list

	return res, nil
}

// 获取用户消费排行
func (s *localAdminStatsService) UserRanking(ctx context.Context, req *api.AdminUserRankingReq) (res *api.AdminUserRankingRes, err error) {
	res = &api.AdminUserRankingRes{}

	// 处理日期范围
	startDate := req.StartDate
	endDate := req.EndDate

	if startDate == "" {
		// 默认查询最近30天
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// 构建查询条件
	model := dao.Order.Ctx(ctx).
		Fields("user_id, COUNT(id) as orders, SUM(total_amount) as amount").
		InnerJoin("user", "order.user_id = user.id")

	if startDate != "" {
		model = model.Where("order.create_time>=?", startDate+" 00:00:00")
	}

	if endDate != "" {
		model = model.Where("order.create_time<=?", endDate+" 23:59:59")
	}

	// 分组查询
	records, err := model.Group("user_id").Order("amount DESC").Limit(req.Limit).All()
	if err != nil {
		return nil, err
	}

	// 构建结果
	list := make([]struct {
		UserId   int64   `json:"userId" description:"用户ID"`
		Username string  `json:"username" description:"用户名"`
		Orders   int     `json:"orders" description:"订单数"`
		Amount   float64 `json:"amount" description:"消费金额"`
	}, 0)

	for _, record := range records {
		userId := record["user_id"].Int64()

		// 查询用户信息
		user, err := dao.User.Ctx(ctx).Where("id=?", userId).One()
		if err != nil {
			return nil, err
		}

		item := struct {
			UserId   int64   `json:"userId" description:"用户ID"`
			Username string  `json:"username" description:"用户名"`
			Orders   int     `json:"orders" description:"订单数"`
			Amount   float64 `json:"amount" description:"消费金额"`
		}{
			UserId:   userId,
			Username: user["nickname"].String(),
			Orders:   record["orders"].Int(),
			Amount:   record["amount"].Float64(),
		}
		list = append(list, item)
	}

	res.List = list

	return res, nil
}

// 获取分类销售统计
func (s *localAdminStatsService) CategorySales(ctx context.Context, req *api.AdminCategorySalesReq) (res *api.AdminCategorySalesRes, err error) {
	res = &api.AdminCategorySalesRes{}

	// 处理日期范围
	startDate := req.StartDate
	endDate := req.EndDate

	if startDate == "" {
		// 默认查询最近30天
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// 构建查询条件
	model := dao.OrderItem.Ctx(ctx).
		Fields("product.category_id, category.name as category_name, SUM(order_item.quantity) as sales, SUM(order_item.price * order_item.quantity) as amount").
		InnerJoin("order", "order_item.order_id = order.id").
		InnerJoin("product", "order_item.product_id = product.id").
		InnerJoin("category", "product.category_id = category.id")

	if startDate != "" {
		model = model.Where("order.create_time>=?", startDate+" 00:00:00")
	}

	if endDate != "" {
		model = model.Where("order.create_time<=?", endDate+" 23:59:59")
	}

	// 分组查询
	records, err := model.Group("product.category_id").Order("amount DESC").All()
	if err != nil {
		return nil, err
	}

	// 计算总金额
	var totalAmount float64
	for _, record := range records {
		totalAmount += record["amount"].Float64()
	}

	// 构建结果
	list := make([]struct {
		CategoryId   int64   `json:"categoryId" description:"分类ID"`
		CategoryName string  `json:"categoryName" description:"分类名称"`
		Sales        int     `json:"sales" description:"销量"`
		Amount       float64 `json:"amount" description:"销售额"`
		Percentage   float64 `json:"percentage" description:"占比"`
	}, 0)

	for _, record := range records {
		amount := record["amount"].Float64()
		percentage := 0.0
		if totalAmount > 0 {
			percentage = amount / totalAmount * 100
		}

		item := struct {
			CategoryId   int64   `json:"categoryId" description:"分类ID"`
			CategoryName string  `json:"categoryName" description:"分类名称"`
			Sales        int     `json:"sales" description:"销量"`
			Amount       float64 `json:"amount" description:"销售额"`
			Percentage   float64 `json:"percentage" description:"占比"`
		}{
			CategoryId:   record["category_id"].Int64(),
			CategoryName: record["category_name"].String(),
			Sales:        record["sales"].Int(),
			Amount:       amount,
			Percentage:   percentage,
		}
		list = append(list, item)
	}

	res.List = list

	return res, nil
}
