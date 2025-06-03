// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Categories is the golang structure of table categories for DAO operations like Where/Data.
type Categories struct {
	g.Meta    `orm:"table:categories, do:true"`
	Id        interface{} // 分类ID
	Name      interface{} // 分类名称
	SortOrder interface{} // 排序，数字越小越靠前
	IsActive  interface{} // 是否激活（是否显示）
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
