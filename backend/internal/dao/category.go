package dao

import (
	"context"

	"backend/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// Category 分类数据访问对象
type Category struct{}

// List 获取分类列表
func (d *Category) List(ctx context.Context, page, limit int, name string) (list []model.Category, total int, err error) {
	m := g.DB().Model("categories").Safe()

	// 如果有名称搜索条件
	if name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}

	// 查询总数
	total, err = m.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	err = m.Page(page, limit).
		Order("sort_order ASC, id DESC").
		Scan(&list)

	return
}

// GetByID 根据ID获取分类
func (d *Category) GetByID(ctx context.Context, id int64) (category *model.Category, err error) {
	err = g.DB().Model("categories").
		Where("id", id).
		Scan(&category)
	return
}

// GetByName 根据名称获取分类
func (d *Category) GetByName(ctx context.Context, name string) (category *model.Category, err error) {
	err = g.DB().Model("categories").
		Where("name", name).
		Scan(&category)
	return
}

// Create 创建分类
func (d *Category) Create(ctx context.Context, data *model.Category) (id int64, err error) {
	result, err := g.DB().Model("categories").Data(data).Insert()
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	return
}

// Update 更新分类
func (d *Category) Update(ctx context.Context, id int64, data g.Map) error {
	_, err := g.DB().Model("categories").
		Where("id", id).
		Data(data).
		Update()
	return err
}

// Delete 删除分类
func (d *Category) Delete(ctx context.Context, id int64) error {
	_, err := g.DB().Model("categories").
		Where("id", id).
		Delete()
	return err
}

// HasProducts 检查分类是否有关联商品
func (d *Category) HasProducts(ctx context.Context, id int64) (bool, error) {
	count, err := g.DB().Model("products").
		Where("category_id", id).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
