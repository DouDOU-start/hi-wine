// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Table is the golang structure of table table for DAO operations like Where/Data.
type Table struct {
	g.Meta     `orm:"table:table, do:true"`
	Id         interface{} // 主键
	Name       interface{} // 桌台名称/编号
	Status     interface{} // 状态（0空闲，1占用）
	CreateTime *gtime.Time // 创建时间
	UpdateTime *gtime.Time // 更新时间
}
