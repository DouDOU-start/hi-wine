package service

import (
	v1 "backend/api/admin/v1"
	productv1 "backend/api/product/v1"
	"backend/internal/model/entity"
	"backend/internal/utility"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PackageService 套餐服务接口
type PackageService interface {
	// GetPackageList 获取套餐列表（支持分页和筛选）
	GetPackageList(ctx context.Context, req *v1.AdminPackageListReq) (list []v1.AdminPackage, total int, err error)

	// GetPackageDetail 获取套餐详情
	GetPackageDetail(ctx context.Context, id int64) (detail *v1.AdminPackage, err error)

	// CreatePackage 创建套餐
	CreatePackage(ctx context.Context, req *v1.AdminPackageCreateReq) (detail *v1.AdminPackage, err error)

	// UpdatePackage 更新套餐
	UpdatePackage(ctx context.Context, req *v1.AdminPackageUpdateReq) (detail *v1.AdminPackage, err error)

	// DeletePackage 删除套餐
	DeletePackage(ctx context.Context, id int64) error

	// GetPackageProducts 获取套餐包含的商品列表
	GetPackageProducts(ctx context.Context, packageID int64) (list []productv1.UserProduct, err error)

	// AddPackageProducts 为套餐添加商品
	AddPackageProducts(ctx context.Context, packageID int64, productIDs []int64) error

	// RemovePackageProduct 从套餐中移除商品
	RemovePackageProduct(ctx context.Context, packageID int64, productID int64) error

	// GetAvailableProducts 获取可添加到套餐的商品列表
	GetAvailableProducts(ctx context.Context, packageID int64, keyword string, page, limit int) (list []productv1.UserProduct, total int, err error)

	// BatchRemovePackageProducts 批量从套餐中移除商品
	BatchRemovePackageProducts(ctx context.Context, packageID int64, productIDs []int64) error

	// GetPackageStats 获取套餐使用统计
	GetPackageStats(ctx context.Context, packageID int64) (*v1.AdminPackageStats, error)

	// CreateUserOrderWithPackage 创建用户套餐购买订单
	CreateUserOrderWithPackage(ctx context.Context, userId int64, packageId int64) (orderId int64, err error)

	// GetPackageWithProducts 获取带商品列表的套餐详情
	GetPackageWithProducts(ctx context.Context, packageID int64) (*v1.AdminPackageWithProducts, error)

	// GetPackageFullDetail 获取套餐详细信息（包含基本信息、统计信息和商品列表）
	GetPackageFullDetail(ctx context.Context, packageID int64) (*v1.AdminPackageFullDetail, error)
}

// packageService 套餐服务实现
type packageService struct{}

// 确保 packageService 实现了 PackageService 接口
var _ PackageService = (*packageService)(nil)

// 单例实例
var packageServiceInstance = packageService{}

// Package 获取套餐服务实例
func Package() PackageService {
	return &packageServiceInstance
}

// GetPackageList 获取套餐列表（支持分页和筛选）
func (s *packageService) GetPackageList(ctx context.Context, req *v1.AdminPackageListReq) (list []v1.AdminPackage, total int, err error) {
	// 构建查询条件
	m := g.DB().Model("drink_all_you_can_packages").As("p")

	// 添加筛选条件
	if req.Name != "" {
		m = m.WhereLike("p.name", "%"+req.Name+"%")
	}

	// 查询总数
	count, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// 查询数据
	var packages []entity.DrinkAllYouCanPackages
	err = m.Order("p.id DESC").
		Limit((page-1)*limit, limit).
		Scan(&packages)
	if err != nil {
		return nil, 0, err
	}

	// 转换为API响应格式
	list = make([]v1.AdminPackage, 0, len(packages))
	for _, item := range packages {
		var adminPackage v1.AdminPackage
		adminPackage.ID = int64(item.Id)
		adminPackage.Name = item.Name
		adminPackage.Description = item.Description
		adminPackage.Price = item.Price
		adminPackage.DurationMinutes = item.DurationMinutes // 直接使用分钟
		adminPackage.IsActive = item.IsActive == 1
		adminPackage.CreatedAt = item.CreatedAt.Format("2006-01-02 15:04:05")
		adminPackage.UpdatedAt = item.UpdatedAt.Format("2006-01-02 15:04:05")

		list = append(list, adminPackage)
	}

	return list, count, nil
}

// GetPackageDetail 获取套餐详情
func (s *packageService) GetPackageDetail(ctx context.Context, id int64) (detail *v1.AdminPackage, err error) {
	// 查询套餐
	var packageInfo entity.DrinkAllYouCanPackages
	err = g.DB().Model("drink_all_you_can_packages").Where("id", id).Scan(&packageInfo)
	if err != nil {
		return nil, err
	}
	if packageInfo.Id == 0 {
		return nil, gerror.New("套餐不存在")
	}

	// 转换为API响应格式
	detail = &v1.AdminPackage{
		ID:              int64(packageInfo.Id),
		Name:            packageInfo.Name,
		Description:     packageInfo.Description,
		Price:           packageInfo.Price,
		DurationMinutes: packageInfo.DurationMinutes,
		IsActive:        packageInfo.IsActive == 1,
		CreatedAt:       packageInfo.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       packageInfo.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return detail, nil
}

// CreatePackage 创建套餐
func (s *packageService) CreatePackage(ctx context.Context, req *v1.AdminPackageCreateReq) (detail *v1.AdminPackage, err error) {
	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 创建套餐记录
		isActive := 1
		if req.IsActive != nil && !*req.IsActive {
			isActive = 0
		}

		now := gtime.Now()
		packageData := entity.DrinkAllYouCanPackages{
			Name:            req.Name,
			Description:     req.Description,
			Price:           req.Price,
			DurationMinutes: req.DurationMinutes,
			IsActive:        isActive,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		// 插入记录
		result, err := tx.Model("drink_all_you_can_packages").Data(packageData).Insert()
		if err != nil {
			return err
		}

		// 获取插入的ID
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// 查询创建的记录
		var createdPackage entity.DrinkAllYouCanPackages
		err = tx.Model("drink_all_you_can_packages").Where("id", lastInsertID).Scan(&createdPackage)
		if err != nil {
			return err
		}

		// 转换为API响应格式
		detail = &v1.AdminPackage{
			ID:              int64(createdPackage.Id),
			Name:            createdPackage.Name,
			Description:     createdPackage.Description,
			Price:           createdPackage.Price,
			DurationMinutes: createdPackage.DurationMinutes,
			IsActive:        createdPackage.IsActive == 1,
			CreatedAt:       createdPackage.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       createdPackage.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return detail, nil
}

// UpdatePackage 更新套餐
func (s *packageService) UpdatePackage(ctx context.Context, req *v1.AdminPackageUpdateReq) (detail *v1.AdminPackage, err error) {
	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询套餐是否存在
		var count int
		count, err = tx.Model("drink_all_you_can_packages").Where("id", req.PackageID).Count()
		if err != nil {
			return err
		}
		if count == 0 {
			return gerror.New("套餐不存在")
		}

		// 构建更新数据
		updateData := g.Map{}

		if req.Name != "" {
			updateData["name"] = req.Name
		}

		if req.Description != "" {
			updateData["description"] = req.Description
		}

		if req.Price > 0 {
			updateData["price"] = req.Price
		}

		if req.DurationMinutes > 0 {
			updateData["duration_minutes"] = req.DurationMinutes
		}

		if req.IsActive != nil {
			if *req.IsActive {
				updateData["is_active"] = 1
			} else {
				updateData["is_active"] = 0
			}
		}

		// 更新时间
		updateData["updated_at"] = gtime.Now()

		// 执行更新
		_, err = tx.Model("drink_all_you_can_packages").
			Where("id", req.PackageID).
			Data(updateData).
			Update()
		if err != nil {
			return err
		}

		// 查询更新后的记录
		var updatedPackage entity.DrinkAllYouCanPackages
		err = tx.Model("drink_all_you_can_packages").Where("id", req.PackageID).Scan(&updatedPackage)
		if err != nil {
			return err
		}

		// 转换为API响应格式
		detail = &v1.AdminPackage{
			ID:              int64(updatedPackage.Id),
			Name:            updatedPackage.Name,
			Description:     updatedPackage.Description,
			Price:           updatedPackage.Price,
			DurationMinutes: updatedPackage.DurationMinutes,
			IsActive:        updatedPackage.IsActive == 1,
			CreatedAt:       updatedPackage.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       updatedPackage.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return detail, nil
}

// DeletePackage 删除套餐
func (s *packageService) DeletePackage(ctx context.Context, id int64) error {
	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询套餐是否存在
		var count int
		var err error
		count, err = tx.Model("drink_all_you_can_packages").Where("id", id).Count()
		if err != nil {
			return err
		}
		if count == 0 {
			return gerror.New("套餐不存在")
		}

		// 检查是否有用户正在使用此套餐
		count, err = tx.Model("user_packages").
			Where("package_id", id).
			Where("status", "active").
			Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return gerror.New("该套餐有用户正在使用，无法删除")
		}

		// 删除套餐与商品的关联
		_, err = tx.Model("package_products").Where("package_id", id).Delete()
		if err != nil {
			return err
		}

		// 删除套餐
		_, err = tx.Model("drink_all_you_can_packages").Where("id", id).Delete()
		if err != nil {
			return err
		}

		return nil
	})
}

// GetPackageProducts 获取套餐包含的商品列表
func (s *packageService) GetPackageProducts(ctx context.Context, packageID int64) (list []productv1.UserProduct, err error) {
	// 查询套餐是否存在
	count, err := g.DB().Model("drink_all_you_can_packages").Where("id", packageID).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.New("套餐不存在")
	}

	// 查询套餐包含的商品
	var products []entity.Products
	err = g.DB().Model("package_products pp").
		LeftJoin("products p", "p.id = pp.product_id").
		Where("pp.package_id", packageID).
		Fields("p.*").
		Scan(&products)
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	list = make([]productv1.UserProduct, 0, len(products))
	for _, p := range products {
		var product productv1.UserProduct
		product.ID = int64(p.Id)
		product.Name = p.Name
		product.Price = p.Price
		product.ImageURL = p.ImageUrl
		product.Stock = p.Stock
		product.Description = p.Description
		product.CategoryID = int64(p.CategoryId)
		product.IsActive = p.IsActive
		list = append(list, product)
	}

	return list, nil
}

// AddPackageProducts 为套餐添加商品
func (s *packageService) AddPackageProducts(ctx context.Context, packageID int64, productIDs []int64) error {
	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询套餐是否存在
		var count int
		var err error
		count, err = tx.Model("drink_all_you_can_packages").Where("id", packageID).Count()
		if err != nil {
			return err
		}
		if count == 0 {
			return gerror.New("套餐不存在")
		}

		// 检查商品是否都存在
		for _, productID := range productIDs {
			count, err = tx.Model("products").Where("id", productID).Count()
			if err != nil {
				return err
			}
			if count == 0 {
				return gerror.Newf("商品ID %d 不存在", productID)
			}

			// 检查关联是否已存在
			count, err = tx.Model("package_products").
				Where("package_id", packageID).
				Where("product_id", productID).
				Count()
			if err != nil {
				return err
			}

			// 如果关联不存在，则添加
			if count == 0 {
				_, err = tx.Model("package_products").Insert(g.Map{
					"package_id": packageID,
					"product_id": productID,
					"created_at": gtime.Now(),
				})
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// RemovePackageProduct 从套餐中移除商品
func (s *packageService) RemovePackageProduct(ctx context.Context, packageID int64, productID int64) error {
	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询套餐是否存在
		var count int
		var err error
		count, err = tx.Model("drink_all_you_can_packages").Where("id", packageID).Count()
		if err != nil {
			return err
		}
		if count == 0 {
			return gerror.New("套餐不存在")
		}

		// 查询商品是否存在
		count, err = tx.Model("products").Where("id", productID).Count()
		if err != nil {
			return err
		}
		if count == 0 {
			return gerror.New("商品不存在")
		}

		// 删除关联
		_, err = tx.Model("package_products").
			Where("package_id", packageID).
			Where("product_id", productID).
			Delete()
		if err != nil {
			return err
		}

		return nil
	})
}

// GetAvailableProducts 获取可添加到套餐的商品列表（未添加到该套餐的商品）
func (s *packageService) GetAvailableProducts(ctx context.Context, packageID int64, keyword string, page, limit int) (list []productv1.UserProduct, total int, err error) {
	// 查询套餐是否存在
	count, err := g.DB().Model("drink_all_you_can_packages").Where("id", packageID).Count()
	if err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return nil, 0, gerror.New("套餐不存在")
	}

	// 构建子查询
	subQuery := fmt.Sprintf("(SELECT product_id FROM package_products WHERE package_id = %d)", packageID)

	// 构建查询条件
	m := g.DB().Model("products p").
		LeftJoin(subQuery+" AS pp", "p.id = pp.product_id").
		Where("pp.product_id IS NULL"). // 只查询未添加到该套餐的商品
		Where("p.is_active", 1)         // 只查询已上架的商品

	// 添加关键字搜索
	if keyword != "" {
		m = m.WhereLike("p.name", "%"+keyword+"%")
	}

	// 查询总数
	total, err = m.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页参数
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// 查询数据
	var products []entity.Products
	err = m.Fields("p.*").
		Order("p.id DESC").
		Limit((page-1)*limit, limit).
		Scan(&products)
	if err != nil {
		return nil, 0, err
	}

	// 转换为API响应格式
	list = make([]productv1.UserProduct, 0, len(products))
	for _, p := range products {
		var product productv1.UserProduct
		product.ID = int64(p.Id)
		product.Name = p.Name
		product.Price = p.Price
		product.ImageURL = p.ImageUrl
		product.Stock = p.Stock
		product.Description = p.Description
		product.CategoryID = int64(p.CategoryId)
		product.IsActive = p.IsActive
		list = append(list, product)
	}

	return list, total, nil
}

// BatchRemovePackageProducts 批量从套餐中移除商品
func (s *packageService) BatchRemovePackageProducts(ctx context.Context, packageID int64, productIDs []int64) error {
	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询套餐是否存在
		var count int
		var err error
		count, err = tx.Model("drink_all_you_can_packages").Where("id", packageID).Count()
		if err != nil {
			return err
		}
		if count == 0 {
			return gerror.New("套餐不存在")
		}

		// 批量删除关联
		_, err = tx.Model("package_products").
			Where("package_id", packageID).
			WhereIn("product_id", productIDs).
			Delete()
		if err != nil {
			return err
		}

		return nil
	})
}

