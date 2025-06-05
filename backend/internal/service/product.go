package service

import (
	"context"
	"strings"

	v1 "backend/api/admin/v1"
	productv1 "backend/api/product/v1"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/utility/minio"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// ProductService 商品服务接口
type ProductService interface {
	// 管理端接口
	// List 获取商品列表（分页、筛选、模糊搜索）
	List(ctx context.Context, req *v1.AdminProductListReq) (list []productv1.UserProduct, total int, err error)

	// Create 创建商品
	Create(ctx context.Context, req *v1.AdminProductCreateReq) (*productv1.UserProduct, error)

	// Update 更新商品
	Update(ctx context.Context, req *v1.AdminProductUpdateReq) (*productv1.UserProduct, error)

	// Delete 删除商品
	Delete(ctx context.Context, id int64) error

	// 用户端接口
	// GetProductsByCategory 获取某分类下的商品列表
	GetProductsByCategory(ctx context.Context, categoryID int64, page, limit int, sortBy, sortOrder string) (list []productv1.UserProduct, total int, err error)

	// GetProductDetail 获取商品详情
	GetProductDetail(ctx context.Context, id int64) (*productv1.UserProductDetail, error)
}

// 商品服务实现
type productService struct{}

// 单例实例
var productServiceInstance = productService{}

// Product 获取商品服务实例
func Product() ProductService {
	return &productServiceInstance
}

// List 获取商品列表（分页、筛选、模糊搜索）
func (s *productService) List(ctx context.Context, req *v1.AdminProductListReq) (list []productv1.UserProduct, total int, err error) {
	// 构建查询条件
	m := dao.Products.Ctx(ctx)

	// 名称模糊搜索
	if req.Name != "" {
		m = m.WhereLike(dao.Products.Columns().Name, "%"+req.Name+"%")
	}

	// 分类ID筛选
	if req.CategoryID > 0 {
		m = m.Where(dao.Products.Columns().CategoryId, req.CategoryID)
	}

	// 是否上架筛选
	if req.IsActive != nil {
		isActive := 0
		if *req.IsActive {
			isActive = 1
		}
		m = m.Where(dao.Products.Columns().IsActive, isActive)
	}

	// 查询总数
	total, err = m.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页参数处理
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}

	// 查询数据
	var products []*entity.Products
	err = m.Page(page, limit).
		Order(dao.Products.Columns().Id + " DESC").
		Scan(&products)
	if err != nil {
		return nil, 0, err
	}

	// 转换为API响应格式
	list = make([]productv1.UserProduct, len(products))
	for i, p := range products {
		list[i] = productv1.UserProduct{
			ID:          int64(p.Id),
			Name:        p.Name,
			Price:       p.Price,
			ImageURL:    p.ImageUrl,
			Stock:       p.Stock,
			Description: p.Description,
			CategoryID:  int64(p.CategoryId),
			Status:      p.IsActive,
		}
	}

	return list, total, nil
}

// Create 创建商品
func (s *productService) Create(ctx context.Context, req *v1.AdminProductCreateReq) (*productv1.UserProduct, error) {
	// 检查分类是否存在
	categoryDao := &Category{}
	category, err := categoryDao.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, gerror.New("分类不存在")
	}

	// 处理是否上架字段
	isActive := 1 // 默认上架
	if req.IsActive != nil && !*req.IsActive {
		isActive = 0
	}

	// 创建商品数据
	product := &entity.Products{
		CategoryId:  int(req.CategoryID),
		Name:        req.Name,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageUrl:    req.ImageURL,
		Description: req.Description,
		IsActive:    isActive,
		CreatedAt:   gtime.Now(),
		UpdatedAt:   gtime.Now(),
	}

	// 插入数据
	result, err := dao.Products.Ctx(ctx).Insert(product)
	if err != nil {
		return nil, err
	}

	// 获取新插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 查询新创建的商品
	var newProduct *entity.Products
	err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, id).Scan(&newProduct)
	if err != nil {
		return nil, err
	}

	// 返回创建后的商品
	return &productv1.UserProduct{
		ID:          int64(newProduct.Id),
		Name:        newProduct.Name,
		Price:       newProduct.Price,
		ImageURL:    newProduct.ImageUrl,
		Stock:       newProduct.Stock,
		Description: newProduct.Description,
		CategoryID:  int64(newProduct.CategoryId),
		Status:      newProduct.IsActive,
	}, nil
}

// Update 更新商品
func (s *productService) Update(ctx context.Context, req *v1.AdminProductUpdateReq) (*productv1.UserProduct, error) {
	// 检查商品是否存在
	var product *entity.Products
	err := dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.ProductID).Scan(&product)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, gerror.New("商品不存在")
	}

	// 如果更新分类，检查分类是否存在
	if req.CategoryID > 0 {
		categoryDao := &Category{}
		category, err := categoryDao.GetByID(ctx, req.CategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, gerror.New("分类不存在")
		}
	}

	// 构建更新数据
	data := g.Map{}
	if req.Name != "" {
		data[dao.Products.Columns().Name] = req.Name
	}
	if req.CategoryID > 0 {
		data[dao.Products.Columns().CategoryId] = req.CategoryID
	}
	if req.Price > 0 {
		data[dao.Products.Columns().Price] = req.Price
	}
	if req.Stock >= 0 {
		data[dao.Products.Columns().Stock] = req.Stock
	}
	if req.ImageURL != "" {
		data[dao.Products.Columns().ImageUrl] = req.ImageURL
	}
	if req.Description != "" {
		data[dao.Products.Columns().Description] = req.Description
	}
	if req.IsActive != nil {
		isActive := 0
		if *req.IsActive {
			isActive = 1
		}
		data[dao.Products.Columns().IsActive] = isActive
	}
	data[dao.Products.Columns().UpdatedAt] = gtime.Now()

	// 更新数据
	_, err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.ProductID).Update(data)
	if err != nil {
		return nil, err
	}

	// 查询更新后的商品
	var updatedProduct *entity.Products
	err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.ProductID).Scan(&updatedProduct)
	if err != nil {
		return nil, err
	}

	// 返回更新后的商品
	return &productv1.UserProduct{
		ID:          int64(updatedProduct.Id),
		Name:        updatedProduct.Name,
		Price:       updatedProduct.Price,
		ImageURL:    updatedProduct.ImageUrl,
		Stock:       updatedProduct.Stock,
		Description: updatedProduct.Description,
		CategoryID:  int64(updatedProduct.CategoryId),
		Status:      updatedProduct.IsActive,
	}, nil
}

