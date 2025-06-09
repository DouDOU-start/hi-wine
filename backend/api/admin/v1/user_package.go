package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 管理端-用户套餐分组

type AdminUserPackage struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	UserName     string  `json:"user_name"`  // 用户名称
	UserPhone    string  `json:"user_phone"` // 用户手机号
	PackageID    int64   `json:"package_id"`
	PackageName  string  `json:"package_name"`  // 套餐名称
	PackagePrice float64 `json:"package_price"` // 套餐价格
	OrderID      int64   `json:"order_id"`
	OrderSN      string  `json:"order_sn"`     // 订单编号
	StartTime    string  `json:"start_time"`   // 开始时间
	EndTime      string  `json:"end_time"`     // 结束时间
	ValidPeriod  string  `json:"valid_period"` // 有效期描述（如：7天）
	Status       string  `json:"status"`       // active, expired, pending, refunded
	StatusDesc   string  `json:"status_desc"`  // 状态描述（如：已激活、已过期等）
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// 查询用户套餐购买和使用记录（分页、筛选）
type AdminUserPackageListReq struct {
	g.Meta    `path:"user-packages" method:"get" tags:"管理端-用户套餐" summary:"查询用户套餐购买和使用记录"`
	UserID    int64  `json:"user_id" in:"query" description:"用户ID筛选"`
	PackageID int64  `json:"package_id" in:"query" description:"套餐ID筛选"`
	Status    string `json:"status" in:"query" description:"状态筛选（active, expired, pending, refunded）"`
	StartDate string `json:"start_date" in:"query" description:"起始日期(YYYY-MM-DD)"`
	EndDate   string `json:"end_date" in:"query" description:"结束日期(YYYY-MM-DD)"`
	Page      int    `json:"page" in:"query" description:"页码，默认1"`
	Limit     int    `json:"limit" in:"query" description:"每页数量，默认10"`
}
type AdminUserPackageListRes struct {
	common.Response[struct {
		List  []AdminUserPackage `json:"list"`
		Total int                `json:"total"`
	}] `json:",inline"`
}

// 获取用户套餐详情
type AdminUserPackageDetailReq struct {
	g.Meta        `path:"user-packages/{user_package_id}" method:"get" tags:"管理端-用户套餐" summary:"获取用户套餐详情"`
	UserPackageID int64 `json:"user_package_id" in:"path" v:"required#用户套餐ID必填"`
}
type AdminUserPackageDetailRes struct {
	common.Response[AdminUserPackage] `json:",inline"`
}

// 创建用户套餐
type AdminUserPackageCreateReq struct {
	g.Meta    `path:"user-packages" method:"post" tags:"管理端-用户套餐" summary:"创建用户套餐"`
	UserID    int64  `json:"user_id" v:"required#用户ID必填"`
	PackageID int64  `json:"package_id" v:"required#套餐ID必填"`
	StartTime string `json:"start_time" description:"开始时间，仅当状态为active时必填"`
	Status    string `json:"status" v:"required#状态必填"`
	OrderID   int64  `json:"order_id" description:"关联订单ID，可选"`
}
type AdminUserPackageCreateRes struct {
	common.Response[AdminUserPackage] `json:",inline"`
}

// 更新用户套餐状态
type AdminUserPackageUpdateStatusReq struct {
	g.Meta        `path:"user-packages/{user_package_id}/status" method:"put" tags:"管理端-用户套餐" summary:"更新用户套餐状态"`
	UserPackageID int64  `json:"user_package_id" in:"path" v:"required#用户套餐ID必填"`
	Status        string `json:"status" v:"required#状态必填"`
	Reason        string `json:"reason" description:"状态变更原因（例如退款原因）"`
}
type AdminUserPackageUpdateStatusRes struct {
	common.Response[AdminUserPackage] `json:",inline"`
}

// 查询用户有效套餐
type AdminUserActivePackagesReq struct {
	g.Meta `path:"users/{user_id}/active-packages" method:"get" tags:"管理端-用户套餐" summary:"查询用户有效套餐"`
	UserID int64 `json:"user_id" in:"path" v:"required#用户ID必填"`
}
type AdminUserActivePackagesRes struct {
	common.Response[struct {
		List []AdminUserPackage `json:"list"`
	}] `json:",inline"`
}

// 获取用户套餐详细信息（包含用户信息和套餐信息）
type AdminUserPackageFullDetailReq struct {
	g.Meta        `path:"user-packages/{user_package_id}/full-detail" method:"get" tags:"管理端-用户套餐" summary:"获取用户套餐详细信息（包含用户信息和套餐信息）"`
	UserPackageID int64 `json:"user_package_id" in:"path" v:"required#用户套餐ID必填"`
}

type AdminUserPackageFullDetail struct {
	// 用户套餐基本信息
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	UserName    string `json:"user_name"` // 用户昵称
	PackageID   int64  `json:"package_id"`
	PackageName string `json:"package_name"` // 套餐名称
	OrderID     int64  `json:"order_id"`
	OrderSN     string `json:"order_sn"` // 订单编号
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	ValidPeriod string `json:"valid_period"` // 有效期描述（如：7天）
	Status      string `json:"status"`       // active, expired, pending, refunded
	StatusDesc  string `json:"status_desc"`  // 状态描述（如：已激活、已过期等）
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`

	// 用户信息
	User struct {
		ID        int64  `json:"id"`
		Nickname  string `json:"nickname"`
		AvatarUrl string `json:"avatar_url"`
		Phone     string `json:"phone"`
	} `json:"user"`

	// 套餐信息
	Package struct {
		ID              int64   `json:"id"`
		Name            string  `json:"name"`
		Description     string  `json:"description"`
		Price           float64 `json:"price"`
		DurationMinutes int     `json:"duration_minutes"`
	} `json:"package"`

	// 订单信息
	Order struct {
		ID        int64   `json:"id"`
		OrderSN   string  `json:"order_sn"`
		TotalFee  float64 `json:"total_fee"`
		PayTime   string  `json:"pay_time"`
		PayStatus string  `json:"pay_status"`
	} `json:"order"`

	// 使用情况统计
	Usage struct {
		TotalUsedTimes int    `json:"total_used_times"` // 总使用次数
		LastUsedTime   string `json:"last_used_time"`   // 最后使用时间
	} `json:"usage"`
}

type AdminUserPackageFullDetailRes struct {
	common.Response[AdminUserPackageFullDetail] `json:",inline"`
}
