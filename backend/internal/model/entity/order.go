// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Order is the golang structure for table order.
type Order struct {
	Id          int64       `json:"id"          orm:"id"           description:"主键"`                        // 主键
	UserId      int64       `json:"userId"      orm:"user_id"      description:"下单用户ID"`                    // 下单用户ID
	TableId     int64       `json:"tableId"     orm:"table_id"     description:"桌台ID"`                      // 桌台ID
	TotalAmount float64     `json:"totalAmount" orm:"total_amount" description:"总金额"`                       // 总金额
	Status      int         `json:"status"      orm:"status"       description:"订单状态（0待支付，1已支付，2已完成，3已取消）"` // 订单状态（0待支付，1已支付，2已完成，3已取消）
	PayTime     *gtime.Time `json:"payTime"     orm:"pay_time"     description:"支付时间"`                      // 支付时间
	CreateTime  *gtime.Time `json:"createTime"  orm:"create_time"  description:"创建时间"`                      // 创建时间
	UpdateTime  *gtime.Time `json:"updateTime"  orm:"update_time"  description:"更新时间"`                      // 更新时间
}
