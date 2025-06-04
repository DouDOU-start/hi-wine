package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 用户端-订单分组

type UserOrder struct {
	ID            int64       `json:"id"`
	OrderSN       string      `json:"order_sn"`
	Status        string      `json:"status"`
	TotalAmount   float64     `json:"total_amount"`
	Items         []OrderItem `json:"items"`
	CreatedAt     string      `json:"created_at"`
	TableQrcodeID int64       `json:"table_qrcode_id"`
	TotalNotes    string      `json:"total_notes,omitempty"`
}

type OrderItem struct {
	ProductID     int64   `json:"product_id"`
	Name          string  `json:"name"`
	Quantity      int     `json:"quantity"`
	ItemPrice     float64 `json:"item_price"`
	IsPackageItem bool    `json:"is_package_item"`
	UserPackageID int64   `json:"user_package_id,omitempty"`
	Notes         string  `json:"notes,omitempty"`
}

// 获取用户订单列表
type UserOrderListReq struct {
	g.Meta `path:"/user/orders" method:"get" tags:"订单" summary:"获取当前登录用户的所有订单列表"`
	Status string `json:"status" in:"query" description:"订单状态筛选"`
	Page   int    `json:"page" in:"query" description:"页码，默认1"`
	Limit  int    `json:"limit" in:"query" description:"每页数量，默认10"`
}
type UserOrderListRes struct {
	common.Response[struct {
		List  []UserOrder `json:"list"`
		Total int         `json:"total"`
	}] `json:",inline"`
}

// 获取订单详情
type UserOrderDetailReq struct {
	g.Meta  `path:"/orders/{order_id}" method:"get" tags:"订单" summary:"获取指定订单的详细信息"`
	OrderID int64 `json:"order_id" in:"path" v:"required#订单ID必填"`
}
type UserOrderDetailRes struct {
	common.Response[UserOrder] `json:",inline"`
}
