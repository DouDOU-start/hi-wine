// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DrinkAllYouCanPackages is the golang structure of table drink_all_you_can_packages for DAO operations like Where/Data.
type DrinkAllYouCanPackages struct {
	g.Meta          `orm:"table:drink_all_you_can_packages, do:true"`
	Id              interface{} // 套餐ID
	Name            interface{} // 套餐名称
	Description     interface{} // 套餐描述
	Price           interface{} // 套餐价格
	DurationMinutes interface{} // 有效时长（分钟），0表示无时间限制
	IsActive        interface{} // 是否激活（是否可售）
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}
