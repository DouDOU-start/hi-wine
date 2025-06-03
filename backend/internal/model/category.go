package model

import (
	"time"
)

// Category 商品分类模型
type Category struct {
	ID        int64     `json:"id" orm:"id,primary"`         // 分类ID
	Name      string    `json:"name" orm:"name"`             // 分类名称
	SortOrder int       `json:"sort_order" orm:"sort_order"` // 排序顺序
	IsActive  bool      `json:"is_active" orm:"is_active"`   // 是否激活
	CreatedAt time.Time `json:"created_at" orm:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" orm:"updated_at"` // 更新时间
}

// TableName 返回表名
func (c *Category) TableName() string {
	return "categories"
}
