package v1

import (
	"backend/api/common"
	productv1 "backend/api/product/v1"

	"github.com/gogf/gf/v2/frame/g"
)

// 管理端-套餐分组

type AdminPackage struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// 套餐使用统计
type AdminPackageStats struct {
	TotalSales    int     `json:"total_sales"`     // 总销售数量
	ActiveUsers   int     `json:"active_users"`    // 当前活跃用户数
	TotalRevenue  float64 `json:"total_revenue"`   // 总收入
	AvgDailyUsage float64 `json:"avg_daily_usage"` // 平均每日使用量
	ProductsCount int     `json:"products_count"`  // 包含商品数量
}

// 获取套餐列表（分页、筛选、模糊搜索）
type AdminPackageListReq struct {
	g.Meta `path:"packages" method:"get" tags:"管理端-套餐" summary:"获取套餐列表（分页、筛选、模糊搜索）"`
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

// 获取套餐详情
type AdminPackageDetailReq struct {
	g.Meta    `path:"packages/{package_id}" method:"get" tags:"管理端-套餐" summary:"获取套餐详情"`
	PackageID int64 `json:"package_id" in:"path" v:"required#套餐ID必填"`
}
type AdminPackageDetailRes struct {
	common.Response[AdminPackage] `json:",inline"`
}

// 创建套餐
type AdminPackageCreateReq struct {
	g.Meta          `path:"packages" method:"post" tags:"管理端-套餐" summary:"创建套餐"`
	Name            string  `json:"name" v:"required#套餐名必填"`
	Price           float64 `json:"price" v:"required#价格必填"`
	DurationMinutes int     `json:"duration_minutes" v:"required#有效时长必填"`
	Description     string  `json:"description"`
	IsActive        *bool   `json:"is_active"`
}
type AdminPackageCreateRes struct {
	common.Response[AdminPackage] `json:",inline"`
}

// 更新套餐
type AdminPackageUpdateReq struct {
	g.Meta          `path:"packages/{package_id}" method:"put" tags:"管理端-套餐" summary:"更新套餐"`
	PackageID       int64   `json:"package_id" in:"path" v:"required#套餐ID必填"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
	Description     string  `json:"description"`
	IsActive        *bool   `json:"is_active"`
}
type AdminPackageUpdateRes struct {
	common.Response[AdminPackage] `json:",inline"`
}

// 删除套餐
type AdminPackageDeleteReq struct {
	g.Meta    `path:"packages/{package_id}" method:"delete" tags:"管理端-套餐" summary:"删除套餐"`
	PackageID int64 `json:"package_id" in:"path" v:"required#套餐ID必填"`
}
type AdminPackageDeleteRes struct {
	common.Response[struct{}] `json:",inline"`
}

// 获取套餐使用统计
type AdminPackageStatsReq struct {
	g.Meta    `path:"packages/{package_id}/stats" method:"get" tags:"管理端-套餐" summary:"获取套餐使用统计"`
	PackageID int64 `json:"package_id" in:"path" v:"required#套餐ID必填"`
}
type AdminPackageStatsRes struct {
	common.Response[AdminPackageStats] `json:",inline"`
}

// 获取带商品列表的套餐详情
type AdminPackageWithProductsReq struct {
	g.Meta    `path:"packages/{package_id}/with-products" method:"get" tags:"管理端-套餐" summary:"获取带商品列表的套餐详情"`
	PackageID int64 `json:"package_id" in:"path" v:"required#套餐ID必填"`
}

type AdminPackageWithProducts struct {
	// 套餐基本信息
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
	DurationDays    int     `json:"duration_days"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`

	// 套餐包含的商品列表
	Products []productv1.UserProduct `json:"products"`

	// 商品数量
	ProductsCount int `json:"products_count"`
}

type AdminPackageWithProductsRes struct {
	common.Response[AdminPackageWithProducts] `json:",inline"`
}

// 获取套餐详细信息（包含基本信息、统计信息和商品列表）
type AdminPackageFullDetailReq struct {
	g.Meta    `path:"packages/{package_id}/full-detail" method:"get" tags:"管理端-套餐" summary:"获取套餐详细信息（包含基本信息、统计信息和商品列表）"`
	PackageID int64 `json:"package_id" in:"path" v:"required#套餐ID必填"`
}

type AdminPackageFullDetail struct {
	// 套餐基本信息
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
	DurationDays    int     `json:"duration_days"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`

	// 套餐统计信息
	Stats AdminPackageStats `json:"stats"`

	// 套餐包含的商品列表
	Products []productv1.UserProduct `json:"products"`

	// 最近购买记录
	RecentPurchases []struct {
		ID           int64  `json:"id"`
		UserID       int64  `json:"user_id"`
		UserName     string `json:"user_name"`
		OrderID      int64  `json:"order_id"`
		OrderSN      string `json:"order_sn"`
		PurchaseTime string `json:"purchase_time"`
		Status       string `json:"status"`
	} `json:"recent_purchases"`
}

type AdminPackageFullDetailRes struct {
	common.Response[AdminPackageFullDetail] `json:",inline"`
}