// GetPackageStats 获取套餐使用统计
func (s *packageService) GetPackageStats(ctx context.Context, packageID int64) (*v1.AdminPackageStats, error) {
	// 查询套餐是否存在
	count, err := g.DB().Model("drink_all_you_can_packages").Where("id", packageID).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.New("套餐不存在")
	}

	stats := &v1.AdminPackageStats{}

	// 总销售数量
	totalSales, err := g.DB().Model("user_packages").
		Where("package_id", packageID).
		Count()
	if err != nil {
		return nil, err
	}
	stats.TotalSales = totalSales

	// 当前活跃用户数
	activeUsers, err := g.DB().Model("user_packages").
		Where("package_id", packageID).
		Where("status", "active").
		Count()
	if err != nil {
		return nil, err
	}
	stats.ActiveUsers = activeUsers

	// 总收入
	var packageInfo entity.DrinkAllYouCanPackages
	err = g.DB().Model("drink_all_you_can_packages").Where("id", packageID).Scan(&packageInfo)
	if err != nil {
		return nil, err
	}
	stats.TotalRevenue = packageInfo.Price * float64(totalSales)

	// 平均每日使用量（近30天）
	dailyUsageCount, err := g.DB().Model("user_packages").
		Where("package_id", packageID).
		Where("created_at >= ?", gtime.Now().AddDate(0, 0, -30).Format("Y-m-d")).
		Count()
	if err != nil {
		return nil, err
	}
	stats.AvgDailyUsage = float64(dailyUsageCount) / 30.0

	// 包含商品数量
	productsCount, err := g.DB().Model("package_products").
		Where("package_id", packageID).
		Count()
	if err != nil {
		return nil, err
	}
	stats.ProductsCount = productsCount

	return stats, nil
}

