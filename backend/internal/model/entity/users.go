// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id        int         `json:"id"        orm:"id"         description:"用户ID"`     // 用户ID
	Openid    string      `json:"openid"    orm:"openid"     description:"微信openid"` // 微信openid
	Nickname  string      `json:"nickname"  orm:"nickname"   description:"微信昵称"`     // 微信昵称
	AvatarUrl string      `json:"avatarUrl" orm:"avatar_url" description:"微信头像URL"`  // 微信头像URL
	Phone     string      `json:"phone"     orm:"phone"      description:"手机号"`      // 手机号
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`     // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`     // 更新时间
}
