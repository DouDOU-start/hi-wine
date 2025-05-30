package service

import (
	"backend/api"
	"backend/internal/dao"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type IAdminCategoryService interface {
	List(ctx context.Context, req *api.AdminCategoryListReq) (res *api.AdminCategoryListRes, err error)
	Detail(ctx context.Context, req *api.AdminCategoryDetailReq) (res *api.AdminCategoryDetailRes, err error)
	Add(ctx context.Context, req *api.AdminCategoryAddReq) (res *api.AdminCategoryAddRes, err error)
	Update(ctx context.Context, req *api.AdminCategoryUpdateReq) (res *api.AdminCategoryUpdateRes, err error)
	Delete(ctx context.Context, req *api.AdminCategoryDeleteReq) (res *api.AdminCategoryDeleteRes, err error)
}

var adminCategoryService = localAdminCategoryService{}

type localAdminCategoryService struct{}

func AdminCategory() IAdminCategoryService {
	return &adminCategoryService
}

// 获取分类列表
func (s *localAdminCategoryService) List(ctx context.Context, req *api.AdminCategoryListReq) (res *api.AdminCategoryListRes, err error) {
	res = &api.AdminCategoryListRes{}

	// 构建查询条件
	model := dao.Category.Ctx(ctx)

	// 查询总数
	count, err := model.Count()
	if err != nil {
		return nil, err
	}

	// 判断是否需要分页
	var list interface{}
	if req.Page > 0 && req.Size > 0 {
		// 使用分页查询
		list, err = model.Page(req.Page, req.Size).Order("sort ASC, id DESC").All()
	} else {
		// 查询全部数据
		list, err = model.Order("sort ASC, id DESC").All()
	}

	if err != nil {
		return nil, err
	}

	res.Total = count
	res.List = list

	return res, nil
}

// 获取分类详情
func (s *localAdminCategoryService) Detail(ctx context.Context, req *api.AdminCategoryDetailReq) (res *api.AdminCategoryDetailRes, err error) {
	res = &api.AdminCategoryDetailRes{}

	// 查询分类信息
	category, err := dao.Category.Ctx(ctx).Where("id=?", req.Id).One()
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, gerror.New("分类不存在")
	}

	res.Category = category

	return res, nil
}

// 添加分类
func (s *localAdminCategoryService) Add(ctx context.Context, req *api.AdminCategoryAddReq) (res *api.AdminCategoryAddRes, err error) {
	res = &api.AdminCategoryAddRes{}

	// 检查分类名称是否已存在
	count, err := dao.Category.Ctx(ctx).Where("name=?", req.Name).Count()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, gerror.New("分类名称已存在")
	}

	// 添加分类
	id, err := dao.Category.Ctx(ctx).Data(g.Map{
		"name": req.Name,
		"icon": req.Icon,
		"sort": req.Sort,
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}

	res.Id = id

	return res, nil
}

// 更新分类
func (s *localAdminCategoryService) Update(ctx context.Context, req *api.AdminCategoryUpdateReq) (res *api.AdminCategoryUpdateRes, err error) {
	res = &api.AdminCategoryUpdateRes{}

	// 检查分类是否存在
	category, err := dao.Category.Ctx(ctx).Where("id=?", req.Id).One()
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, gerror.New("分类不存在")
	}

	// 检查分类名称是否已存在
	count, err := dao.Category.Ctx(ctx).Where("name=? AND id<>?", req.Name, req.Id).Count()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, gerror.New("分类名称已存在")
	}

	// 更新分类
	_, err = dao.Category.Ctx(ctx).Data(g.Map{
		"name": req.Name,
		"icon": req.Icon,
		"sort": req.Sort,
	}).Where("id=?", req.Id).Update()
	if err != nil {
		return nil, err
	}

	res.Success = true

	return res, nil
}

// 删除分类
func (s *localAdminCategoryService) Delete(ctx context.Context, req *api.AdminCategoryDeleteReq) (res *api.AdminCategoryDeleteRes, err error) {
	res = &api.AdminCategoryDeleteRes{}

	// 检查分类是否存在
	category, err := dao.Category.Ctx(ctx).Where("id=?", req.Id).One()
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, gerror.New("分类不存在")
	}

	// 检查分类是否有关联商品
	count, err := dao.Product.Ctx(ctx).Where("category_id=?", req.Id).Count()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, gerror.New("该分类下有关联商品，无法删除")
	}

	// 删除分类
	_, err = dao.Category.Ctx(ctx).Where("id=?", req.Id).Delete()
	if err != nil {
		return nil, err
	}

	res.Success = true

	return res, nil
}
