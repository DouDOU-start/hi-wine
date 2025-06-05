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
