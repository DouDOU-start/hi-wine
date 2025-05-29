package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 管理员商品列表请求
type AdminProductListReq struct {
	g.Meta     `path:"/product/list" method:"get" summary:"获取商品列表(管理员)"`
	Name       string `json:"name" description:"商品名称" v:""`
	CategoryId int64  `json:"categoryId" description:"分类ID" v:""`
	Status     string `json:"status" description:"状态" v:""`
	Page       int    `json:"page" description:"页码" v:"min:1" d:"1"`
	Size       int    `json:"size" description:"每页数量" v:"max:50" d:"10"`
}

// 管理员商品列表响应
type AdminProductListRes struct {
	List  interface{} `json:"list" description:"商品列表"`
	Total int         `json:"total" description:"总数"`
}

// 管理员商品详情请求
type AdminProductDetailReq struct {
	g.Meta `path:"/product/detail" method:"get" summary:"获取商品详情(管理员)"`
	Id     int64 `json:"id" description:"商品ID" v:"required"`
}

// 管理员商品详情响应
type AdminProductDetailRes struct {
	Product interface{} `json:"product" description:"商品详情"`
}

// 管理员添加商品请求
type AdminProductAddReq struct {
	g.Meta      `path:"/product/add" method:"post" summary:"添加商品(管理员)"`
	Name        string  `json:"name" description:"商品名称" v:"required"`
	CategoryId  int64   `json:"categoryId" description:"分类ID" v:"required"`
	Price       float64 `json:"price" description:"价格" v:"required|min:0.01"`
	Stock       int     `json:"stock" description:"库存" v:"required|min:0"`
	Image       string  `json:"image" description:"商品图片" v:"required"`
	Status      int     `json:"status" description:"状态（0下架，1上架）" v:"required|in:0,1"`
	Description string  `json:"description" description:"商品描述" v:""`
}

// 管理员添加商品响应
type AdminProductAddRes struct {
	Id int64 `json:"id" description:"商品ID"`
}

// 管理员更新商品请求
type AdminProductUpdateReq struct {
	g.Meta      `path:"/product/update" method:"post" summary:"更新商品(管理员)"`
	Id          int64   `json:"id" description:"商品ID" v:"required"`
	Name        string  `json:"name" description:"商品名称" v:"required"`
	CategoryId  int64   `json:"categoryId" description:"分类ID" v:"required"`
	Price       float64 `json:"price" description:"价格" v:"required|min:0.01"`
	Stock       int     `json:"stock" description:"库存" v:"required|min:0"`
	Image       string  `json:"image" description:"商品图片" v:"required"`
	Status      int     `json:"status" description:"状态（0下架，1上架）" v:"required|in:0,1"`
	Description string  `json:"description" description:"商品描述" v:""`
}

// 管理员更新商品响应
type AdminProductUpdateRes struct {
	Success bool `json:"success" description:"是否成功"`
}

// 管理员删除商品请求
type AdminProductDeleteReq struct {
	g.Meta `path:"/product/delete" method:"post" summary:"删除商品(管理员)"`
	Id     int64 `json:"id" description:"商品ID" v:"required"`
}

// 管理员删除商品响应
type AdminProductDeleteRes struct {
	Success bool `json:"success" description:"是否成功"`
}

// 管理员更新商品状态请求
type AdminProductStatusUpdateReq struct {
	g.Meta `path:"/product/status" method:"post" summary:"更新商品状态(管理员)"`
	Id     int64 `json:"id" description:"商品ID" v:"required"`
	Status int   `json:"status" description:"状态（0下架，1上架）" v:"required|in:0,1"`
}

// 管理员更新商品状态响应
type AdminProductStatusUpdateRes struct {
	Success bool `json:"success" description:"是否成功"`
}
