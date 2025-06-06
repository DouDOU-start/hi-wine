package utility

import (
	adminv1 "backend/api/admin/v1"
	userv1 "backend/api/user/v1"
	"backend/internal/consts"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// FormatUserPackageTimes 格式化用户套餐的时间字段
func FormatUserPackageTimes(src *entity.UserPackages, dst *adminv1.AdminUserPackage) {
	if src.StartTime != nil {
		dst.StartTime = src.StartTime.Format(consts.TimeFormatStandard)
	}
	if src.EndTime != nil {
		dst.EndTime = src.EndTime.Format(consts.TimeFormatStandard)
	}
	if src.CreatedAt != nil {
		dst.CreatedAt = src.CreatedAt.Format(consts.TimeFormatStandard)
	}
	if src.UpdatedAt != nil {
		dst.UpdatedAt = src.UpdatedAt.Format(consts.TimeFormatStandard)
	}
}

// FormatUserMyPackage 格式化用户个人套餐的时间字段并计算剩余时间
func FormatUserMyPackage(src *entity.UserPackages, dst *userv1.UserMyPackage) {
	// 基本字段转换
	if err := gconv.Struct(src, dst); err != nil {
		return
	}

	// 格式化时间
	if src.StartTime != nil {
		dst.StartTime = src.StartTime.Format(consts.TimeFormatStandard)
	}
	if src.EndTime != nil {
		dst.EndTime = src.EndTime.Format(consts.TimeFormatStandard)
	}
	if src.CreatedAt != nil {
		dst.CreatedAt = src.CreatedAt.Format(consts.TimeFormatStandard)
	}

	// 计算剩余时间（秒）
	if src.Status == consts.PackageStatusActive && src.EndTime != nil {
		now := gtime.Now()
		if src.EndTime.After(now) {
			dst.RemainingTime = int64(src.EndTime.Sub(now) / 1000000000) // 转换为秒
		}
	}
}
