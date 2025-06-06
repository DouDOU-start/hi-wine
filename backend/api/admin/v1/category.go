package v1

import (
	"backend/api/common"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// 管理端-分类分组

type AdminCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// 获取分类列表
type AdminCategoryListReq struct {
	g.Meta `path:"categories" method:"get" tags:"管理端-分类" summary:"获取分类列表"`
	Page   int    `json:"page" in:"query" description:"页码，默认1"`
	Limit  int    `json:"limit" in:"query" description:"每页数量，默认10"`
	Name   string `json:"name" in:"query" description:"分类名称模糊搜索"`
}
type AdminCategoryListRes struct {
	common.Response[struct {
		List  []AdminCategory `json:"list"`
		Total int             `json:"total"`
	}] `json:",inline"`
}

// 创建分类
type AdminCreateCategoryReq struct {
	g.Meta    `path:"categories" method:"post" tags:"管理端-分类" summary:"创建分类"`
	Name      string `json:"name" v:"required#分类名称不能为空"`
	SortOrder int    `json:"sort_order,omitempty"`
	IsActive  bool   `json:"is_active,omitempty"`
}
type AdminCreateCategoryRes struct {
	common.Response[AdminCategory] `json:",inline"`
}

// 更新分类
type AdminUpdateCategoryReq struct {
	g.Meta    `path:"categories/{id}" method:"put" tags:"管理端-分类" summary:"更新分类"`
	ID        int64  `json:"id" in:"path" v:"required#分类ID不能为空"`
	Name      string `json:"name,omitempty"`
	SortOrder int    `json:"sort_order,omitempty"`
	IsActive  bool   `json:"is_active,omitempty"`
}
type AdminUpdateCategoryRes struct {
	common.Response[AdminCategory] `json:",inline"`
}

// 删除分类
type AdminDeleteCategoryReq struct {
	g.Meta `path:"categories/{id}" method:"delete" tags:"管理端-分类" summary:"删除分类"`
	ID     int64 `json:"id" in:"path" v:"required#分类ID不能为空"`
}
type AdminDeleteCategoryRes struct {
	common.Response[interface{}] `json:",inline"`
}

// 获取分类详情
type AdminCategoryDetailReq struct {
	g.Meta `path:"categories/{id}" method:"get" tags:"管理端-分类" summary:"获取分类详情"`
	ID     int64 `json:"id" in:"path" v:"required#分类ID不能为空"`
}
type AdminCategoryDetailRes struct {
	common.Response[AdminCategory] `json:",inline"`
}
