// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// OrderItem is the golang structure of table order_item for DAO operations like Where/Data.
type OrderItem struct {
	g.Meta      `orm:"table:order_item, do:true"`
	Id          interface{} // 主键
	OrderId     interface{} // 订单ID
	ProductId   interface{} // 商品ID
	ProductName interface{} // 商品名称
	Price       interface{} // 商品单价
	Quantity    interface{} // 数量
}
