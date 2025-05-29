package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 管理员分类列表请求
type AdminCategoryListReq struct {
	g.Meta `path:"/category/list" method:"get" summary:"获取分类列表(管理员)"`
	Page   int `json:"page" description:"页码" v:"min:1" d:"1"`
	Size   int `json:"size" description:"每页数量" v:"max:50" d:"10"`
}

// 管理员分类列表响应
type AdminCategoryListRes struct {
	List  interface{} `json:"list" description:"分类列表"`
	Total int         `json:"total" description:"总数"`
}

// 管理员分类详情请求
type AdminCategoryDetailReq struct {
	g.Meta `path:"/category/detail" method:"get" summary:"获取分类详情(管理员)"`
	Id     int64 `json:"id" description:"分类ID" v:"required"`
}

// 管理员分类详情响应
type AdminCategoryDetailRes struct {
	Category interface{} `json:"category" description:"分类详情"`
}

// 管理员添加分类请求
type AdminCategoryAddReq struct {
	g.Meta `path:"/category/add" method:"post" summary:"添加分类(管理员)"`
	Name   string `json:"name" description:"分类名称" v:"required"`
	Icon   string `json:"icon" description:"分类图标" v:""`
	Sort   int    `json:"sort" description:"排序" v:"min:0" d:"0"`
}

// 管理员添加分类响应
type AdminCategoryAddRes struct {
	Id int64 `json:"id" description:"分类ID"`
}

// 管理员更新分类请求
type AdminCategoryUpdateReq struct {
	g.Meta `path:"/category/update" method:"post" summary:"更新分类(管理员)"`
	Id     int64  `json:"id" description:"分类ID" v:"required"`
	Name   string `json:"name" description:"分类名称" v:"required"`
	Icon   string `json:"icon" description:"分类图标" v:""`
	Sort   int    `json:"sort" description:"排序" v:"min:0" d:"0"`
}

// 管理员更新分类响应
type AdminCategoryUpdateRes struct {
	Success bool `json:"success" description:"是否成功"`
}

// 管理员删除分类请求
type AdminCategoryDeleteReq struct {
	g.Meta `path:"/category/delete" method:"post" summary:"删除分类(管理员)"`
	Id     int64 `json:"id" description:"分类ID" v:"required"`
}

// 管理员删除分类响应
type AdminCategoryDeleteRes struct {
	Success bool `json:"success" description:"是否成功"`
}
