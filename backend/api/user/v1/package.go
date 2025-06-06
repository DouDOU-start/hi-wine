package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 用户端-套餐分组

type UserPackage struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	DurationHours int     `json:"duration_hours"`
	IsActive      bool    `json:"is_active"`
}

type UserPackageDetail struct {
	UserPackage
	// 可扩展更多详情字段
}

// 获取套餐列表 - 公开API
type UserPackageListReq struct {
	g.Meta `path:"/packages" method:"get" tags:"套餐" summary:"获取所有可购买的套餐列表"`
	Page   int    `json:"page" in:"query" description:"页码，默认1"`
	Limit  int    `json:"limit" in:"query" description:"每页数量，默认10"`
	Name   string `json:"name" in:"query" description:"套餐名模糊搜索"`
}
type UserPackageListRes struct {
	common.Response[struct {
		List  []UserPackage `json:"list"`
		Total int           `json:"total"`
	}] `json:",inline"`
}

// 获取套餐详情 - 公开API
type UserPackageDetailReq struct {
	g.Meta    `path:"/packages/{package_id}" method:"get" tags:"套餐" summary:"获取套餐详情"`
	PackageID int64 `json:"package_id" in:"path" description:"套餐ID" v:"required#套餐ID必填"`
}
type UserPackageDetailRes struct {
	common.Response[UserPackageDetail] `json:",inline"`
}

// 购买套餐 - 需要认证
type UserBuyPackageReq struct {
	g.Meta    `path:"/packages/{package_id}/buy" method:"post" tags:"套餐" summary:"购买套餐"`
	PackageID int64 `json:"package_id" in:"path" description:"套餐ID" v:"required#套餐ID必填"`
}
type UserBuyPackageRes struct {
	common.Response[struct {
		OrderID int64 `json:"order_id"`
	}] `json:",inline"`
}

// 获取用户个人套餐列表
// GET /api/v1/user/my-packages
// 返回: 用户套餐列表
type UserMyPackagesReq struct {
	g.Meta `path:"/my-packages" method:"get" tags:"套餐" summary:"获取当前登录用户的套餐列表"`
	Status string `json:"status" in:"query" description:"套餐状态筛选 (active-激活中, pending-待激活, expired-已过期, 不传则查询全部)"`
}

type UserMyPackagesRes struct {
	common.Response[struct {
		List []UserMyPackage `json:"list"`
	}] `json:",inline"`
}

// 用户套餐信息结构体
type UserMyPackage struct {
	ID            int64   `json:"id"`             // 用户套餐ID
	PackageID     int64   `json:"package_id"`     // 套餐ID
	PackageName   string  `json:"package_name"`   // 套餐名称
	Price         float64 `json:"price"`          // 套餐价格
	Status        string  `json:"status"`         // 套餐状态 (active-激活中, pending-待激活, expired-已过期)
	StartTime     string  `json:"start_time"`     // 开始时间
	EndTime       string  `json:"end_time"`       // 结束时间
	OrderID       int64   `json:"order_id"`       // 关联订单ID
	OrderSN       string  `json:"order_sn"`       // 关联订单号
	CreatedAt     string  `json:"created_at"`     // 创建时间
	RemainingTime int64   `json:"remaining_time"` // 剩余时间（秒）
}
