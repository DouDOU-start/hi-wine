// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure of table users for DAO operations like Where/Data.
type Users struct {
	g.Meta    `orm:"table:users, do:true"`
	Id        interface{} // 用户ID
	Openid    interface{} // 微信openid
	Nickname  interface{} // 微信昵称
	AvatarUrl interface{} // 微信头像URL
	Phone     interface{} // 手机号
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
