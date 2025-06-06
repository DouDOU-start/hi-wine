package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 管理员登录
type AdminLoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" tags:"管理员" summary:"管理员登录"`
	Username string `json:"username" description:"用户名" v:"required#用户名不能为空"`
	Password string `json:"password" description:"密码" v:"required#密码不能为空"`
}

type AdminLoginRes struct {
	common.Response[struct {
		Token     string    `json:"token"`
		AdminUser AdminUser `json:"admin_user"`
	}] `json:",inline"`
}

// 获取当前管理员信息
type AdminProfileReq struct {
	g.Meta `path:"/profile" method:"get" tags:"管理员" summary:"获取当前管理员信息"`
}

type AdminProfileRes struct {
	common.Response[AdminUser] `json:",inline"`
}

// 管理员信息结构体
type AdminUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
