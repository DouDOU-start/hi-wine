package service

import (
	v1 "backend/api/admin/v1"
	"backend/internal/model/entity"
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// UserPackageService 用户套餐服务接口
type UserPackageService interface {
	// GetUserPackageList 获取用户套餐列表（支持分页和筛选）
	GetUserPackageList(ctx context.Context, req *v1.AdminUserPackageListReq) (list []v1.AdminUserPackage, total int, err error)

	// GetUserPackageDetail 获取用户套餐详情
	GetUserPackageDetail(ctx context.Context, id int64) (detail *v1.AdminUserPackage, err error)

	// CreateUserPackage 创建用户套餐
	CreateUserPackage(ctx context.Context, req *v1.AdminUserPackageCreateReq) (detail *v1.AdminUserPackage, err error)

	// UpdateUserPackageStatus 更新用户套餐状态
	UpdateUserPackageStatus(ctx context.Context, req *v1.AdminUserPackageUpdateStatusReq) (detail *v1.AdminUserPackage, err error)

	// GetUserActivePackages 获取用户有效套餐
	GetUserActivePackages(ctx context.Context, userID int64) (list []v1.AdminUserPackage, err error)

	// ActivateUserPackage 激活用户套餐（首次使用时）
	ActivateUserPackage(ctx context.Context, userPackageID int64) error

	// GetUserPendingPackages 获取用户待激活的套餐列表
	GetUserPendingPackages(ctx context.Context, userID int64) (list []v1.AdminUserPackage, err error)

	// CreateUserPackageAfterPurchase 用户购买套餐后创建待支付的套餐记录
	CreateUserPackageAfterPurchase(ctx context.Context, userID, packageID, orderID int64) (userPackageID int64, err error)

	// ActivateUserPackageAfterPayment 支付成功后激活用户套餐
	ActivateUserPackageAfterPayment(ctx context.Context, orderID int64) error

	// CheckAndUpdatePackageStatus 检查并更新单个用户套餐的状态
	CheckAndUpdatePackageStatus(ctx context.Context, userPackage *entity.UserPackages) error
}

// userPackageService 用户套餐服务实现
type userPackageService struct{}

// 确保 userPackageService 实现了 UserPackageService 接口
var _ UserPackageService = (*userPackageService)(nil)

// 单例实例
var userPackageServiceInstance = userPackageService{}

// UserPackage 获取用户套餐服务实例
func UserPackage() UserPackageService {
	return &userPackageServiceInstance
}

// CheckAndUpdatePackageStatus 检查并更新单个用户套餐的状态
func (s *userPackageService) CheckAndUpdatePackageStatus(ctx context.Context, userPackage *entity.UserPackages) error {
	// 如果套餐状态已经是过期状态，则无需更新
	if userPackage.Status == "expired" {
		return nil
	}

	// 如果套餐状态是激活状态，且结束时间已过期，则更新为过期状态
	if userPackage.Status == "active" && userPackage.EndTime != nil && userPackage.EndTime.Before(gtime.Now()) {
		_, err := g.DB().Model("user_packages").
			Data(g.Map{
				"status":     "expired",
				"updated_at": gtime.Now(),
			}).
			Where("id", userPackage.Id).
			Update()

		if err != nil {
			return err
		}

		// 更新当前对象的状态
		userPackage.Status = "expired"
		userPackage.UpdatedAt = gtime.Now()
	}

	return nil
}

// GetUserPackageList 获取用户套餐列表（支持分页和筛选）
func (s *userPackageService) GetUserPackageList(ctx context.Context, req *v1.AdminUserPackageListReq) (list []v1.AdminUserPackage, total int, err error) {
	// 构建查询条件
	m := g.DB().Model("user_packages").As("up")

	// 关联用户表和套餐表
	m = m.LeftJoin("users u", "u.id = up.user_id")
	m = m.LeftJoin("drink_all_you_can_packages p", "p.id = up.package_id")

	// 添加筛选条件
	if req.UserID > 0 {
		m = m.Where("up.user_id", req.UserID)
	}
	if req.PackageID > 0 {
		m = m.Where("up.package_id", req.PackageID)
	}
	if req.Status != "" {
		m = m.Where("up.status", req.Status)
	}
	if req.StartDate != "" {
		m = m.WhereGTE("up.created_at", req.StartDate)
	}
	if req.EndDate != "" {
		m = m.WhereLTE("up.created_at", req.EndDate+" 23:59:59")
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
	var userPackages []entity.UserPackages
	err = m.Fields("up.*").
		Order("up.id DESC").
		Limit((page-1)*limit, limit).
		Scan(&userPackages)
	if err != nil {
		return nil, 0, err
	}

	// 实时检查并更新套餐状态
	for i := range userPackages {
		if err = s.CheckAndUpdatePackageStatus(ctx, &userPackages[i]); err != nil {
			g.Log().Errorf(ctx, "检查更新套餐状态失败: %v", err)
		}
	}

	// 转换为API响应格式
	list = make([]v1.AdminUserPackage, 0, len(userPackages))
	for _, item := range userPackages {
		var adminUserPackage v1.AdminUserPackage
		if err = gconv.Struct(item, &adminUserPackage); err != nil {
			return nil, 0, err
		}

		// 格式化时间
		if item.StartTime != nil {
			adminUserPackage.StartTime = item.StartTime.Format("Y-m-d H:i:s")
		}
		if item.EndTime != nil {
			adminUserPackage.EndTime = item.EndTime.Format("Y-m-d H:i:s")
		}
		if item.CreatedAt != nil {
			adminUserPackage.CreatedAt = item.CreatedAt.Format("Y-m-d H:i:s")
		}
		if item.UpdatedAt != nil {
			adminUserPackage.UpdatedAt = item.UpdatedAt.Format("Y-m-d H:i:s")
		}

		list = append(list, adminUserPackage)
	}

	return list, count, nil
}

// GetUserPackageDetail 获取用户套餐详情
func (s *userPackageService) GetUserPackageDetail(ctx context.Context, id int64) (detail *v1.AdminUserPackage, err error) {
	// 查询用户套餐
	var userPackage entity.UserPackages
	err = g.DB().Model("user_packages").Where("id", id).Scan(&userPackage)
	if err != nil {
		return nil, err
	}
	if userPackage.Id == 0 {
		return nil, gerror.New("用户套餐不存在")
	}

	// 实时检查并更新套餐状态
	if err = s.CheckAndUpdatePackageStatus(ctx, &userPackage); err != nil {
		g.Log().Errorf(ctx, "检查更新套餐状态失败: %v", err)
	}

	// 转换为API响应格式
	detail = &v1.AdminUserPackage{}
	if err = gconv.Struct(userPackage, detail); err != nil {
		return nil, err
	}

	// 格式化时间
	if userPackage.StartTime != nil {
		detail.StartTime = userPackage.StartTime.Format("Y-m-d H:i:s")
	}
	if userPackage.EndTime != nil {
		detail.EndTime = userPackage.EndTime.Format("Y-m-d H:i:s")
	}
	if userPackage.CreatedAt != nil {
		detail.CreatedAt = userPackage.CreatedAt.Format("Y-m-d H:i:s")
	}
	if userPackage.UpdatedAt != nil {
		detail.UpdatedAt = userPackage.UpdatedAt.Format("Y-m-d H:i:s")
	}

	return detail, nil
}

// CreateUserPackage 创建用户套餐
func (s *userPackageService) CreateUserPackage(ctx context.Context, req *v1.AdminUserPackageCreateReq) (detail *v1.AdminUserPackage, err error) {
	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询套餐信息，获取套餐时长
		var packageInfo entity.DrinkAllYouCanPackages
		err = tx.Model("drink_all_you_can_packages").Where("id", req.PackageID).Scan(&packageInfo)
		if err != nil {
			return err
		}
		if packageInfo.Id == 0 {
			return gerror.New("套餐不存在")
		}

		// 验证用户是否存在
		var userCount int
		userCount, err = tx.Model("users").Where("id", req.UserID).Count()
		if err != nil {
			return err
		}
		if userCount == 0 {
			return gerror.New("用户不存在")
		}

		// 创建用户套餐记录
		userPackage := entity.UserPackages{
			UserId:    int(req.UserID),
			PackageId: int(req.PackageID),
			Status:    req.Status,
		}

		// 设置关联订单ID（如果有）
		if req.OrderID > 0 {
			userPackage.OrderId = int(req.OrderID)
		}

		// 仅当套餐状态为active时才设置开始时间和结束时间
		if req.Status == "active" {
			if req.StartTime == "" {
				return gerror.New("状态为active时，开始时间必填")
			}

			// 解析开始时间
			startTime, err := gtime.StrToTime(req.StartTime)
			if err != nil {
				return gerror.New("开始时间格式错误")
			}
			userPackage.StartTime = startTime

			// 计算结束时间
			if packageInfo.DurationMinutes > 0 {
				endTimeStr := startTime.Add(time.Minute * time.Duration(packageInfo.DurationMinutes)).Format("Y-m-d H:i:s")
				endTime := gtime.NewFromStr(endTimeStr)
				userPackage.EndTime = endTime
			}
		}

		// 插入记录
		result, err := tx.Model("user_packages").Data(userPackage).Insert()
		if err != nil {
			return err
		}

		// 获取插入的ID
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// 查询创建的记录
		var createdPackage entity.UserPackages
		err = tx.Model("user_packages").Where("id", lastInsertID).Scan(&createdPackage)
		if err != nil {
			return err
		}

		// 转换为API响应格式
		detail = &v1.AdminUserPackage{}
		if err = gconv.Struct(createdPackage, detail); err != nil {
			return err
		}

		// 格式化时间
		if createdPackage.StartTime != nil {
			detail.StartTime = createdPackage.StartTime.Format("Y-m-d H:i:s")
		}
		if createdPackage.EndTime != nil {
			detail.EndTime = createdPackage.EndTime.Format("Y-m-d H:i:s")
		}
		if createdPackage.CreatedAt != nil {
			detail.CreatedAt = createdPackage.CreatedAt.Format("Y-m-d H:i:s")
		}
		if createdPackage.UpdatedAt != nil {
			detail.UpdatedAt = createdPackage.UpdatedAt.Format("Y-m-d H:i:s")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return detail, nil
}

// UpdateUserPackageStatus 更新用户套餐状态
func (s *userPackageService) UpdateUserPackageStatus(ctx context.Context, req *v1.AdminUserPackageUpdateStatusReq) (detail *v1.AdminUserPackage, err error) {
	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询用户套餐
		var userPackage entity.UserPackages
		err = tx.Model("user_packages").Where("id", req.UserPackageID).Scan(&userPackage)
		if err != nil {
			return err
		}
		if userPackage.Id == 0 {
			return gerror.New("用户套餐不存在")
		}

		// 更新状态
		_, err = tx.Model("user_packages").
			Data(g.Map{
				"status":     req.Status,
				"updated_at": gtime.Now(),
			}).
			Where("id", req.UserPackageID).
			Update()
		if err != nil {
			return err
		}

		// 查询更新后的记录
		err = tx.Model("user_packages").Where("id", req.UserPackageID).Scan(&userPackage)
		if err != nil {
			return err
		}

		// 转换为API响应格式
		detail = &v1.AdminUserPackage{}
		if err = gconv.Struct(userPackage, detail); err != nil {
			return err
		}

		// 格式化时间
		if userPackage.StartTime != nil {
			detail.StartTime = userPackage.StartTime.Format("Y-m-d H:i:s")
		}
		if userPackage.EndTime != nil {
			detail.EndTime = userPackage.EndTime.Format("Y-m-d H:i:s")
		}
		if userPackage.CreatedAt != nil {
			detail.CreatedAt = userPackage.CreatedAt.Format("Y-m-d H:i:s")
		}
		if userPackage.UpdatedAt != nil {
			detail.UpdatedAt = userPackage.UpdatedAt.Format("Y-m-d H:i:s")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return detail, nil
}

// GetUserActivePackages 获取用户有效套餐
func (s *userPackageService) GetUserActivePackages(ctx context.Context, userID int64) (list []v1.AdminUserPackage, err error) {
	// 查询用户激活状态的套餐
	var userPackages []entity.UserPackages
	err = g.DB().Model("user_packages").
		Where("user_id", userID).
		Where("status", "active").
		Order("id DESC").
		Scan(&userPackages)
	if err != nil {
		return nil, err
	}

	// 实时检查并更新套餐状态
	for i := range userPackages {
		if err = s.CheckAndUpdatePackageStatus(ctx, &userPackages[i]); err != nil {
			g.Log().Errorf(ctx, "检查更新套餐状态失败: %v", err)
		}
	}

	// 过滤出仍然是激活状态的套餐
	activePackages := make([]entity.UserPackages, 0, len(userPackages))
	for _, pkg := range userPackages {
		if pkg.Status == "active" {
			activePackages = append(activePackages, pkg)
		}
	}

	// 转换为API响应格式
	list = make([]v1.AdminUserPackage, 0, len(activePackages))
	for _, item := range activePackages {
		var adminUserPackage v1.AdminUserPackage
		if err = gconv.Struct(item, &adminUserPackage); err != nil {
			return nil, err
		}

		// 格式化时间
		if item.StartTime != nil {
			adminUserPackage.StartTime = item.StartTime.Format("Y-m-d H:i:s")
		}
		if item.EndTime != nil {
			adminUserPackage.EndTime = item.EndTime.Format("Y-m-d H:i:s")
		}
		if item.CreatedAt != nil {
			adminUserPackage.CreatedAt = item.CreatedAt.Format("Y-m-d H:i:s")
		}
		if item.UpdatedAt != nil {
			adminUserPackage.UpdatedAt = item.UpdatedAt.Format("Y-m-d H:i:s")
		}

		list = append(list, adminUserPackage)
	}

	return list, nil
}

// ActivateUserPackage 激活用户套餐（首次使用时）
func (s *userPackageService) ActivateUserPackage(ctx context.Context, userPackageID int64) error {
	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询用户套餐
		var userPackage entity.UserPackages
		err := tx.Model("user_packages").Where("id", userPackageID).Scan(&userPackage)
		if err != nil {
			return err
		}
		if userPackage.Id == 0 {
			return gerror.New("用户套餐不存在")
		}

		// 检查套餐状态
		if userPackage.Status != "pending" {
			return gerror.New("只有待激活的套餐可以被激活")
		}

		// 检查是否已激活 - 通过检查start_time是否为空或特殊值来判断
		specialDate := gtime.NewFromStr("0001-01-01 00:00:00")
		if userPackage.StartTime != nil && !userPackage.StartTime.Equal(specialDate) {
			return gerror.New("套餐已经被激活")
		}

		// 查询套餐信息，获取套餐时长
		var packageInfo entity.DrinkAllYouCanPackages
		err = tx.Model("drink_all_you_can_packages").Where("id", userPackage.PackageId).Scan(&packageInfo)
		if err != nil {
			return err
		}

		// 设置开始时间为当前时间
		now := gtime.Now()

		// 计算结束时间
		var endTime *gtime.Time
		if packageInfo.DurationMinutes > 0 {
			endTimeStr := now.Add(time.Minute * time.Duration(packageInfo.DurationMinutes)).Format("Y-m-d H:i:s")
			endTime = gtime.NewFromStr(endTimeStr)
		}

		// 更新套餐状态
		_, err = tx.Model("user_packages").
			Data(g.Map{
				"start_time": now,
				"end_time":   endTime,
				"status":     "active",
				"updated_at": now,
			}).
			Where("id", userPackageID).
			Update()

		return err
	})
}

// GetUserPendingPackages 获取用户待激活的套餐列表
// 注意：在新逻辑下，套餐在购买时就会被激活，此方法主要用于处理旧数据
func (s *userPackageService) GetUserPendingPackages(ctx context.Context, userID int64) (list []v1.AdminUserPackage, err error) {
	// 查询用户待激活套餐
	var userPackages []entity.UserPackages
	err = g.DB().Model("user_packages").
		Where("user_id", userID).
		Where("status", "pending").
		Order("id DESC").
		Scan(&userPackages)
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	list = make([]v1.AdminUserPackage, 0, len(userPackages))
	for _, item := range userPackages {
		var adminUserPackage v1.AdminUserPackage
		if err = gconv.Struct(item, &adminUserPackage); err != nil {
			return nil, err
		}

		// 格式化时间
		if item.StartTime != nil {
			adminUserPackage.StartTime = item.StartTime.Format("Y-m-d H:i:s")
		}
		if item.EndTime != nil {
			adminUserPackage.EndTime = item.EndTime.Format("Y-m-d H:i:s")
		}
		if item.CreatedAt != nil {
			adminUserPackage.CreatedAt = item.CreatedAt.Format("Y-m-d H:i:s")
		}
		if item.UpdatedAt != nil {
			adminUserPackage.UpdatedAt = item.UpdatedAt.Format("Y-m-d H:i:s")
		}

		list = append(list, adminUserPackage)
	}

	return list, nil
}

// CreateUserPackageAfterPurchase 用户购买套餐后创建待支付的套餐记录
func (s *userPackageService) CreateUserPackageAfterPurchase(ctx context.Context, userID, packageID, orderID int64) (userPackageID int64, err error) {
	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询套餐信息
		var packageInfo entity.DrinkAllYouCanPackages
		err = tx.Model("drink_all_you_can_packages").Where("id", packageID).Scan(&packageInfo)
		if err != nil {
			return err
		}
		if packageInfo.Id == 0 {
			return gerror.New("套餐不存在")
		}

		// 验证用户是否存在
		var userCount int
		userCount, err = tx.Model("users").Where("id", userID).Count()
		if err != nil {
			return err
		}
		if userCount == 0 {
			return gerror.New("用户不存在")
		}

		// 验证订单是否存在
		var orderCount int
		orderCount, err = tx.Model("orders").Where("id", orderID).Count()
		if err != nil {
			return err
		}
		if orderCount == 0 {
			return gerror.New("订单不存在")
		}

		// 获取当前时间
		now := gtime.Now()

		// 创建初始状态为 pending 的套餐记录，不设置开始时间和结束时间
		// 待支付成功后再设置开始时间和结束时间
		_, err = tx.Exec(
			"INSERT INTO user_packages(user_id, package_id, order_id, start_time, end_time, status, created_at, updated_at) "+
				"VALUES(?, ?, ?, '0001-01-01 00:00:00', '0001-01-01 00:00:00', 'pending', ?, ?)",
			userID, packageID, orderID, now, now,
		)
		if err != nil {
			return err
		}

		// 获取插入的ID
		var result gdb.Value
		result, err = tx.GetValue("SELECT LAST_INSERT_ID()")
		if err != nil {
			return err
		}
		userPackageID = result.Int64()

		return nil
	})

	return userPackageID, err
}

