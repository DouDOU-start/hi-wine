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
