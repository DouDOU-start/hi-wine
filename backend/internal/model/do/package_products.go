// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PackageProducts is the golang structure of table package_products for DAO operations like Where/Data.
type PackageProducts struct {
	g.Meta    `orm:"table:package_products, do:true"`
	PackageId interface{} // 关联畅饮套餐ID
	ProductId interface{} // 关联商品ID
	CreatedAt *gtime.Time // 创建时间
}
