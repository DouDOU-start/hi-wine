// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Product is the golang structure for table product.
type Product struct {
	Id          int64       `json:"id"          orm:"id"          description:"主键"`          // 主键
	Name        string      `json:"name"        orm:"name"        description:"商品名称"`        // 商品名称
	CategoryId  int64       `json:"categoryId"  orm:"category_id" description:"分类ID"`        // 分类ID
	Price       float64     `json:"price"       orm:"price"       description:"售价"`          // 售价
	Stock       int         `json:"stock"       orm:"stock"       description:"库存"`          // 库存
	Image       string      `json:"image"       orm:"image"       description:"商品图片"`        // 商品图片
	Status      int         `json:"status"      orm:"status"      description:"状态（0下架，1上架）"` // 状态（0下架，1上架）
	Description string      `json:"description" orm:"description" description:"商品描述"`        // 商品描述
	CreateTime  *gtime.Time `json:"createTime"  orm:"create_time" description:"创建时间"`        // 创建时间
	UpdateTime  *gtime.Time `json:"updateTime"  orm:"update_time" description:"更新时间"`        // 更新时间
}
