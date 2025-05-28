package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 创建订单请求
type OrderCreateReq struct {
	g.Meta  `path:"/order/create" method:"post" summary:"创建订单"`
	TableId int64          `json:"tableId" description:"桌台ID" v:""`
	Items   []OrderItemReq `json:"items" description:"订单商品列表" v:"required"`
}

// 订单商品请求
type OrderItemReq struct {
	ProductId int64 `json:"productId" description:"商品ID" v:"required"`
	Quantity  int   `json:"quantity" description:"数量" v:"required|min:1"`
}

// 创建订单响应
type OrderCreateRes struct {
	Id    int64       `json:"id" description:"订单ID"`
	Order interface{} `json:"order" description:"订单信息"`
}

// 订单列表请求
type OrderListReq struct {
	g.Meta `path:"/order/list" method:"get" summary:"获取订单列表"`
	Status int `json:"status" description:"订单状态" v:""`
	Page   int `json:"page" description:"页码" v:"min:1" d:"1"`
	Size   int `json:"size" description:"每页数量" v:"max:50" d:"10"`
}

// 订单列表响应
type OrderListRes struct {
	List  interface{} `json:"list" description:"订单列表"`
	Total int         `json:"total" description:"总数"`
}

// 订单详情请求
type OrderDetailReq struct {
	g.Meta `path:"/order/detail" method:"get" summary:"获取订单详情"`
	Id     int64 `json:"id" description:"订单ID" v:"required"`
}

// 订单详情响应
type OrderDetailRes struct {
	Order      interface{} `json:"order" description:"订单信息"`
	OrderItems interface{} `json:"orderItems" description:"订单商品列表"`
}

// 更新订单状态请求
type OrderUpdateStatusReq struct {
	g.Meta `path:"/order/updateStatus" method:"post" summary:"更新订单状态"`
	Id     int64 `json:"id" description:"订单ID" v:"required"`
	Status int   `json:"status" description:"订单状态" v:"required|in:0,1,2,3"`
}

// 更新订单状态响应
type OrderUpdateStatusRes struct {
	Status bool `json:"status" description:"更新结果"`
}
