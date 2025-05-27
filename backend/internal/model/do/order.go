// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Order is the golang structure of table order for DAO operations like Where/Data.
type Order struct {
	g.Meta      `orm:"table:order, do:true"`
	Id          interface{} // 主键
	UserId      interface{} // 下单用户ID
	TableId     interface{} // 桌台ID
	TotalAmount interface{} // 总金额
	Status      interface{} // 订单状态（0待支付，1已支付，2已完成，3已取消）
	PayTime     *gtime.Time // 支付时间
	CreateTime  *gtime.Time // 创建时间
	UpdateTime  *gtime.Time // 更新时间
}
