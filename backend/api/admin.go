package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 管理员登录请求
type AdminLoginReq struct {
	g.Meta   `path:"/login" method:"post" summary:"管理员登录"`
	Username string `json:"username" description:"用户名" v:"required"`
	Password string `json:"password" description:"密码" v:"required"`
}

// 管理员登录响应
type AdminLoginRes struct {
	Token    string      `json:"token" description:"登录token"`
	UserInfo interface{} `json:"userInfo" description:"用户信息"`
}

// 管理员信息请求
type AdminInfoReq struct {
	g.Meta `path:"/admin/info" method:"get" summary:"获取管理员信息"`
}

// 管理员信息响应
type AdminInfoRes struct {
	User interface{} `json:"user" description:"用户信息"`
}

// 用户列表请求
type UserListReq struct {
	g.Meta   `path:"/user/list" method:"get" summary:"获取用户列表"`
	Username string `json:"username" description:"用户名" v:""`
	Phone    string `json:"phone" description:"手机号" v:""`
	Status   string `json:"status" description:"状态" v:""`
	Page     int    `json:"page" description:"页码" v:"min:1" d:"1"`
	Size     int    `json:"size" description:"每页数量" v:"max:50" d:"10"`
}

// 用户列表响应
type UserListRes struct {
	List  interface{} `json:"list" description:"用户列表"`
	Total int         `json:"total" description:"总数"`
}

// 用户详情请求
type UserDetailReq struct {
	g.Meta `path:"/user/detail" method:"get" summary:"获取用户详情"`
	Id     int64 `json:"id" description:"用户ID" v:"required"`
}

// 用户详情响应
type UserDetailRes struct {
	User interface{} `json:"user" description:"用户详情"`
}

// 更新用户状态请求
type UserStatusUpdateReq struct {
	g.Meta `path:"/user/status" method:"post" summary:"更新用户状态"`
	Id     int64 `json:"id" description:"用户ID" v:"required"`
	Status int   `json:"status" description:"状态" v:"required|in:0,1"`
}

// 更新用户状态响应
type UserStatusUpdateRes struct {
	Success bool `json:"success" description:"是否成功"`
}
