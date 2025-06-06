package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

// CronService 定时任务服务接口
type CronService interface {
	// StartCronJobs 启动所有定时任务
	StartCronJobs() error

	// CheckExpiredPackages 检查并更新过期的用户套餐
	CheckExpiredPackages(ctx context.Context) error
}

// cronService 定时任务服务实现
type cronService struct{}

// 确保 cronService 实现了 CronService 接口
var _ CronService = (*cronService)(nil)

// 单例实例
var cronServiceInstance = cronService{}

// Cron 获取定时任务服务实例
func Cron() CronService {
	return &cronServiceInstance
}

// StartCronJobs 启动所有定时任务
func (s *cronService) StartCronJobs() error {
	ctx := gctx.New()

	// 每天凌晨0点执行一次，检查过期的用户套餐
	// 使用GoFrame的cron表达式格式：秒 分 时 日 月 周
	_, err := gcron.Add(ctx, "0 0 0 * * *", func(ctx context.Context) {
		if err := s.CheckExpiredPackages(ctx); err != nil {
			g.Log().Error(ctx, "检查过期套餐失败:", err)
		}
	}, "check_expired_packages")

	return err
}

// CheckExpiredPackages 检查并更新过期的用户套餐
func (s *cronService) CheckExpiredPackages(ctx context.Context) error {
	g.Log().Info(ctx, "开始检查过期的用户套餐...")

	// 当前时间
	now := gtime.Now()

	// 批量更新所有已激活但已过期的套餐
	result, err := g.DB().Model("user_packages").
		Data(g.Map{
			"status":     "expired",
			"updated_at": gtime.Now(),
		}).
		Where("status", "active").
		WhereLT("end_time", now).
		Update()

	if err != nil {
		g.Log().Errorf(ctx, "批量更新过期套餐状态失败: %v", err)
		return err
	}

	// 获取影响的行数
	affected, err := result.RowsAffected()
	if err != nil {
		g.Log().Errorf(ctx, "获取影响行数失败: %v", err)
		return err
	}

	g.Log().Infof(ctx, "成功将 %d 个已过期的套餐状态更新为已过期", affected)

	return nil
}