// CreateUserOrderWithPackage 创建用户套餐购买订单
func (s *packageService) CreateUserOrderWithPackage(ctx context.Context, userId int64, packageId int64) (orderId int64, err error) {
	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 查询套餐信息
		var packageInfo entity.DrinkAllYouCanPackages
		err = tx.Model("drink_all_you_can_packages").Where("id", packageId).Where("is_active", 1).Scan(&packageInfo)
		if err != nil {
			return err
		}
		if packageInfo.Id == 0 {
			return gerror.New("套餐不存在或已下架")
		}

		// 2. 验证用户是否存在
		var userCount int
		userCount, err = tx.Model("users").Where("id", userId).Count()
		if err != nil {
			return err
		}
		if userCount == 0 {
			return gerror.New("用户不存在")
		}

		// 3. 创建订单
		now := gtime.Now()
		orderSN := utility.GenerateOrderSN()

		// 3.1 创建订单主表记录 - 直接插入SQL以绕过外键约束
		_, err = tx.Exec(
			"INSERT INTO orders(order_sn, user_id, table_qrcode_id, total_amount, payment_status, order_status, payment_method, transaction_id, created_at, updated_at) "+
				"VALUES(?, ?, NULL, ?, 'pending', 'new', '', '', ?, ?)",
			orderSN, userId, packageInfo.Price, now, now,
		)
		if err != nil {
			return err
		}

		// 3.2 获取最后插入的订单ID
		var result gdb.Value
		result, err = tx.GetValue("SELECT LAST_INSERT_ID()")
		if err != nil {
			return err
		}
		orderId = result.Int64()

		// 4. 创建用户套餐记录（默认为pending状态，等待支付成功后激活）
		_, err = UserPackage().CreateUserPackageAfterPurchase(ctx, userId, packageId, orderId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return orderId, nil
}

// GetPackageWithProducts 获取带商品列表的套餐详情
func (s *packageService) GetPackageWithProducts(ctx context.Context, packageID int64) (*v1.AdminPackageWithProducts, error) {
	// 1. 查询套餐基本信息
	var packageInfo entity.DrinkAllYouCanPackages
	err := g.DB().Model("drink_all_you_can_packages").Where("id", packageID).Scan(&packageInfo)
	if err != nil {
		g.Log().Errorf(ctx, "查询套餐信息失败, ID:%d, 错误:%v", packageID, err)
		return nil, err
	}
	if packageInfo.Id == 0 {
		return nil, gerror.New("套餐不存在")
	}

	// 2. 创建返回结构
	detail := &v1.AdminPackageWithProducts{
		ID:              int64(packageInfo.Id),
		Name:            packageInfo.Name,
		Description:     packageInfo.Description,
		Price:           packageInfo.Price,
		DurationMinutes: packageInfo.DurationMinutes,
		DurationDays:    packageInfo.DurationMinutes / (60 * 24), // 计算天数
		IsActive:        packageInfo.IsActive == 1,
		CreatedAt:       packageInfo.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       packageInfo.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// 3. 获取套餐包含的商品列表
	products, err := s.GetPackageProducts(ctx, packageID)
	if err != nil {
		g.Log().Errorf(ctx, "获取套餐商品列表失败, ID:%d, 错误:%v", packageID, err)
		// 不返回错误，保持空商品列表
		detail.Products = make([]productv1.UserProduct, 0)
		detail.ProductsCount = 0
	} else {
		detail.Products = products
		detail.ProductsCount = len(products)
	}

	return detail, nil
}

// GetPackageFullDetail 获取套餐详细信息（包含基本信息、统计信息和商品列表）
func (s *packageService) GetPackageFullDetail(ctx context.Context, packageID int64) (*v1.AdminPackageFullDetail, error) {
	detail := &v1.AdminPackageFullDetail{}

	// 1. 查询套餐基本信息
	var packageInfo entity.DrinkAllYouCanPackages
	err := g.DB().Model("drink_all_you_can_packages").Where("id", packageID).Scan(&packageInfo)
	if err != nil {
		g.Log().Errorf(ctx, "查询套餐信息失败, ID:%d, 错误:%v", packageID, err)
		return nil, err
	}
	if packageInfo.Id == 0 {
		return nil, gerror.New("套餐不存在")
	}

	// 设置基本信息
	detail.ID = int64(packageInfo.Id)
	detail.Name = packageInfo.Name
	detail.Description = packageInfo.Description
	detail.Price = packageInfo.Price
	detail.DurationMinutes = packageInfo.DurationMinutes
	detail.DurationDays = packageInfo.DurationMinutes / (60 * 24) // 计算天数
	detail.IsActive = packageInfo.IsActive == 1
	detail.CreatedAt = packageInfo.CreatedAt.Format("Y-m-d H:i:s")
	detail.UpdatedAt = packageInfo.UpdatedAt.Format("Y-m-d H:i:s")

	// 2. 获取套餐统计信息
	stats, err := s.GetPackageStats(ctx, packageID)
	if err != nil {
		g.Log().Errorf(ctx, "获取套餐统计信息失败, ID:%d, 错误:%v", packageID, err)
		// 不返回错误，继续获取其他信息
	} else {
		detail.Stats = *stats
	}

	// 3. 获取套餐包含的商品列表
	products, err := s.GetPackageProducts(ctx, packageID)
	if err != nil {
		g.Log().Errorf(ctx, "获取套餐商品列表失败, ID:%d, 错误:%v", packageID, err)
		// 不返回错误，继续获取其他信息
	} else {
		detail.Products = products
	}

	// 4. 获取最近购买记录
	type PurchaseRecord struct {
		ID           int64       `json:"id"`
		UserID       int64       `json:"user_id"`
		UserName     string      `json:"user_name"`
		OrderID      int64       `json:"order_id"`
		OrderSN      string      `json:"order_sn"`
		PurchaseTime *gtime.Time `json:"purchase_time"`
		Status       string      `json:"status"`
	}

	var recentPurchases []PurchaseRecord
	err = g.DB().Model("user_packages up").
		LeftJoin("users u", "u.id = up.user_id").
		LeftJoin("orders o", "o.id = up.order_id").
		Fields("up.id, up.user_id, u.nickname as user_name, up.order_id, o.order_sn, o.created_at as purchase_time, up.status").
		Where("up.package_id", packageID).
		Order("up.created_at DESC").
		Limit(10).
		Scan(&recentPurchases)

	if err != nil {
		g.Log().Errorf(ctx, "获取套餐最近购买记录失败, ID:%d, 错误:%v", packageID, err)
		// 不返回错误，继续处理
	} else {
		// 转换为API响应格式
		for _, record := range recentPurchases {
			purchase := struct {
				ID           int64  `json:"id"`
				UserID       int64  `json:"user_id"`
				UserName     string `json:"user_name"`
				OrderID      int64  `json:"order_id"`
				OrderSN      string `json:"order_sn"`
				PurchaseTime string `json:"purchase_time"`
				Status       string `json:"status"`
			}{
				ID:       record.ID,
				UserID:   record.UserID,
				UserName: record.UserName,
				OrderID:  record.OrderID,
				OrderSN:  record.OrderSN,
				Status:   record.Status,
			}

			if record.PurchaseTime != nil {
				purchase.PurchaseTime = record.PurchaseTime.Format("Y-m-d H:i:s")
			}

			detail.RecentPurchases = append(detail.RecentPurchases, purchase)
		}
	}

	return detail, nil
}