// ActivateUserPackageAfterPayment 支付成功后激活用户套餐
func (s *userPackageService) ActivateUserPackageAfterPayment(ctx context.Context, orderID int64) error {
	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 查询与订单关联的用户套餐
		var userPackage entity.UserPackages
		err := tx.Model("user_packages").Where("order_id", orderID).Scan(&userPackage)
		if err != nil {
			return err
		}
		if userPackage.Id == 0 {
			return gerror.New("未找到与订单关联的用户套餐")
		}

		// 2. 检查套餐状态
		if userPackage.Status != "pending" {
			// 如果套餐不是待支付状态，则不需要处理
			return nil
		}

		// 3. 查询套餐信息，获取套餐时长
		var packageInfo entity.DrinkAllYouCanPackages
		err = tx.Model("drink_all_you_can_packages").Where("id", userPackage.PackageId).Scan(&packageInfo)
		if err != nil {
			return err
		}

		// 4. 设置开始时间为当前时间
		now := gtime.Now()

		// 5. 计算结束时间
		var endTime *gtime.Time
		if packageInfo.DurationMinutes > 0 {
			endTimeStr := now.Add(time.Minute * time.Duration(packageInfo.DurationMinutes)).Format("Y-m-d H:i:s")
			endTime = gtime.NewFromStr(endTimeStr)
		} else {
			// 如果没有设置时长，则默认30天
			endTimeStr := now.AddDate(0, 1, 0).Format("Y-m-d H:i:s")
			endTime = gtime.NewFromStr(endTimeStr)
		}

		// 6. 更新套餐状态为激活
		_, err = tx.Model("user_packages").
			Data(g.Map{
				"start_time": now,
				"end_time":   endTime,
				"status":     "active",
				"updated_at": now,
			}).
			Where("id", userPackage.Id).
			Update()

		return err
	})
}