// Delete 删除商品
func (s *productService) Delete(ctx context.Context, id int64) error {
	// 检查商品是否存在
	var product *entity.Products
	err := dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, id).Scan(&product)
	if err != nil {
		return err
	}
	if product == nil {
		return gerror.New("商品不存在")
	}

	// 检查是否有订单关联此商品
	count, err := dao.OrderItems.Ctx(ctx).Where("product_id", id).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("该商品已有订单关联，无法删除")
	}

	imageURL := product.ImageUrl

	// 删除商品记录
	_, err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, id).Delete()
	if err != nil {
		return err
	}

	// 如果有图片，则删除图片文件
	if imageURL != "" {
		// 从URL中提取文件路径
		// 图片URL格式为: http://domain/api/v1/file/directory/filename.ext 或 http://domain/file/directory/filename.ext
		// 需要提取的是 directory/filename.ext 部分
		var objectName string

		if strings.Contains(imageURL, "/api/v1/file/") {
			objectName = strings.TrimPrefix(imageURL, strings.Split(imageURL, "/api/v1/file/")[0]+"/api/v1/file/")
		} else if strings.Contains(imageURL, "/file/") {
			objectName = strings.TrimPrefix(imageURL, strings.Split(imageURL, "/file/")[0]+"/file/")
		}

		if objectName != "" {
			// 获取MinIO客户端
			minioClient := minio.GetClient()
			if minioClient != nil {
				// 删除文件，忽略可能的错误
				deleteErr := minioClient.DeleteFile(ctx, objectName)
				if deleteErr != nil {
					g.Log().Warning(ctx, "删除商品图片失败: ", deleteErr, ", 文件路径: ", objectName)
				} else {
					g.Log().Info(ctx, "成功删除商品图片: ", objectName)
				}
			}
		}
	}

	return nil
}

// GetProductsByCategory 获取某分类下的商品列表（用户端）
func (s *productService) GetProductsByCategory(ctx context.Context, categoryID int64, page, limit int, sortBy, sortOrder string) (list []productv1.UserProduct, total int, err error) {
	// 构建查询条件
	m := dao.Products.Ctx(ctx).
		Where(dao.Products.Columns().CategoryId, categoryID).
		Where(dao.Products.Columns().IsActive, 1) // 只查询已上架商品

	// 查询总数
	total, err = m.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页参数处理
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// 处理排序
	orderBy := dao.Products.Columns().Id + " DESC" // 默认按ID倒序排序

	if sortBy != "" {
		// 映射API字段名到数据库字段名
		var dbField string
		switch sortBy {
		case "price":
			dbField = dao.Products.Columns().Price
			// 暂时不支持销量排序，因为数据库没有该字段
			// case "sales_count":
			//	dbField = dao.Products.Columns().SalesCount
		}

		if dbField != "" {
			if sortOrder == "asc" {
				orderBy = dbField + " ASC"
			} else {
				orderBy = dbField + " DESC"
			}
		}
	}

	// 查询数据
	var products []*entity.Products
	err = m.Page(page, limit).
		Order(orderBy).
		Scan(&products)
	if err != nil {
		return nil, 0, err
	}

	// 转换为API响应格式
	list = make([]productv1.UserProduct, len(products))
	for i, p := range products {
		if err = gconv.Struct(p, &list[i]); err != nil {
			return nil, 0, err
		}
		list[i].ID = int64(p.Id)
		list[i].ImageURL = p.ImageUrl
		list[i].CategoryID = int64(p.CategoryId)
		list[i].Status = p.IsActive
		// 暂时不设置销量
		// list[i].SalesCount = p.SalesCount
	}

	return list, total, nil
}

// GetProductDetail 获取商品详情（用户端）
func (s *productService) GetProductDetail(ctx context.Context, id int64) (*productv1.UserProductDetail, error) {
	var product *entity.Products
	err := dao.Products.Ctx(ctx).
		Where(dao.Products.Columns().Id, id).
		Where(dao.Products.Columns().IsActive, 1). // 只查询已上架商品
		Scan(&product)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, gerror.New("商品不存在或已下架")
	}

	detail := &productv1.UserProductDetail{}
	detail.UserProduct = productv1.UserProduct{
		ID:          int64(product.Id),
		Name:        product.Name,
		Price:       product.Price,
		ImageURL:    product.ImageUrl,
		Stock:       product.Stock,
		Description: product.Description,
		CategoryID:  int64(product.CategoryId),
		Status:      product.IsActive,
	}

	return detail, nil
}
