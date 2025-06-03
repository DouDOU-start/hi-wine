// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Orders is the golang structure of table orders for DAO operations like Where/Data.
type Orders struct {
	g.Meta        `orm:"table:orders, do:true"`
	Id            interface{} // 订单ID
	OrderSn       interface{} // 订单号，唯一
	UserId        interface{} // 用户ID
	TableQrcodeId interface{} // 关联的桌号二维码ID
	TotalAmount   interface{} // 订单总金额
	PaymentStatus interface{} // 支付状态
	OrderStatus   interface{} // 订单状态
	PaymentMethod interface{} // 支付方式（例如：wechat_pay）
	TransactionId interface{} // 微信支付交易ID
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
	PaidAt        *gtime.Time // 支付时间
}
