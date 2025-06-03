package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 获取所有商品分类
// 管理端接口
// type AdminCategoryListReq struct {
// 	g.Meta `path:"/api/v1/admin/categories" method:"get" tags:"管理端-商品" summary:"获取所有商品分类（管理端）"`
// }

type CategoryListReq struct {
	g.Meta `path:"/api/v1/categories" method:"get" tags:"商品" summary:"获取所有商品分类"`
}

type CategoryListRes struct {
	List []Category `json:"list"`
}

// 获取某分类下商品列表
// 管理端接口
// type AdminProductListReq struct {
// 	g.Meta     `path:"/api/v1/admin/categories/{category_id}/products" method:"get" tags:"管理端-商品" summary:"获取分类下商品列表（管理端）"`
// }

type ProductListReq struct {
	g.Meta     `path:"/api/v1/categories/{category_id}/products" method:"get" tags:"商品" summary:"获取分类下商品列表"`
	CategoryID int64 `json:"category_id" in:"path" description:"分类ID" v:"required#分类ID不能为空"`
	Page       int   `json:"page" in:"query" description:"页码"`
	Limit      int   `json:"limit" in:"query" description:"每页数量"`
}

type ProductListRes struct {
	List []Product `json:"list"`
}

// 获取商品详情
// 管理端接口
// type AdminProductDetailReq struct {
// 	g.Meta    `path:"/api/v1/admin/products/{product_id}" method:"get" tags:"管理端-商品" summary:"获取商品详情（管理端）"`
// }

type ProductDetailReq struct {
	g.Meta    `path:"/api/v1/products/{product_id}" method:"get" tags:"商品" summary:"获取商品详情"`
	ProductID int64 `json:"product_id" in:"path" description:"商品ID" v:"required#商品ID不能为空"`
}

type ProductDetailRes struct {
	ProductDetail
}

// Category 商品分类
// GET /api/v1/categories
type Category struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
	ImageURL  string `json:"image_url,omitempty"`
}

// Product 商品信息
type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
	Description string  `json:"description,omitempty"`
}

// ProductDetail 商品详情
// GET /api/v1/products/{product_id}
type ProductDetail struct {
	Product
	// 可扩展更多详情字段
}

// 用户端-商品分组

type UserProduct struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
	Description string  `json:"description,omitempty"`
}

type UserProductDetail struct {
	UserProduct
	// 可扩展更多详情字段
}

// 获取某分类下的商品列表
type UserProductListReq struct {
	g.Meta     `path:"/categories/{category_id}/products" method:"get" tags:"商品" summary:"获取指定分类下的所有已激活商品"`
	CategoryID int64 `json:"category_id" in:"path" description:"分类ID" v:"required#分类ID必填"`
	Page       int   `json:"page" in:"query" description:"页码，默认1"`
	Limit      int   `json:"limit" in:"query" description:"每页数量，默认10"`
}
type UserProductListRes struct {
	common.Response[struct {
		List  []UserProduct `json:"list"`
		Total int           `json:"total"`
	}] `json:",inline"`
}

// 获取商品详情
type UserProductDetailReq struct {
	g.Meta    `path:"/products/{product_id}" method:"get" tags:"商品" summary:"获取单个商品的详细信息"`
	ProductID int64 `json:"product_id" in:"path" description:"商品ID" v:"required#商品ID必填"`
}
type UserProductDetailRes struct {
	common.Response[UserProductDetail] `json:",inline"`
}
