// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure of table user for DAO operations like Where/Data.
type User struct {
	g.Meta     `orm:"table:user, do:true"`
	Id         interface{} // 主键
	Openid     interface{} // 微信openid
	Nickname   interface{} // 昵称
	Avatar     interface{} // 头像
	Phone      interface{} // 手机号
	Role       interface{} // 角色（0顾客，1店员，2管理员）
	CreateTime *gtime.Time // 创建时间
	UpdateTime *gtime.Time // 更新时间
}
