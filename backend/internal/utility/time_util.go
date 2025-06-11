package utility

import (
	adminv1 "backend/api/admin/v1"
	userv1 "backend/api/user/v1"
	"backend/internal/consts"
	"backend/internal/model/entity"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// FormatTimeOrEmpty 格式化时间，如果时间为nil或零值则返回空字符串
func FormatTimeOrEmpty(t *gtime.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(consts.TimeFormatStandard)
}

// FormatStdTimeOrEmpty 格式化标准time.Time时间，如果时间为零值则返回空字符串
func FormatStdTimeOrEmpty(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatUserPackageTimes 格式化用户套餐的时间字段
func FormatUserPackageTimes(src *entity.UserPackages, dst *adminv1.AdminUserPackage) {
	dst.StartTime = FormatTimeOrEmpty(src.StartTime)
	dst.EndTime = FormatTimeOrEmpty(src.EndTime)
	dst.CreatedAt = FormatTimeOrEmpty(src.CreatedAt)
	dst.UpdatedAt = FormatTimeOrEmpty(src.UpdatedAt)
}

// FormatUserMyPackage 格式化用户个人套餐的时间字段并计算剩余时间
func FormatUserMyPackage(src *entity.UserPackages, dst *userv1.UserMyPackage) {
	// 基本字段转换
	if err := gconv.Struct(src, dst); err != nil {
		return
	}

	// 格式化时间
	dst.StartTime = FormatTimeOrEmpty(src.StartTime)
	dst.EndTime = FormatTimeOrEmpty(src.EndTime)
	dst.CreatedAt = FormatTimeOrEmpty(src.CreatedAt)

	// 计算剩余时间（秒）
	if src.Status == consts.PackageStatusActive && src.EndTime != nil {
		now := gtime.Now()
		if src.EndTime.After(now) {
			dst.RemainingTime = int64(src.EndTime.Sub(now) / 1000000000) // 转换为秒
		}
	}
}
