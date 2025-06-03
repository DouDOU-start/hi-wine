// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Products is the golang structure of table products for DAO operations like Where/Data.
type Products struct {
	g.Meta      `orm:"table:products, do:true"`
	Id          interface{} // 商品ID
	CategoryId  interface{} // 所属分类ID
	Name        interface{} // 商品名称
	Description interface{} // 商品描述
	Price       interface{} // 商品价格
	ImageUrl    interface{} // 商品图片URL
	Stock       interface{} // 库存数量
	IsActive    interface{} // 是否上架
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
