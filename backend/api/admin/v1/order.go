package v1

import (
	"backend/api/common"
	orderv1 "backend/api/order/v1"

	"github.com/gogf/gf/v2/frame/g"
)

// 管理端-订单分组

type AdminOrder struct {
	ID            int64               `json:"id"`
	OrderSN       string              `json:"order_sn"`
	UserID        int64               `json:"user_id"`
	UserName      string              `json:"user_name,omitempty"`     // 用户名称
	UserNickname  string              `json:"user_nickname,omitempty"` // 用户昵称
	UserPhone     string              `json:"user_phone,omitempty"`    // 用户手机号
	TableQrcodeID int64               `json:"table_qrcode_id"`
	TableNumber   string              `json:"table_number,omitempty"` // 桌号
	TotalAmount   float64             `json:"total_amount"`
	ItemCount     int                 `json:"item_count"` // 商品总数量
	PaymentStatus string              `json:"payment_status"`
	OrderStatus   string              `json:"order_status"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
	PaidAt        string              `json:"paid_at,omitempty"`     // 支付时间
	Items         []orderv1.OrderItem `json:"items,omitempty"`       // 订单项列表
	TotalNotes    string              `json:"total_notes,omitempty"` // 订单备注
}

// 获取订单列表（分页、筛选、模糊搜索）
type AdminOrderListReq struct {
	g.Meta      `path:"orders" method:"get" tags:"管理端-订单" summary:"获取订单列表（分页、筛选、模糊搜索）"`
	Status      string `json:"status" in:"query" description:"订单状态筛选"`
	Page        int    `json:"page" in:"query" description:"页码，默认1"`
	Limit       int    `json:"limit" in:"query" description:"每页数量，默认10"`
	OrderSN     string `json:"order_sn" in:"query" description:"订单号搜索"`
	StartDate   string `json:"start_date" in:"query" description:"起始日期(YYYY-MM-DD)"`
	EndDate     string `json:"end_date" in:"query" description:"结束日期(YYYY-MM-DD)"`
	TableNumber string `json:"table_number" in:"query" description:"桌号搜索"`
	UserID      int64  `json:"user_id" in:"query" description:"用户ID搜索"`
}
type AdminOrderListRes struct {
	common.Response[struct {
		List  []AdminOrder `json:"list"`
		Total int          `json:"total"`
	}] `json:",inline"`
}

// 获取订单详情
type AdminOrderDetailReq struct {
	g.Meta  `path:"orders/{order_id}" method:"get" tags:"管理端-订单" summary:"获取订单详情"`
	OrderID int64 `json:"order_id" in:"path" v:"required#订单ID必填"`
}
type AdminOrderDetailRes struct {
	common.Response[AdminOrder] `json:",inline"`
}

// 更新订单状态
type AdminOrderUpdateStatusReq struct {
	g.Meta  `path:"orders/{order_id}/status" method:"put" tags:"管理端-订单" summary:"更新订单状态"`
	OrderID int64  `json:"order_id" in:"path" v:"required#订单ID必填"`
	Status  string `json:"status" v:"required#新状态必填"`
	Reason  string `json:"reason" description:"变更原因"`
}
type AdminOrderUpdateStatusRes struct {
	common.Response[AdminOrder] `json:",inline"`
}
