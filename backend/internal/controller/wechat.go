package controller

import (
	"backend/api"
	"backend/internal/service"
	"context"
)

type WechatController struct{}

func NewWechat() *WechatController {
	return &WechatController{}
}

// 微信小程序登录
func (c *WechatController) Login(ctx context.Context, req *api.WechatLoginReq) (res *api.WechatLoginRes, err error) {
	return service.Wechat().Login(ctx, req)
}
