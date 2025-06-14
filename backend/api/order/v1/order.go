package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 创建订单
type CreateOrderReq struct {
	g.Meta        `path:"/orders" method:"post" tags:"订单" summary:"创建订单"`
	TableQrcodeID int64              `json:"table_qrcode_id" description:"桌号二维码ID" v:"required#桌号二维码ID不能为空"`
	Items         []OrderItemRequest `json:"items" description:"订单商品列表" v:"required#订单商品不能为空"`
	TotalNotes    string             `json:"total_notes,omitempty" description:"订单备注"`
}

type CreateOrderRes struct {
	OrderID     int64   `json:"order_id"`
	OrderSN     string  `json:"order_sn"`
	TotalAmount float64 `json:"total_amount"`
	PrepayID    string  `json:"prepay_id"`
}

// 订单商品项
type OrderItemRequest struct {
	ProductID int64  `json:"product_id" description:"商品ID" v:"required#商品ID不能为空"`
	Quantity  int    `json:"quantity" description:"数量" v:"required#数量不能为空"`
	Notes     string `json:"notes,omitempty" description:"商品备注"`
}

// 订单结构体
type Order struct {
	ID            int64       `json:"id"`
	OrderSN       string      `json:"order_sn"`
	Status        string      `json:"status"`
	TotalAmount   float64     `json:"total_amount"`
	Items         []OrderItem `json:"items"`
	CreatedAt     string      `json:"created_at"`
	TableQrcodeID int64       `json:"table_qrcode_id"`
	TotalNotes    string      `json:"total_notes,omitempty"`
}

// 套餐信息
type PackageInfo struct {
	Status    string `json:"status"`     // 套餐状态
	StartTime string `json:"start_time"` // 开始时间
	EndTime   string `json:"end_time"`   // 结束时间
}

type OrderItem struct {
	ProductID     int64        `json:"product_id"`
	Name          string       `json:"name"`
	ImageURL      string       `json:"image_url,omitempty"` // 商品图片地址
	Quantity      int          `json:"quantity"`
	ItemPrice     float64      `json:"item_price"`
	IsPackageItem bool         `json:"is_package_item"`
	UserPackageID int64        `json:"user_package_id,omitempty"`
	Notes         string       `json:"notes,omitempty"`
	PackageInfo   *PackageInfo `json:"package_info,omitempty"` // 套餐信息，仅当IsPackageItem为true时有值
}
