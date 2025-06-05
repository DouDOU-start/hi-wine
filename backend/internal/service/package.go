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
		adminPackage.DurationHours = item.DurationMinutes / 60 // 转换分钟为小时
		adminPackage.IsActive = item.IsActive == 1

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
		ID:            int64(packageInfo.Id),
		Name:          packageInfo.Name,
		Description:   packageInfo.Description,
		Price:         packageInfo.Price,
		DurationHours: packageInfo.DurationMinutes / 60, // 转换分钟为小时
		IsActive:      packageInfo.IsActive == 1,
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

		// 转换小时为分钟
		durationMinutes := req.DurationHours * 60

		packageData := entity.DrinkAllYouCanPackages{
			Name:            req.Name,
			Description:     req.Description,
			Price:           req.Price,
			DurationMinutes: durationMinutes,
			IsActive:        isActive,
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
			ID:            int64(createdPackage.Id),
			Name:          createdPackage.Name,
			Description:   createdPackage.Description,
			Price:         createdPackage.Price,
			DurationHours: createdPackage.DurationMinutes / 60, // 转换分钟为小时
			IsActive:      createdPackage.IsActive == 1,
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

		if req.DurationHours > 0 {
			updateData["duration_minutes"] = req.DurationHours * 60 // 转换小时为分钟
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
			ID:            int64(updatedPackage.Id),
			Name:          updatedPackage.Name,
			Description:   updatedPackage.Description,
			Price:         updatedPackage.Price,
			DurationHours: updatedPackage.DurationMinutes / 60, // 转换分钟为小时
			IsActive:      updatedPackage.IsActive == 1,
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
		product.Status = p.IsActive
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
		product.Status = p.IsActive
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
