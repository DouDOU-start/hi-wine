// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DrinkAllYouCanPackages is the golang structure for table drink_all_you_can_packages.
type DrinkAllYouCanPackages struct {
	Id              int         `json:"id"              orm:"id"               description:"套餐ID"`              // 套餐ID
	Name            string      `json:"name"            orm:"name"             description:"套餐名称"`              // 套餐名称
	Description     string      `json:"description"     orm:"description"      description:"套餐描述"`              // 套餐描述
	Price           float64     `json:"price"           orm:"price"            description:"套餐价格"`              // 套餐价格
	DurationMinutes int         `json:"durationMinutes" orm:"duration_minutes" description:"有效时长（分钟），0表示无时间限制"` // 有效时长（分钟），0表示无时间限制
	IsActive        int         `json:"isActive"        orm:"is_active"        description:"是否激活（是否可售）"`        // 是否激活（是否可售）
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:"创建时间"`              // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       description:"更新时间"`              // 更新时间
}
