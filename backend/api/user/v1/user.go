package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 用户端-用户分组

// 微信登录/注册
// 微信登录
// POST /auth/wechat-login
// code: 微信登录凭证
// 返回: token, user信息
type WechatLoginReq struct {
	g.Meta `path:"/auth/wechat-login" method:"post" tags:"用户" summary:"微信登录/注册"`
	Code   string `json:"code" description:"微信登录凭证" v:"required#code不能为空"`
}

type WechatLoginRes struct {
	common.Response[struct {
		Token string      `json:"token"`
		User  UserProfile `json:"user"`
	}] `json:",inline"`
}

// 获取用户个人信息
// GET /user/profile
// 返回: user信息
type UserProfileReq struct {
	g.Meta `path:"/user/profile" method:"get" tags:"用户" summary:"获取当前登录用户的个人信息"`
}

type UserProfileRes struct {
	common.Response[UserProfile] `json:",inline"`
}

// 更新用户个人信息
// PUT /user/profile
// 参数: phone, nickname, avatar_url
// 返回: user信息
type UpdateUserProfileReq struct {
	g.Meta    `path:"/user/profile" method:"put" tags:"用户" summary:"更新当前登录用户的个人信息"`
	Phone     string `json:"phone,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type UpdateUserProfileRes struct {
	common.Response[UserProfile] `json:",inline"`
}

// 用户信息结构体
type UserProfile struct {
	ID        int64  `json:"id"`
	Openid    string `json:"openid"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Phone     string `json:"phone"`
}
