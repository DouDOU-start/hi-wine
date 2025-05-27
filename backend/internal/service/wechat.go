package service

import (
	"backend/api"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v4"
)

type sWechat struct{}

var (
	wechatService = sWechat{}
)

func Wechat() *sWechat {
	return &wechatService
}

// Login 处理微信小程序登录
func (s *sWechat) Login(ctx context.Context, req *api.WechatLoginReq) (res *api.WechatLoginRes, err error) {
	// 1. 读取配置
	cfg, _ := gcfg.Instance().Data(ctx)
	appid := gconv.String(cfg["wechat.appid"])
	secret := gconv.String(cfg["wechat.secret"])
	jwtSecret := []byte("hi-wine-jwt-secret") // 建议放配置

	// 2. 请求微信服务器换 openid
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appid, secret, req.Code)
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	var wxResp struct {
		Openid     string `json:"openid"`
		SessionKey string `json:"session_key"`
		Errcode    int    `json:"errcode"`
		Errmsg     string `json:"errmsg"`
	}
	if err := json.NewDecoder(rsp.Body).Decode(&wxResp); err != nil {
		return nil, err
	}
	if wxResp.Openid == "" {
		return nil, fmt.Errorf("微信登录失败: %s", wxResp.Errmsg)
	}

	// 3. 查找或创建用户，并同步昵称和头像
	user := &entity.User{}
	dbUser, err := dao.User.Ctx(ctx).Where("openid", wxResp.Openid).One()
	if err != nil {
		return nil, err
	}
	if dbUser.IsEmpty() {
		// 新用户，插入
		user.Openid = wxResp.Openid
		user.Nickname = req.Nickname
		user.Avatar = req.Avatar
		user.Role = 0
		id, err := dao.User.Ctx(ctx).InsertAndGetId(user)
		if err != nil {
			return nil, err
		}
		user.Id = id
	} else {
		if err := dbUser.Struct(user); err != nil {
			return nil, err
		}
		// 老用户，更新昵称和头像
		_, err = dao.User.Ctx(ctx).Where("id", user.Id).Data(g.Map{
			"nickname": req.Nickname,
			"avatar":   req.Avatar,
		}).Update()
		if err != nil {
			return nil, err
		}
		user.Nickname = req.Nickname
		user.Avatar = req.Avatar
	}

	// 4. 生成 JWT token
	claims := jwt.MapClaims{
		"userId": user.Id,
		"openid": user.Openid,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(), // 7天有效
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenObj.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	return &api.WechatLoginRes{
		Token:    token,
		UserInfo: user,
	}, nil
}
