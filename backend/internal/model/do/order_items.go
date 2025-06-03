// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderItems is the golang structure of table order_items for DAO operations like Where/Data.
type OrderItems struct {
	g.Meta        `orm:"table:order_items, do:true"`
	Id            interface{} // 订单项ID
	OrderId       interface{} // 订单ID
	ProductId     interface{} // 商品ID
	ProductName   interface{} // 冗余商品名称
	Price         interface{} // 下单时商品单价
	Quantity      interface{} // 购买数量
	Subtotal      interface{} // 小计
	IsPackageItem interface{} // 是否为畅饮套餐内商品
	UserPackageId interface{} // 关联的用户套餐购买记录ID（如果为套餐商品）
	ItemPrice     interface{} // 该订单项的实际结算价格（畅饮则为0）
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
