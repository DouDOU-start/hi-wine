// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Category is the golang structure for table category.
type Category struct {
	Id         int64       `json:"id"         orm:"id"          description:"主键"`   // 主键
	Name       string      `json:"name"       orm:"name"        description:"分类名称"` // 分类名称
	Sort       int         `json:"sort"       orm:"sort"        description:"排序"`   // 排序
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"` // 创建时间
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"` // 更新时间
}
