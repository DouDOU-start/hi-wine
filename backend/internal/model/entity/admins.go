// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Admins is the golang structure for table admins.
type Admins struct {
	Id         int         `json:"id"         orm:"id"          description:"管理员ID"`
	Username   string      `json:"username"   orm:"username"    description:"用户名"`
	Password   string      `json:"password"   orm:"password"    description:"密码（哈希）"`
	Role       string      `json:"role"       orm:"role"        description:"角色"`
	IsActive   int         `json:"isActive"   orm:"is_active"   description:"是否激活"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:"更新时间"`
	LastLoginAt *gtime.Time `json:"lastLoginAt" orm:"last_login_at" description:"最后登录时间"`
} 