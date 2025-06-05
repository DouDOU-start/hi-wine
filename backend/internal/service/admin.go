package service

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/utility/jwt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/crypto/bcrypt"
)

// AdminService 管理员服务接口
type AdminService interface {
	// Login 管理员登录
	Login(ctx context.Context, username, password string) (*v1.AdminLoginRes, error)
}

// 管理员服务实现
type adminService struct{}

// 单例实例
var adminServiceInstance = adminService{}

// Admin 获取管理员服务实例
func Admin() AdminService {
	return &adminServiceInstance
}

// Login 管理员登录
func (s *adminService) Login(ctx context.Context, username, password string) (*v1.AdminLoginRes, error) {
	// 查询管理员信息
	var admin *entity.Admins
	err := dao.Admins.Ctx(ctx).
		Where(dao.Admins.Columns().Username, username).
		Scan(&admin)
	if err != nil {
		return nil, err
	}

	// 检查管理员是否存在
	if admin == nil {
		return nil, gerror.New("用户名或密码错误")
	}

	// 检查管理员是否激活
	if admin.IsActive == 0 {
		return nil, gerror.New("账号已被禁用，请联系超级管理员")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return nil, gerror.New("用户名或密码错误")
	}

	// 生成JWT令牌
	tokenData := g.Map{
		"id":       admin.Id,
		"username": admin.Username,
		"role":     admin.Role,
	}
	token, err := jwt.GenerateToken(tokenData)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	_, err = dao.Admins.Ctx(ctx).
		Data(g.Map{
			dao.Admins.Columns().LastLoginAt: gtime.Now(),
		}).
		Where(dao.Admins.Columns().Id, admin.Id).
		Update()
	if err != nil {
		g.Log().Error(ctx, "更新管理员最后登录时间失败:", err)
		// 不返回错误，继续执行
	}

	// 返回登录结果
	return &v1.AdminLoginRes{
		Token: token,
		AdminUser: v1.AdminUser{
			ID:       int64(admin.Id),
			Username: admin.Username,
			Role:     admin.Role,
		},
	}, nil
}
