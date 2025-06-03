// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Products is the golang structure for table products.
type Products struct {
	Id          int         `json:"id"          orm:"id"          description:"商品ID"`    // 商品ID
	CategoryId  int         `json:"categoryId"  orm:"category_id" description:"所属分类ID"`  // 所属分类ID
	Name        string      `json:"name"        orm:"name"        description:"商品名称"`    // 商品名称
	Description string      `json:"description" orm:"description" description:"商品描述"`    // 商品描述
	Price       float64     `json:"price"       orm:"price"       description:"商品价格"`    // 商品价格
	ImageUrl    string      `json:"imageUrl"    orm:"image_url"   description:"商品图片URL"` // 商品图片URL
	Stock       int         `json:"stock"       orm:"stock"       description:"库存数量"`    // 库存数量
	IsActive    int         `json:"isActive"    orm:"is_active"   description:"是否上架"`    // 是否上架
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:"创建时间"`    // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:"更新时间"`    // 更新时间
}
