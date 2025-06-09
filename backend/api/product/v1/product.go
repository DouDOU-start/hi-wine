package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// ===== 用户端商品API =====

// UserProduct 商品信息
type UserProduct struct {
	ID          int64   `json:"id"`                    // 商品ID
	Name        string  `json:"name"`                  // 商品名称
	Price       float64 `json:"price"`                 // 商品价格
	ImageURL    string  `json:"image_url"`             // 商品图片URL
	Stock       int     `json:"stock"`                 // 库存数量
	Description string  `json:"description,omitempty"` // 商品描述
	CategoryID  int64   `json:"category_id,omitempty"` // 所属分类ID
	Status      int     `json:"status,omitempty"`      // 商品状态(1:上架, 0:下架)
	SalesCount  int     `json:"sales_count,omitempty"` // 销量
	CreatedAt   string  `json:"created_at,omitempty"`  // 创建时间
	UpdatedAt   string  `json:"updated_at,omitempty"`  // 更新时间
}

// UserProductDetail 商品详情
type UserProductDetail struct {
	UserProduct
	// 扩展字段
	Images         []string            `json:"images,omitempty"`         // 商品图片列表
	Specifications []map[string]string `json:"specifications,omitempty"` // 商品规格
}

// 获取某分类下的商品列表
type UserProductListReq struct {
	g.Meta     `path:"/categories/{category_id}/products" method:"get" tags:"商品" summary:"获取指定分类下的所有已激活商品"`
	CategoryID int64  `json:"category_id" in:"path" description:"分类ID" v:"required#分类ID必填"`
	Page       int    `json:"page" in:"query" description:"页码，默认1"`
	Limit      int    `json:"limit" in:"query" description:"每页数量，默认10"`
	SortBy     string `json:"sort_by" in:"query" description:"排序字段(price,sales_count)"`
	SortOrder  string `json:"sort_order" in:"query" description:"排序方式(asc,desc)"`
}

type UserProductListRes struct {
	common.Response[struct {
		List  []UserProduct `json:"list"`  // 商品列表
		Total int           `json:"total"` // 总数量
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

// 搜索商品
type UserProductSearchReq struct {
	g.Meta     `path:"/products/search" method:"get" tags:"商品" summary:"搜索商品"`
	Keyword    string `json:"keyword" in:"query" description:"搜索关键词"`
	CategoryID int64  `json:"category_id" in:"query" description:"分类ID(可选)"`
	Page       int    `json:"page" in:"query" description:"页码，默认1"`
	Limit      int    `json:"limit" in:"query" description:"每页数量，默认10"`
	SortBy     string `json:"sort_by" in:"query" description:"排序字段(price,sales_count)"`
	SortOrder  string `json:"sort_order" in:"query" description:"排序方式(asc,desc)"`
}

type UserProductSearchRes struct {
	common.Response[struct {
		List  []UserProduct `json:"list"`  // 商品列表
		Total int           `json:"total"` // 总数量
	}] `json:",inline"`
}
