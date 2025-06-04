package v1

import (
	"backend/api/common"
	productv1 "backend/api/product/v1"

	"github.com/gogf/gf/v2/frame/g"
)

// 管理端-套餐商品分组

// 为指定套餐添加可畅饮的酒水
type AdminPackageAddProductsReq struct {
	g.Meta     `path:"/api/v1/admin/packages/{package_id}/products" method:"post" tags:"管理端-套餐商品" summary:"为指定套餐添加可畅饮的酒水"`
	PackageID  int64   `json:"package_id" in:"path" v:"required#套餐ID必填"`
	ProductIDs []int64 `json:"product_ids" v:"required#商品ID列表必填"`
}
type AdminPackageAddProductsRes struct {
	common.Response[struct{}] `json:",inline"`
}

// 从指定套餐中移除某个酒水
type AdminPackageRemoveProductReq struct {
	g.Meta    `path:"/api/v1/admin/packages/{package_id}/products/{product_id}" method:"delete" tags:"管理端-套餐商品" summary:"从指定套餐中移除某个酒水"`
	PackageID int64 `json:"package_id" in:"path" v:"required#套餐ID必填"`
	ProductID int64 `json:"product_id" in:"path" v:"required#商品ID必填"`
}
type AdminPackageRemoveProductRes struct {
	common.Response[struct{}] `json:",inline"`
}

// 查询套餐包含的商品列表
type AdminPackageProductListReq struct {
	g.Meta    `path:"/api/v1/admin/packages/{package_id}/products" method:"get" tags:"管理端-套餐商品" summary:"查询套餐包含的商品列表"`
	PackageID int64 `json:"package_id" in:"path" v:"required#套餐ID必填"`
}
type AdminPackageProductListRes struct {
	common.Response[struct {
		List []productv1.Product `json:"list"`
	}] `json:",inline"`
}

// 获取可添加到套餐的商品列表（未添加到该套餐的商品）
type AdminPackageAvailableProductsReq struct {
	g.Meta    `path:"/api/v1/admin/packages/{package_id}/available-products" method:"get" tags:"管理端-套餐商品" summary:"获取可添加到套餐的商品列表"`
	PackageID int64  `json:"package_id" in:"path" v:"required#套餐ID必填"`
	Keyword   string `json:"keyword" in:"query" description:"商品名称关键字搜索"`
	Page      int    `json:"page" in:"query" description:"页码，默认1"`
	Limit     int    `json:"limit" in:"query" description:"每页数量，默认10"`
}
type AdminPackageAvailableProductsRes struct {
	common.Response[struct {
		List  []productv1.Product `json:"list"`
		Total int                 `json:"total"`
	}] `json:",inline"`
}

// 批量从套餐中移除商品
type AdminPackageBatchRemoveProductsReq struct {
	g.Meta     `path:"/api/v1/admin/packages/{package_id}/batch-remove-products" method:"post" tags:"管理端-套餐商品" summary:"批量从套餐中移除商品"`
	PackageID  int64   `json:"package_id" in:"path" v:"required#套餐ID必填"`
	ProductIDs []int64 `json:"product_ids" v:"required#商品ID列表必填"`
}
type AdminPackageBatchRemoveProductsRes struct {
	common.Response[struct{}] `json:",inline"`
}
