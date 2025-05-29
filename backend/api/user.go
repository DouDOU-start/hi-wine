package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 用户信息更新请求
type UserUpdateReq struct {
	g.Meta   `path:"/user/update" method:"post" summary:"更新用户信息"`
	Nickname string `json:"nickname" description:"昵称" v:"required"`
	Avatar   string `json:"avatar" description:"头像" v:"required"`
}

// 用户信息更新响应
type UserUpdateRes struct {
	Success bool        `json:"success" description:"更新结果"`
	User    interface{} `json:"user" description:"更新后的用户信息"`
}

// 用户信息获取请求
type UserInfoReq struct {
	g.Meta `path:"/user/info" method:"get" summary:"获取用户信息"`
}

// 用户信息获取响应
type UserInfoRes struct {
	User interface{} `json:"user" description:"用户信息"`
}
