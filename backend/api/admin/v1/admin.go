package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 管理员登录
type AdminLoginReq struct {
	g.Meta   `path:"/api/v1/admin/login" method:"post" tags:"管理员" summary:"管理员登录"`
	Username string `json:"username" description:"用户名" v:"required#用户名不能为空"`
	Password string `json:"password" description:"密码" v:"required#密码不能为空"`
}

type AdminLoginRes struct {
	Token     string    `json:"token"`
	AdminUser AdminUser `json:"admin_user"`
}

// 管理员信息结构体
type AdminUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
