package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 管理端-用户分组

// 管理员用户列表查询
type AdminUserListReq struct {
	g.Meta    `path:"/users" method:"get" tags:"管理端-用户" summary:"获取用户列表"`
	Page      int    `json:"page" in:"query" description:"页码，默认1"`
	Limit     int    `json:"limit" in:"query" description:"每页数量，默认10"`
	Keyword   string `json:"keyword" in:"query" description:"关键词搜索（用户昵称、手机号）"`
	StartDate string `json:"start_date" in:"query" description:"起始日期(YYYY-MM-DD)"`
	EndDate   string `json:"end_date" in:"query" description:"结束日期(YYYY-MM-DD)"`
}

type AdminUserListRes struct {
	common.Response[struct {
		List  []UserInfo `json:"list"`
		Total int        `json:"total"`
	}] `json:",inline"`
}

// 管理员用户详情查询
type AdminUserDetailReq struct {
	g.Meta `path:"/users/{user_id}" method:"get" tags:"管理端-用户" summary:"获取用户详情"`
	UserID int64 `json:"user_id" in:"path" v:"required#用户ID必填"`
}

type AdminUserDetailRes struct {
	common.Response[UserInfo] `json:",inline"`
}

// 用户信息结构体（管理端视图）
type UserInfo struct {
	ID        int64  `json:"id"`
	Openid    string `json:"openid"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
