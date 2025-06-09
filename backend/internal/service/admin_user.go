package service

import (
	"context"

	adminv1 "backend/api/admin/v1"
	"backend/internal/dao"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
)

// AdminUserService 管理员用户服务接口
type AdminUserService interface {
	// GetUserList 获取用户列表
	GetUserList(ctx context.Context, req *adminv1.AdminUserListReq) (list []adminv1.UserInfo, total int, err error)

	// GetUserDetail 获取用户详情
	GetUserDetail(ctx context.Context, userId int64) (*adminv1.UserInfo, error)
}

// 管理员用户服务实现
type adminUserService struct{}

// 单例实例
var adminUserServiceInstance = adminUserService{}

// AdminUser 获取管理员用户服务实例
func AdminUser() AdminUserService {
	return &adminUserServiceInstance
}

// GetUserList 获取用户列表
func (s *adminUserService) GetUserList(ctx context.Context, req *adminv1.AdminUserListReq) (list []adminv1.UserInfo, total int, err error) {
	// 1. 构建查询条件
	m := dao.Users.Ctx(ctx)

	// 1.1 关键词搜索
	if req.Keyword != "" {
		m = m.WhereOr(
			dao.Users.Columns().Nickname+" LIKE ?", "%"+req.Keyword+"%",
		).WhereOr(
			dao.Users.Columns().Phone+" LIKE ?", "%"+req.Keyword+"%",
		)
	}

	// 1.2 日期范围筛选
	if req.StartDate != "" {
		m = m.WhereGTE(dao.Users.Columns().CreatedAt, req.StartDate+" 00:00:00")
	}
	if req.EndDate != "" {
		m = m.WhereLTE(dao.Users.Columns().CreatedAt, req.EndDate+" 23:59:59")
	}

	// 2. 查询总数
	total, err = m.Count()
	if err != nil {
		return nil, 0, err
	}

	// 3. 分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// 4. 查询数据
	var users []*entity.Users
	err = m.Page(page, limit).
		Order(dao.Users.Columns().Id + " DESC").
		Scan(&users)
	if err != nil {
		return nil, 0, err
	}

	// 5. 转换为API响应格式
	list = make([]adminv1.UserInfo, len(users))
	for i, user := range users {
		list[i] = adminv1.UserInfo{
			ID:        int64(user.Id),
			Openid:    user.Openid,
			Nickname:  user.Nickname,
			AvatarURL: user.AvatarUrl,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		}
	}

	return list, total, nil
}

// GetUserDetail 获取用户详情
func (s *adminUserService) GetUserDetail(ctx context.Context, userId int64) (*adminv1.UserInfo, error) {
	// 1. 查询用户信息
	var user *entity.Users
	err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, userId).
		Scan(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	// 2. 转换为API响应格式
	userInfo := &adminv1.UserInfo{
		ID:        int64(user.Id),
		Openid:    user.Openid,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarUrl,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}

	return userInfo, nil
}
