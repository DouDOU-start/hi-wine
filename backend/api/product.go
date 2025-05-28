package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 商品列表请求
type ProductListReq struct {
	g.Meta     `path:"/product/list" method:"get" summary:"获取商品列表"`
	CategoryId int64  `json:"categoryId" description:"分类ID" v:""`
	Keyword    string `json:"keyword" description:"搜索关键词" v:""`
	Page       int    `json:"page" description:"页码" v:"min:1" d:"1"`
	Size       int    `json:"size" description:"每页数量" v:"max:50" d:"10"`
}

// 商品列表响应
type ProductListRes struct {
	List  interface{} `json:"list" description:"商品列表"`
	Total int         `json:"total" description:"总数"`
}

// 商品详情请求
type ProductDetailReq struct {
	g.Meta `path:"/product/detail" method:"get" summary:"获取商品详情"`
	Id     int64 `json:"id" description:"商品ID" v:"required"`
}

// 商品详情响应
type ProductDetailRes struct {
	Product interface{} `json:"product" description:"商品详情"`
}

// 商品分类列表请求
type CategoryListReq struct {
	g.Meta `path:"/category/list" method:"get" summary:"获取商品分类列表"`
}

// 商品分类列表响应
type CategoryListRes struct {
	List interface{} `json:"list" description:"分类列表"`
}
