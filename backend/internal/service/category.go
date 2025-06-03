package service

import (
	"context"
	"time"

	v1 "backend/api/admin/v1"
	"backend/internal/dao"
	"backend/internal/model"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// Category 分类服务
type Category struct{}

// List 获取分类列表
func (s *Category) List(ctx context.Context, req *v1.AdminCategoryListReq) (list []v1.AdminCategory, total int, err error) {
	// 调用DAO层获取数据
	categoryDao := dao.Category{}
	categories, total, err := categoryDao.List(ctx, req.Page, req.Limit, req.Name)
	if err != nil {
		return nil, 0, err
	}

	// 转换为API响应格式
	list = make([]v1.AdminCategory, len(categories))
	for i, category := range categories {
		list[i] = v1.AdminCategory{
			ID:        category.ID,
			Name:      category.Name,
			SortOrder: category.SortOrder,
			IsActive:  category.IsActive,
		}
	}

	return list, total, nil
}

// GetByID 根据ID获取分类
func (s *Category) GetByID(ctx context.Context, id int64) (*v1.AdminCategory, error) {
	categoryDao := dao.Category{}
	category, err := categoryDao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, gerror.New("分类不存在")
	}

	return &v1.AdminCategory{
		ID:        category.ID,
		Name:      category.Name,
		SortOrder: category.SortOrder,
		IsActive:  category.IsActive,
	}, nil
}

// Create 创建分类
func (s *Category) Create(ctx context.Context, req *v1.AdminCreateCategoryReq) (*v1.AdminCategory, error) {
	categoryDao := dao.Category{}

	// 检查分类名称是否已存在
	existingCategory, err := categoryDao.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existingCategory != nil {
		return nil, gerror.New("分类名称已存在")
	}

	// 创建模型
	category := &model.Category{
		Name:      req.Name,
		SortOrder: req.SortOrder,
		IsActive:  req.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 默认值处理
	if category.SortOrder == 0 {
		category.SortOrder = 100 // 默认排序值
	}
	if !req.IsActive {
		category.IsActive = true // 默认激活
	}

	// 插入数据
	id, err := categoryDao.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	// 返回创建后的分类
	return s.GetByID(ctx, id)
}

// Update 更新分类
func (s *Category) Update(ctx context.Context, req *v1.AdminUpdateCategoryReq) (*v1.AdminCategory, error) {
	categoryDao := dao.Category{}

	// 检查分类是否存在
	exists, err := categoryDao.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if exists == nil {
		return nil, gerror.New("分类不存在")
	}

	// 如果要更改名称，检查新名称是否已被使用
	if req.Name != "" && req.Name != exists.Name {
		existingCategory, err := categoryDao.GetByName(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		if existingCategory != nil {
			return nil, gerror.New("分类名称已存在")
		}
	}

	// 构建更新数据
	data := g.Map{}
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.SortOrder > 0 {
		data["sort_order"] = req.SortOrder
	}
	data["is_active"] = req.IsActive
	data["updated_at"] = time.Now()

	// 更新数据
	if err := categoryDao.Update(ctx, req.ID, data); err != nil {
		return nil, err
	}

	// 返回更新后的分类
	return s.GetByID(ctx, req.ID)
}

// Delete 删除分类
func (s *Category) Delete(ctx context.Context, id int64) error {
	categoryDao := dao.Category{}

	// 检查分类是否存在
	exists, err := categoryDao.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if exists == nil {
		return gerror.New("分类不存在")
	}

	// 检查是否有商品关联此分类
	hasProducts, err := categoryDao.HasProducts(ctx, id)
	if err != nil {
		return err
	}
	if hasProducts {
		return gerror.New("该分类下有商品，无法删除")
	}

	// 删除分类
	return categoryDao.Delete(ctx, id)
}
