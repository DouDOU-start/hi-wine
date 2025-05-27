// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Payment is the golang structure of table payment for DAO operations like Where/Data.
type Payment struct {
	g.Meta  `orm:"table:payment, do:true"`
	Id      interface{} // 主键
	OrderId interface{} // 订单ID
	PayType interface{} // 支付方式（0微信，1支付宝等）
	Amount  interface{} // 支付金额
	PayTime *gtime.Time // 支付时间
	Status  interface{} // 状态（0未支付，1已支付）
}
