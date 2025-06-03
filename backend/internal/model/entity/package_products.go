// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PackageProducts is the golang structure for table package_products.
type PackageProducts struct {
	PackageId int         `json:"packageId" orm:"package_id" description:"关联畅饮套餐ID"` // 关联畅饮套餐ID
	ProductId int         `json:"productId" orm:"product_id" description:"关联商品ID"`   // 关联商品ID
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`     // 创建时间
}
