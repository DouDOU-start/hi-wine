// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Table is the golang structure for table table.
type Table struct {
	Id         int64       `json:"id"         orm:"id"          description:"主键"`          // 主键
	Name       string      `json:"name"       orm:"name"        description:"桌台名称/编号"`     // 桌台名称/编号
	Status     int         `json:"status"     orm:"status"      description:"状态（0空闲，1占用）"` // 状态（0空闲，1占用）
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`        // 创建时间
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`        // 更新时间
}
