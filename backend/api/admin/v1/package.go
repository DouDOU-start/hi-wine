package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 管理端-套餐分组

type AdminPackage struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	DurationHours int     `json:"duration_hours"`
	IsActive      bool    `json:"is_active"`
}

// 获取套餐列表（分页、筛选、模糊搜索）
type AdminPackageListReq struct {
	g.Meta `path:"/api/v1/admin/packages" method:"get" tags:"管理端-套餐" summary:"获取套餐列表（分页、筛选、模糊搜索）"`
	Page   int    `json:"page" in:"query" description:"页码，默认1"`
	Limit  int    `json:"limit" in:"query" description:"每页数量，默认10"`
	Name   string `json:"name" in:"query" description:"套餐名模糊搜索"`
}
type AdminPackageListRes struct {
	common.Response[struct {
		List  []AdminPackage `json:"list"`
		Total int            `json:"total"`
	}] `json:",inline"`
}

// 创建套餐
type AdminPackageCreateReq struct {
	g.Meta        `path:"/api/v1/admin/packages" method:"post" tags:"管理端-套餐" summary:"创建套餐"`
	Name          string  `json:"name" v:"required#套餐名必填"`
	Price         float64 `json:"price" v:"required#价格必填"`
	DurationHours int     `json:"duration_hours" v:"required#有效时长必填"`
	Description   string  `json:"description"`
	IsActive      *bool   `json:"is_active"`
}
type AdminPackageCreateRes struct {
	common.Response[AdminPackage] `json:",inline"`
}

// 更新套餐
type AdminPackageUpdateReq struct {
	g.Meta        `path:"/api/v1/admin/packages/{package_id}" method:"put" tags:"管理端-套餐" summary:"更新套餐"`
	PackageID     int64   `json:"package_id" in:"path" v:"required#套餐ID必填"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	DurationHours int     `json:"duration_hours"`
	Description   string  `json:"description"`
	IsActive      *bool   `json:"is_active"`
}
type AdminPackageUpdateRes struct {
	common.Response[AdminPackage] `json:",inline"`
}

// 删除套餐
type AdminPackageDeleteReq struct {
	g.Meta    `path:"/api/v1/admin/packages/{package_id}" method:"delete" tags:"管理端-套餐" summary:"删除套餐"`
	PackageID int64 `json:"package_id" in:"path" v:"required#套餐ID必填"`
}
type AdminPackageDeleteRes struct {
	common.Response[struct{}] `json:",inline"`
}
