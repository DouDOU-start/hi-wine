package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 用户端-商品分类分组

// UserCategory 商品分类信息
type UserCategory struct {
	ID           int64  `json:"id"`                      // 分类ID
	Name         string `json:"name"`                    // 分类名称
	SortOrder    int    `json:"sort_order"`              // 排序顺序
	ImageURL     string `json:"image_url,omitempty"`     // 分类图片URL
	ProductCount int    `json:"product_count,omitempty"` // 分类下商品数量
}

// 获取所有商品分类
type UserCategoryListReq struct {
	g.Meta `path:"/categories" method:"get" tags:"商品" summary:"获取所有已激活的商品分类列表"`
}
type UserCategoryListRes struct {
	common.Response[struct {
		List  []UserCategory `json:"list"`
		Total int            `json:"total"`
	}] `json:",inline"`
}
