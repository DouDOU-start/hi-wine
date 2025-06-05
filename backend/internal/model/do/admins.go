// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Admins is the golang structure of table admins for DAO operations like Where/Data.
type Admins struct {
	g.Meta      `orm:"table:admins, do:true"`
	Id          interface{} // 管理员ID
	Username    interface{} // 用户名
	Password    interface{} // 密码（哈希）
	Role        interface{} // 角色
	IsActive    interface{} // 是否激活
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	LastLoginAt *gtime.Time // 最后登录时间
} 