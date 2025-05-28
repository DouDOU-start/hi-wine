package service

import (
	"backend/api"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

type sProduct struct{}

var (
	productService = sProduct{}
)

// Product 商品服务
func Product() *sProduct {
	return &productService
}

// List 获取商品列表
func (s *sProduct) List(ctx context.Context, req *api.ProductListReq) (res *api.ProductListRes, err error) {
	res = &api.ProductListRes{}

	// 构建查询条件
	m := dao.Product.Ctx(ctx)

	// 按分类筛选
	if req.CategoryId > 0 {
		m = m.Where(dao.Product.Columns().CategoryId, req.CategoryId)
	}

	// 按关键词搜索
	if req.Keyword != "" {
		m = m.WhereLike(dao.Product.Columns().Name, "%"+req.Keyword+"%")
	}

	// 只查询上架商品
	m = m.Where(dao.Product.Columns().Status, 1)

	// 查询总数
	count, err := m.Count()
	if err != nil {
		return nil, err
	}
	res.Total = count

	// 分页查询
	var products []*entity.Product
	err = m.Page(req.Page, req.Size).Order(dao.Product.Columns().Id + " DESC").Scan(&products)
	if err != nil {
		return nil, err
	}

	res.List = products
	return res, nil
}

// Detail 获取商品详情
func (s *sProduct) Detail(ctx context.Context, req *api.ProductDetailReq) (res *api.ProductDetailRes, err error) {
	res = &api.ProductDetailRes{}

	// 查询商品
	product, err := dao.Product.Ctx(ctx).Where(dao.Product.Columns().Id, req.Id).One()
	if err != nil {
		return nil, err
	}
	if product.IsEmpty() {
		return nil, gerror.New("商品不存在")
	}

	var productEntity *entity.Product
	if err = product.Struct(&productEntity); err != nil {
		return nil, err
	}

	res.Product = productEntity
	return res, nil
}

// CategoryList 获取商品分类列表
func (s *sProduct) CategoryList(ctx context.Context, req *api.CategoryListReq) (res *api.CategoryListRes, err error) {
	res = &api.CategoryListRes{}

	// 查询分类列表
	var categories []*entity.Category
	err = dao.Category.Ctx(ctx).Order(dao.Category.Columns().Sort + " ASC, " + dao.Category.Columns().Id + " ASC").Scan(&categories)
	if err != nil {
		return nil, err
	}

	res.List = categories
	return res, nil
}
