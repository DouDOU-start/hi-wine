package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 管理员仪表盘统计请求
type AdminDashboardStatsReq struct {
	g.Meta `path:"/admin/stats/dashboard" method:"get" summary:"获取仪表盘统计数据(管理员)"`
}

// 管理员仪表盘统计响应
type AdminDashboardStatsRes struct {
	TotalUsers       int     `json:"totalUsers" description:"用户总数"`
	TotalProducts    int     `json:"totalProducts" description:"商品总数"`
	TotalOrders      int     `json:"totalOrders" description:"订单总数"`
	TotalSales       float64 `json:"totalSales" description:"总销售额"`
	TodayOrders      int     `json:"todayOrders" description:"今日订单数"`
	TodaySales       float64 `json:"todaySales" description:"今日销售额"`
	PendingOrders    int     `json:"pendingOrders" description:"待处理订单数"`
	CompletedOrders  int     `json:"completedOrders" description:"已完成订单数"`
	ProductsOutStock int     `json:"productsOutStock" description:"缺货商品数"`
}

// 管理员销售统计请求
type AdminSalesStatsReq struct {
	g.Meta    `path:"/admin/stats/sales" method:"get" summary:"获取销售统计数据(管理员)"`
	StartDate string `json:"startDate" description:"开始日期" v:""`
	EndDate   string `json:"endDate" description:"结束日期" v:""`
}

// 管理员销售统计响应
type AdminSalesStatsRes struct {
	Dates   []string  `json:"dates" description:"日期列表"`
	Sales   []float64 `json:"sales" description:"销售额列表"`
	Orders  []int     `json:"orders" description:"订单数列表"`
	Summary struct {
		TotalSales  float64 `json:"totalSales" description:"总销售额"`
		TotalOrders int     `json:"totalOrders" description:"总订单数"`
		AvgSales    float64 `json:"avgSales" description:"平均日销售额"`
		AvgOrders   float64 `json:"avgOrders" description:"平均日订单数"`
	} `json:"summary" description:"汇总数据"`
}

// 管理员商品销量排行请求
type AdminProductRankingReq struct {
	g.Meta    `path:"/admin/stats/products/ranking" method:"get" summary:"获取商品销量排行(管理员)"`
	StartDate string `json:"startDate" description:"开始日期" v:""`
	EndDate   string `json:"endDate" description:"结束日期" v:""`
	Limit     int    `json:"limit" description:"返回数量" v:"max:50" d:"10"`
}

// 管理员商品销量排行响应
type AdminProductRankingRes struct {
	List []struct {
		ProductId   int64   `json:"productId" description:"商品ID"`
		ProductName string  `json:"productName" description:"商品名称"`
		Sales       int     `json:"sales" description:"销量"`
		Amount      float64 `json:"amount" description:"销售额"`
	} `json:"list" description:"商品销量排行列表"`
}

// 管理员用户消费排行请求
type AdminUserRankingReq struct {
	g.Meta    `path:"/admin/stats/users/ranking" method:"get" summary:"获取用户消费排行(管理员)"`
	StartDate string `json:"startDate" description:"开始日期" v:""`
	EndDate   string `json:"endDate" description:"结束日期" v:""`
	Limit     int    `json:"limit" description:"返回数量" v:"max:50" d:"10"`
}

// 管理员用户消费排行响应
type AdminUserRankingRes struct {
	List []struct {
		UserId   int64   `json:"userId" description:"用户ID"`
		Username string  `json:"username" description:"用户名"`
		Orders   int     `json:"orders" description:"订单数"`
		Amount   float64 `json:"amount" description:"消费金额"`
	} `json:"list" description:"用户消费排行列表"`
}

// 管理员分类销售统计请求
type AdminCategorySalesReq struct {
	g.Meta    `path:"/admin/stats/category/sales" method:"get" summary:"获取分类销售统计(管理员)"`
	StartDate string `json:"startDate" description:"开始日期" v:""`
	EndDate   string `json:"endDate" description:"结束日期" v:""`
}

// 管理员分类销售统计响应
type AdminCategorySalesRes struct {
	List []struct {
		CategoryId   int64   `json:"categoryId" description:"分类ID"`
		CategoryName string  `json:"categoryName" description:"分类名称"`
		Sales        int     `json:"sales" description:"销量"`
		Amount       float64 `json:"amount" description:"销售额"`
		Percentage   float64 `json:"percentage" description:"占比"`
	} `json:"list" description:"分类销售统计列表"`
}
