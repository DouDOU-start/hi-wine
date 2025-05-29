package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 管理员订单列表请求
type AdminOrderListReq struct {
	g.Meta    `path:"/admin/order/list" method:"get" summary:"获取订单列表(管理员)"`
	OrderId   string `json:"orderId" description:"订单号" v:""`
	Username  string `json:"username" description:"用户名" v:""`
	Status    string `json:"status" description:"订单状态" v:""`
	StartDate string `json:"startDate" description:"开始日期" v:""`
	EndDate   string `json:"endDate" description:"结束日期" v:""`
	Page      int    `json:"page" description:"页码" v:"min:1" d:"1"`
	Size      int    `json:"size" description:"每页数量" v:"max:50" d:"10"`
}

// 管理员订单列表响应
type AdminOrderListRes struct {
	List  interface{} `json:"list" description:"订单列表"`
	Total int         `json:"total" description:"总数"`
}

// 管理员订单详情请求
type AdminOrderDetailReq struct {
	g.Meta `path:"/admin/order/detail" method:"get" summary:"获取订单详情(管理员)"`
	Id     int64 `json:"id" description:"订单ID" v:"required"`
}

// 管理员订单详情响应
type AdminOrderDetailRes struct {
	Order interface{} `json:"order" description:"订单信息"`
	Items interface{} `json:"items" description:"订单商品列表"`
}

// 管理员更新订单状态请求
type AdminOrderUpdateStatusReq struct {
	g.Meta `path:"/admin/order/updateStatus" method:"post" summary:"更新订单状态(管理员)"`
	Id     int64 `json:"id" description:"订单ID" v:"required"`
	Status int   `json:"status" description:"订单状态" v:"required|in:0,1,2,3"`
}

// 管理员更新订单状态响应
type AdminOrderUpdateStatusRes struct {
	Success bool `json:"success" description:"是否成功"`
}

// 管理员订单统计请求
type AdminOrderStatsReq struct {
	g.Meta `path:"/admin/order/stats" method:"get" summary:"获取订单统计数据(管理员)"`
}

// 管理员订单统计响应
type AdminOrderStatsRes struct {
	TotalCount     int     `json:"totalCount" description:"订单总数"`
	TodayCount     int     `json:"todayCount" description:"今日订单数"`
	TotalAmount    float64 `json:"totalAmount" description:"总金额"`
	TodayAmount    float64 `json:"todayAmount" description:"今日金额"`
	PendingCount   int     `json:"pendingCount" description:"待支付订单数"`
	CompletedCount int     `json:"completedCount" description:"已完成订单数"`
}

// 管理员导出订单请求
type AdminOrderExportReq struct {
	g.Meta    `path:"/admin/order/export" method:"get" summary:"导出订单数据(管理员)"`
	OrderId   string `json:"orderId" description:"订单号" v:""`
	Username  string `json:"username" description:"用户名" v:""`
	Status    string `json:"status" description:"订单状态" v:""`
	StartDate string `json:"startDate" description:"开始日期" v:""`
	EndDate   string `json:"endDate" description:"结束日期" v:""`
}
