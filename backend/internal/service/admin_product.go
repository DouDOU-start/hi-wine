package service

import (
	"backend/api"
	"backend/internal/dao"
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type IAdminProductService interface {
	List(ctx context.Context, req *api.AdminProductListReq) (res *api.AdminProductListRes, err error)
	Detail(ctx context.Context, req *api.AdminProductDetailReq) (res *api.AdminProductDetailRes, err error)
	Add(ctx context.Context, req *api.AdminProductAddReq) (res *api.AdminProductAddRes, err error)
	Update(ctx context.Context, req *api.AdminProductUpdateReq) (res *api.AdminProductUpdateRes, err error)
	Delete(ctx context.Context, req *api.AdminProductDeleteReq) (res *api.AdminProductDeleteRes, err error)
	UpdateStatus(ctx context.Context, req *api.AdminProductStatusUpdateReq) (res *api.AdminProductStatusUpdateRes, err error)
}

var adminProductService = localAdminProductService{}

type localAdminProductService struct{}

func AdminProduct() IAdminProductService {
	return &adminProductService
}

// 获取商品列表
func (s *localAdminProductService) List(ctx context.Context, req *api.AdminProductListReq) (res *api.AdminProductListRes, err error) {
	res = &api.AdminProductListRes{}

	// 构建查询条件
	model := dao.Product.Ctx(ctx)

	if req.Name != "" {
		model = model.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.CategoryId > 0 {
		model = model.Where("category_id=?", req.CategoryId)
	}

	if req.Status != "" {
		model = model.Where("status=?", req.Status)
	}

	// 分页查询
	count, err := model.Count()
	if err != nil {
		return nil, err
	}

	list, err := model.Page(req.Page, req.Size).Order("id DESC").All()
	if err != nil {
		return nil, err
	}

	res.Total = count
	res.List = list

	return res, nil
}

// 获取商品详情
func (s *localAdminProductService) Detail(ctx context.Context, req *api.AdminProductDetailReq) (res *api.AdminProductDetailRes, err error) {
	res = &api.AdminProductDetailRes{}

	// 查询商品信息
	product, err := dao.Product.Ctx(ctx).Where("id=?", req.Id).One()
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, gerror.New("商品不存在")
	}

	res.Product = product

	return res, nil
}

// 添加商品
func (s *localAdminProductService) Add(ctx context.Context, req *api.AdminProductAddReq) (res *api.AdminProductAddRes, err error) {
	res = &api.AdminProductAddRes{}

	// 检查分类是否存在
	category, err := dao.Category.Ctx(ctx).Where("id=?", req.CategoryId).One()
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, gerror.New("分类不存在")
	}

	// 检查商品名称是否已存在
	count, err := dao.Product.Ctx(ctx).Where("name=?", req.Name).Count()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, gerror.New("商品名称已存在")
	}

	// 添加商品
	id, err := dao.Product.Ctx(ctx).Data(g.Map{
		"name":        req.Name,
		"category_id": req.CategoryId,
		"price":       req.Price,
		"stock":       req.Stock,
		"image":       req.Image,
		"status":      req.Status,
		"description": req.Description,
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}

	res.Id = id

	return res, nil
}

// 更新商品
func (s *localAdminProductService) Update(ctx context.Context, req *api.AdminProductUpdateReq) (res *api.AdminProductUpdateRes, err error) {
	res = &api.AdminProductUpdateRes{}

	// 检查商品是否存在
	product, err := dao.Product.Ctx(ctx).Where("id=?", req.Id).One()
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, gerror.New("商品不存在")
	}

	// 检查分类是否存在
	category, err := dao.Category.Ctx(ctx).Where("id=?", req.CategoryId).One()
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, gerror.New("分类不存在")
	}

	// 检查商品名称是否已存在
	count, err := dao.Product.Ctx(ctx).Where("name=? AND id<>?", req.Name, req.Id).Count()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, gerror.New("商品名称已存在")
	}

	// 如果图片发生变化，尝试删除旧图片
	oldImageUrl := product["image"].String()
	if oldImageUrl != req.Image && oldImageUrl != "" &&
		(strings.HasPrefix(oldImageUrl, "http://") || strings.HasPrefix(oldImageUrl, "https://")) {
		// 忽略删除图片的错误，不影响商品更新
		_ = Minio().Delete(ctx, oldImageUrl)
	}

	// 更新商品
	_, err = dao.Product.Ctx(ctx).Data(g.Map{
		"name":        req.Name,
		"category_id": req.CategoryId,
		"price":       req.Price,
		"stock":       req.Stock,
		"image":       req.Image,
		"status":      req.Status,
		"description": req.Description,
	}).Where("id=?", req.Id).Update()
	if err != nil {
		return nil, err
	}

	res.Success = true

	return res, nil
}

// 删除商品
func (s *localAdminProductService) Delete(ctx context.Context, req *api.AdminProductDeleteReq) (res *api.AdminProductDeleteRes, err error) {
	res = &api.AdminProductDeleteRes{}

	// 检查商品是否存在
	product, err := dao.Product.Ctx(ctx).Where("id=?", req.Id).One()
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, gerror.New("商品不存在")
	}

	// 检查商品是否有关联订单
	count, err := dao.OrderItem.Ctx(ctx).Where("product_id=?", req.Id).Count()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, gerror.New("该商品已有订单关联，无法删除")
	}

	// 尝试删除商品图片
	imageUrl := product["image"].String()
	if imageUrl != "" && (strings.HasPrefix(imageUrl, "http://") || strings.HasPrefix(imageUrl, "https://")) {
		// 忽略删除图片的错误，不影响商品删除
		_ = Minio().Delete(ctx, imageUrl)
	}

	// 删除商品
	_, err = dao.Product.Ctx(ctx).Where("id=?", req.Id).Delete()
	if err != nil {
		return nil, err
	}

	res.Success = true

	return res, nil
}

// 更新商品状态
func (s *localAdminProductService) UpdateStatus(ctx context.Context, req *api.AdminProductStatusUpdateReq) (res *api.AdminProductStatusUpdateRes, err error) {
	res = &api.AdminProductStatusUpdateRes{}

	// 检查商品是否存在
	product, err := dao.Product.Ctx(ctx).Where("id=?", req.Id).One()
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, gerror.New("商品不存在")
	}

	// 更新商品状态
	_, err = dao.Product.Ctx(ctx).Data(g.Map{
		"status": req.Status,
	}).Where("id=?", req.Id).Update()
	if err != nil {
		return nil, err
	}

	res.Success = true

	return res, nil
}
