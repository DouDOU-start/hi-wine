package v1

import (
	"backend/api/common"
	productv1 "backend/api/product/v1"

	"github.com/gogf/gf/v2/frame/g"
)

// 管理端-商品分组

// 获取商品列表（分页、筛选、模糊搜索）
type AdminProductListReq struct {
	g.Meta     `path:"products" method:"get" tags:"管理端-商品" summary:"获取商品列表（分页、筛选、模糊搜索）"`
	Page       int    `json:"page" in:"query" description:"页码，默认1"`
	Limit      int    `json:"limit" in:"query" description:"每页数量，默认10"`
	Name       string `json:"name" in:"query" description:"商品名模糊搜索"`
	CategoryID int64  `json:"category_id" in:"query" description:"分类ID精确筛选"`
	IsActive   *bool  `json:"is_active" in:"query" description:"是否上架筛选"`
}
type AdminProductListRes struct {
	common.Response[struct {
		List  []productv1.UserProduct `json:"list"`
		Total int                     `json:"total"`
	}] `json:",inline"`
}

// 创建商品
type AdminProductCreateReq struct {
	g.Meta      `path:"products" method:"post" tags:"管理端-商品" summary:"创建商品"`
	Name        string  `json:"name" v:"required#商品名必填"`
	CategoryID  int64   `json:"category_id" v:"required#分类必填"`
	Price       float64 `json:"price" v:"required#价格必填"`
	Stock       int     `json:"stock" v:"required#库存必填"`
	ImageURL    string  `json:"image_url" v:"required#商品图片必填"`
	Description string  `json:"description"`
	IsActive    *bool   `json:"is_active"`
}
type AdminProductCreateRes struct {
	common.Response[productv1.UserProduct] `json:",inline"`
}

// 更新商品
type AdminProductUpdateReq struct {
	g.Meta      `path:"products/{product_id}" method:"put" tags:"管理端-商品" summary:"更新商品"`
	ProductID   int64   `json:"product_id" in:"path" v:"required#商品ID必填"`
	Name        string  `json:"name"`
	CategoryID  int64   `json:"category_id"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	IsActive    *bool   `json:"is_active"`
}
type AdminProductUpdateRes struct {
	common.Response[productv1.UserProduct] `json:",inline"`
}

// 删除商品
type AdminProductDeleteReq struct {
	g.Meta    `path:"products/{product_id}" method:"delete" tags:"管理端-商品" summary:"删除商品"`
	ProductID int64 `json:"product_id" in:"path" v:"required#商品ID必填"`
}
type AdminProductDeleteRes struct {
	common.Response[struct{}] `json:",inline"`
}

// 获取商品详情
type AdminProductDetailReq struct {
	g.Meta    `path:"products/{product_id}" method:"get" tags:"管理端-商品" summary:"获取商品详情"`
	ProductID int64 `json:"product_id" in:"path" v:"required#商品ID必填"`
}
type AdminProductDetailRes struct {
	common.Response[productv1.UserProduct] `json:",inline"`
}
