// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Product is the golang structure of table product for DAO operations like Where/Data.
type Product struct {
	g.Meta      `orm:"table:product, do:true"`
	Id          interface{} // 主键
	Name        interface{} // 商品名称
	CategoryId  interface{} // 分类ID
	Price       interface{} // 售价
	Stock       interface{} // 库存
	Image       interface{} // 商品图片
	Status      interface{} // 状态（0下架，1上架）
	Description interface{} // 商品描述
	CreateTime  *gtime.Time // 创建时间
	UpdateTime  *gtime.Time // 更新时间
}
