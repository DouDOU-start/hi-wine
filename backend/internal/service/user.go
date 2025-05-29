package service

import (
	"backend/api"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sUser struct{}

var (
	userService = sUser{}
)

func User() *sUser {
	return &userService
}

// Update 更新用户信息
func (s *sUser) Update(ctx context.Context, req *api.UserUpdateReq) (res *api.UserUpdateRes, err error) {
	res = &api.UserUpdateRes{}

	// 获取当前用户ID
	userId := ctx.Value("userId")
	if userId == nil {
		return nil, gerror.New("用户未登录")
	}

	// 更新用户信息
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().Id, userId).
		Data(g.Map{
			dao.User.Columns().Nickname: req.Nickname,
			dao.User.Columns().Avatar:   req.Avatar,
		}).
		Update()
	if err != nil {
		return nil, err
	}

	// 查询更新后的用户信息
	user := &entity.User{}
	err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().Id, userId).
		Scan(user)
	if err != nil {
		return nil, err
	}

	res.Success = true
	res.User = user
	return res, nil
}

// GetInfo 获取用户信息
func (s *sUser) GetInfo(ctx context.Context, req *api.UserInfoReq) (res *api.UserInfoRes, err error) {
	res = &api.UserInfoRes{}

	// 获取当前用户ID
	userId := ctx.Value("userId")
	if userId == nil {
		return nil, gerror.New("用户未登录")
	}

	// 查询用户信息
	user := &entity.User{}
	err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().Id, userId).
		Scan(user)
	if err != nil {
		return nil, err
	}

	res.User = user
	return res, nil
}
