package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

type WechatLoginReq struct {
	g.Meta   `path:"/wechat/login" method:"post" summary:"微信小程序登录"`
	Code     string `json:"code" description:"微信登录凭证code"`
	Nickname string `json:"nickname" description:"昵称"`
	Avatar   string `json:"avatar" description:"头像"`
}

type WechatLoginRes struct {
	Token    string      `json:"token" description:"登录token"`
	UserInfo interface{} `json:"userInfo" description:"用户信息"`
}
