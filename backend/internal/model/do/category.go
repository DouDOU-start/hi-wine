// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Category is the golang structure of table category for DAO operations like Where/Data.
type Category struct {
	g.Meta     `orm:"table:category, do:true"`
	Id         interface{} // 主键
	Name       interface{} // 分类名称
	Sort       interface{} // 排序
	CreateTime *gtime.Time // 创建时间
	UpdateTime *gtime.Time // 更新时间
}
