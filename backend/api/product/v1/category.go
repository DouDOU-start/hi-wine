package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 用户端-商品分类分组

type UserCategory struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

// 获取所有商品分类
type UserCategoryListReq struct {
	g.Meta `path:"/categories" method:"get" tags:"商品" summary:"获取所有已激活的商品分类列表"`
}
type UserCategoryListRes struct {
	common.Response[struct {
		List []UserCategory `json:"list"`
	}] `json:",inline"`
}
