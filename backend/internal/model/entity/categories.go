// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Categories is the golang structure for table categories.
type Categories struct {
	Id        int         `json:"id"        orm:"id"         description:"分类ID"`       // 分类ID
	Name      string      `json:"name"      orm:"name"       description:"分类名称"`       // 分类名称
	SortOrder int         `json:"sortOrder" orm:"sort_order" description:"排序，数字越小越靠前"` // 排序，数字越小越靠前
	IsActive  int         `json:"isActive"  orm:"is_active"  description:"是否激活（是否显示）"` // 是否激活（是否显示）
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`       // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`       // 更新时间
}
