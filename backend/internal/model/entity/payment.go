// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Payment is the golang structure for table payment.
type Payment struct {
	Id      int64       `json:"id"      orm:"id"       description:"主键"`              // 主键
	OrderId int64       `json:"orderId" orm:"order_id" description:"订单ID"`            // 订单ID
	PayType int         `json:"payType" orm:"pay_type" description:"支付方式（0微信，1支付宝等）"` // 支付方式（0微信，1支付宝等）
	Amount  float64     `json:"amount"  orm:"amount"   description:"支付金额"`            // 支付金额
	PayTime *gtime.Time `json:"payTime" orm:"pay_time" description:"支付时间"`            // 支付时间
	Status  int         `json:"status"  orm:"status"   description:"状态（0未支付，1已支付）"`   // 状态（0未支付，1已支付）
}
