package service

import (
	"backend/api"
	"backend/internal/dao"
	"backend/internal/utility/jwt"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type IAdminService interface {
	Login(ctx context.Context, req *api.AdminLoginReq) (res *api.AdminLoginRes, err error)
	Info(ctx context.Context, req *api.AdminInfoReq) (res *api.AdminInfoRes, err error)
	UserList(ctx context.Context, req *api.UserListReq) (res *api.UserListRes, err error)
	UserDetail(ctx context.Context, req *api.UserDetailReq) (res *api.UserDetailRes, err error)
	UserStatusUpdate(ctx context.Context, req *api.UserStatusUpdateReq) (res *api.UserStatusUpdateRes, err error)
}

var adminService = localAdminService{}

type localAdminService struct{}

func Admin() IAdminService {
	return &adminService
}

// 管理员登录
func (s *localAdminService) Login(ctx context.Context, req *api.AdminLoginReq) (res *api.AdminLoginRes, err error) {
	// 这里实现登录逻辑
	res = &api.AdminLoginRes{}

	// 验证用户名和密码
	user, err := dao.User.Ctx(ctx).Where("nickname=? AND role=2", req.Username).One()
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, gerror.New("用户名或密码错误")
	}

	// 生成token
	token, err := jwt.GenerateToken(user["id"].Int64())
	if err != nil {
		return nil, err
	}

	res.Token = token
	res.UserInfo = user

	return res, nil
}

// 获取管理员信息
func (s *localAdminService) Info(ctx context.Context, req *api.AdminInfoReq) (res *api.AdminInfoRes, err error) {
	res = &api.AdminInfoRes{}

	// 从上下文中获取用户ID
	userId := Auth().GetUserId(ctx)
	if userId == 0 {
		return nil, gerror.New("未登录")
	}

	// 查询用户信息
	user, err := dao.User.Ctx(ctx).Where("id=? AND role=2", userId).One()
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	res.User = user

	return res, nil
}

// 获取用户列表
func (s *localAdminService) UserList(ctx context.Context, req *api.UserListReq) (res *api.UserListRes, err error) {
	res = &api.UserListRes{}

	// 构建查询条件
	model := dao.User.Ctx(ctx)

	if req.Username != "" {
		model = model.Where("nickname LIKE ?", "%"+req.Username+"%")
	}

	if req.Phone != "" {
		model = model.Where("phone LIKE ?", "%"+req.Phone+"%")
	}

	if req.Status != "" {
		model = model.Where("status=?", req.Status)
	}

	// 分页查询
	count, err := model.Count()
	if err != nil {
		return nil, err
	}

	list, err := model.Page(req.Page, req.Size).Order("id DESC").All()
	if err != nil {
		return nil, err
	}

	res.Total = count
	res.List = list

	return res, nil
}

// 获取用户详情
func (s *localAdminService) UserDetail(ctx context.Context, req *api.UserDetailReq) (res *api.UserDetailRes, err error) {
	res = &api.UserDetailRes{}

	// 查询用户信息
	user, err := dao.User.Ctx(ctx).Where("id=?", req.Id).One()
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	res.User = user

	return res, nil
}

// 更新用户状态
func (s *localAdminService) UserStatusUpdate(ctx context.Context, req *api.UserStatusUpdateReq) (res *api.UserStatusUpdateRes, err error) {
	res = &api.UserStatusUpdateRes{}

	// 查询用户是否存在
	user, err := dao.User.Ctx(ctx).Where("id=?", req.Id).One()
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	// 更新用户状态
	_, err = dao.User.Ctx(ctx).Data(g.Map{
		"status": req.Status,
	}).Where("id=?", req.Id).Update()
	if err != nil {
		return nil, err
	}

	res.Success = true

	return res, nil
}
