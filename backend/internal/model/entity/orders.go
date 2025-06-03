// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Orders is the golang structure for table orders.
type Orders struct {
	Id            int         `json:"id"            orm:"id"              description:"订单ID"`                // 订单ID
	OrderSn       string      `json:"orderSn"       orm:"order_sn"        description:"订单号，唯一"`              // 订单号，唯一
	UserId        int         `json:"userId"        orm:"user_id"         description:"用户ID"`                // 用户ID
	TableQrcodeId int         `json:"tableQrcodeId" orm:"table_qrcode_id" description:"关联的桌号二维码ID"`          // 关联的桌号二维码ID
	TotalAmount   float64     `json:"totalAmount"   orm:"total_amount"    description:"订单总金额"`               // 订单总金额
	PaymentStatus string      `json:"paymentStatus" orm:"payment_status"  description:"支付状态"`                // 支付状态
	OrderStatus   string      `json:"orderStatus"   orm:"order_status"    description:"订单状态"`                // 订单状态
	PaymentMethod string      `json:"paymentMethod" orm:"payment_method"  description:"支付方式（例如：wechat_pay）"` // 支付方式（例如：wechat_pay）
	TransactionId string      `json:"transactionId" orm:"transaction_id"  description:"微信支付交易ID"`            // 微信支付交易ID
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:"创建时间"`                // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      description:"更新时间"`                // 更新时间
	PaidAt        *gtime.Time `json:"paidAt"        orm:"paid_at"         description:"支付时间"`                // 支付时间
}
