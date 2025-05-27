// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Id         int64       `json:"id"         orm:"id"          description:"主键"`               // 主键
	Openid     string      `json:"openid"     orm:"openid"      description:"微信openid"`         // 微信openid
	Nickname   string      `json:"nickname"   orm:"nickname"    description:"昵称"`               // 昵称
	Avatar     string      `json:"avatar"     orm:"avatar"      description:"头像"`               // 头像
	Phone      string      `json:"phone"      orm:"phone"       description:"手机号"`              // 手机号
	Role       int         `json:"role"       orm:"role"        description:"角色（0顾客，1店员，2管理员）"` // 角色（0顾客，1店员，2管理员）
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`             // 创建时间
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`             // 更新时间
}
