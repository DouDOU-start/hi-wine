package service

import (
	v1 "backend/api/admin/v1"
	userv1 "backend/api/user/v1"
	"backend/internal/consts"
	"backend/internal/model/entity"
	"backend/internal/utility"
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

	// GetUserMyPackages 获取用户个人套餐列表
	GetUserMyPackages(ctx context.Context, userID int64, status string) (list []userv1.UserMyPackage, err error)
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
	if userPackage.Status == consts.PackageStatusExpired {
		return nil
	}

	// 如果套餐状态是激活状态，且结束时间已过期，则更新为过期状态
	if userPackage.Status == consts.PackageStatusActive && userPackage.EndTime != nil && userPackage.EndTime.Before(gtime.Now()) {
		_, err := g.DB().Model("user_packages").
			Data(g.Map{
				"status":     consts.PackageStatusExpired,
				"updated_at": gtime.Now(),
			}).
			Where("id", userPackage.Id).
			Update()

		if err != nil {
			g.Log().Errorf(ctx, "更新套餐状态失败，套餐ID:%d, 错误:%v", userPackage.Id, err)
			return err
		}

		// 更新当前对象的状态
		userPackage.Status = consts.PackageStatusExpired
		userPackage.UpdatedAt = gtime.Now()
		g.Log().Debugf(ctx, "已将套餐 ID:%d 状态更新为已过期", userPackage.Id)
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
		// 如果查询激活状态的套餐，直接在SQL中过滤过期的
		if req.Status == consts.PackageStatusActive {
			m = m.WhereGT("up.end_time", gtime.Now())
		}
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
		g.Log().Errorf(ctx, "查询用户套餐总数失败: %v", err)
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
		g.Log().Errorf(ctx, "查询用户套餐列表失败: %v", err)
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

		// 使用工具函数格式化时间
		utility.FormatUserPackageTimes(&item, &adminUserPackage)

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
		g.Log().Errorf(ctx, "查询用户套餐失败, ID:%d, 错误:%v", id, err)
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

	// 使用工具函数格式化时间
	utility.FormatUserPackageTimes(&userPackage, detail)

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
				endTimeStr := startTime.Add(time.Minute * time.Duration(packageInfo.DurationMinutes)).Format(consts.TimeFormatStandard)
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
	// 查询用户激活状态的套餐，并且结束时间大于当前时间
	var userPackages []entity.UserPackages
	err = g.DB().Model("user_packages").
		Where("user_id", userID).
		Where("status", consts.PackageStatusActive).
		WhereGT("end_time", gtime.Now()).
		Order("id DESC").
		Scan(&userPackages)
	if err != nil {
		g.Log().Errorf(ctx, "查询用户有效套餐失败, 用户ID:%d, 错误:%v", userID, err)
		return nil, err
	}

	// 转换为API响应格式
	list = make([]v1.AdminUserPackage, 0, len(userPackages))
	for _, item := range userPackages {
		var adminUserPackage v1.AdminUserPackage
		if err = gconv.Struct(item, &adminUserPackage); err != nil {
			return nil, err
		}

		// 使用工具函数格式化时间
		utility.FormatUserPackageTimes(&item, &adminUserPackage)

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
			g.Log().Errorf(ctx, "查询用户套餐失败, ID:%d, 错误:%v", userPackageID, err)
			return err
		}
		if userPackage.Id == 0 {
			return gerror.New("用户套餐不存在")
		}

		// 检查套餐状态
		if userPackage.Status != consts.PackageStatusPending {
			return gerror.New("只有待激活的套餐可以被激活")
		}

		// 查询套餐详情，获取有效时长
		var packageInfo entity.DrinkAllYouCanPackages
		err = tx.Model("drink_all_you_can_packages").Where("id", userPackage.PackageId).Scan(&packageInfo)
		if err != nil {
			g.Log().Errorf(ctx, "查询套餐详情失败, ID:%d, 错误:%v", userPackage.PackageId, err)
			return err
		}
		if packageInfo.Id == 0 {
			return gerror.New("套餐不存在")
		}

		// 计算开始时间和结束时间
		startTime := gtime.Now()
		var endTime *gtime.Time
		if packageInfo.DurationMinutes > 0 {
			// 使用字符串格式化方式创建结束时间
			endTimeStr := startTime.Add(time.Minute * time.Duration(packageInfo.DurationMinutes)).Format(consts.TimeFormatStandard)
			endTime = gtime.NewFromStr(endTimeStr)
		}

		// 更新套餐状态为激活
		_, err = tx.Model("user_packages").
			Data(g.Map{
				"status":     consts.PackageStatusActive,
				"start_time": startTime,
				"end_time":   endTime,
				"updated_at": gtime.Now(),
			}).
			Where("id", userPackageID).
			Update()

		if err != nil {
			g.Log().Errorf(ctx, "激活套餐失败, ID:%d, 错误:%v", userPackageID, err)
			return err
		}

		g.Log().Infof(ctx, "已成功激活用户套餐, ID:%d", userPackageID)
		return nil
	})
}

// GetUserPendingPackages 获取用户待激活的套餐列表
func (s *userPackageService) GetUserPendingPackages(ctx context.Context, userID int64) (list []v1.AdminUserPackage, err error) {
	// 查询用户待激活套餐
	var userPackages []entity.UserPackages
	err = g.DB().Model("user_packages").
		Where("user_id", userID).
		Where("status", consts.PackageStatusPending).
		Order("id DESC").
		Scan(&userPackages)
	if err != nil {
		g.Log().Errorf(ctx, "查询用户待激活套餐失败, 用户ID:%d, 错误:%v", userID, err)
		return nil, err
	}

	// 转换为API响应格式
	list = make([]v1.AdminUserPackage, 0, len(userPackages))
	for _, item := range userPackages {
		var adminUserPackage v1.AdminUserPackage
		if err = gconv.Struct(item, &adminUserPackage); err != nil {
			return nil, err
		}

		// 使用工具函数格式化时间
		utility.FormatUserPackageTimes(&item, &adminUserPackage)

		list = append(list, adminUserPackage)
	}

	return list, nil
}

// CreateUserPackageAfterPurchase 用户购买套餐后创建待支付的套餐记录
func (s *userPackageService) CreateUserPackageAfterPurchase(ctx context.Context, userID, packageID, orderID int64) (userPackageID int64, err error) {
	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询套餐是否存在
		var packageInfo entity.DrinkAllYouCanPackages
		err = tx.Model("drink_all_you_can_packages").Where("id", packageID).Scan(&packageInfo)
		if err != nil {
			g.Log().Errorf(ctx, "查询套餐失败, ID:%d, 错误:%v", packageID, err)
			return err
		}
		if packageInfo.Id == 0 {
			return gerror.New("套餐不存在")
		}

		// 验证套餐是否可购买
		if packageInfo.IsActive != 1 {
			return gerror.New("该套餐当前不可购买")
		}

		// 创建用户套餐记录
		now := gtime.Now()
		result, err := tx.Model("user_packages").Insert(g.Map{
			"user_id":    userID,
			"package_id": packageID,
			"order_id":   orderID,
			"status":     consts.PackageStatusPending, // 待支付状态
			"created_at": now,
			"updated_at": now,
		})
		if err != nil {
			g.Log().Errorf(ctx, "创建用户套餐记录失败, 用户ID:%d, 套餐ID:%d, 错误:%v", userID, packageID, err)
			return err
		}

		// 获取新创建的记录ID
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			g.Log().Errorf(ctx, "获取新创建套餐ID失败: %v", err)
			return err
		}

		userPackageID = lastInsertID
		g.Log().Infof(ctx, "成功创建用户套餐记录, ID:%d", userPackageID)
		return nil
	})

	return userPackageID, err
}

// ActivateUserPackageAfterPayment 支付成功后激活用户套餐
func (s *userPackageService) ActivateUserPackageAfterPayment(ctx context.Context, orderID int64) error {
	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询订单关联的套餐
		var userPackages []entity.UserPackages
		err := tx.Model("user_packages").Where("order_id", orderID).Scan(&userPackages)
		if err != nil {
			g.Log().Errorf(ctx, "查询订单关联套餐失败, 订单ID:%d, 错误:%v", orderID, err)
			return err
		}
		if len(userPackages) == 0 {
			g.Log().Infof(ctx, "订单未关联套餐, 订单ID:%d", orderID)
			return nil
		}

		// 遍历激活所有套餐
		for _, userPackage := range userPackages {
			// 如果套餐已经是激活状态，则跳过
			if userPackage.Status == consts.PackageStatusActive {
				continue
			}

			// 查询套餐详情，获取有效时长
			var packageInfo entity.DrinkAllYouCanPackages
			err = tx.Model("drink_all_you_can_packages").Where("id", userPackage.PackageId).Scan(&packageInfo)
			if err != nil {
				g.Log().Errorf(ctx, "查询套餐详情失败, ID:%d, 错误:%v", userPackage.PackageId, err)
				return err
			}
			if packageInfo.Id == 0 {
				g.Log().Errorf(ctx, "套餐不存在, ID:%d", userPackage.PackageId)
				return gerror.Newf("套餐ID:%d 不存在", userPackage.PackageId)
			}

			// 计算开始时间和结束时间
			now := gtime.Now()
			var endTime *gtime.Time
			if packageInfo.DurationMinutes > 0 {
				endTimeStr := now.Add(time.Minute * time.Duration(packageInfo.DurationMinutes)).Format(consts.TimeFormatStandard)
				endTime = gtime.NewFromStr(endTimeStr)
			} else {
				// 如果没有设置时长，则默认30天
				endTimeStr := now.AddDate(0, 1, 0).Format(consts.TimeFormatStandard)
				endTime = gtime.NewFromStr(endTimeStr)
			}

			// 更新套餐状态为激活
			_, err = tx.Model("user_packages").
				Data(g.Map{
					"status":     consts.PackageStatusActive,
					"start_time": now,
					"end_time":   endTime,
					"updated_at": now,
				}).
				Where("id", userPackage.Id).
				Update()

			if err != nil {
				g.Log().Errorf(ctx, "激活套餐失败, ID:%d, 错误:%v", userPackage.Id, err)
				return err
			}

			g.Log().Infof(ctx, "已成功激活用户套餐, ID:%d", userPackage.Id)
		}

		return nil
	})
}

// GetUserMyPackages 获取用户个人套餐列表
func (s *userPackageService) GetUserMyPackages(ctx context.Context, userID int64, status string) (list []userv1.UserMyPackage, err error) {
	// 构建查询条件
	m := g.DB().Model("user_packages up")

	// 关联套餐表和订单表
	m = m.LeftJoin("drink_all_you_can_packages p", "p.id = up.package_id")
	m = m.LeftJoin("orders o", "o.id = up.order_id")

	// 添加筛选条件
	m = m.Where("up.user_id", userID)
	if status != "" {
		m = m.Where("up.status", status)
		// 如果查询激活状态的套餐，直接在SQL中过滤过期的
		if status == consts.PackageStatusActive {
			m = m.WhereGT("up.end_time", gtime.Now())
		}
	}

	// 查询数据，直接在SQL中计算剩余时间
	type JoinedResult struct {
		entity.UserPackages         // 用户套餐信息
		PackageName         string  `json:"package_name"` // 套餐名称
		Price               float64 `json:"price"`        // 套餐价格
		OrderSN             string  `json:"order_sn"`     // 订单号
	}

	var joinedResults []JoinedResult
	err = m.Fields("up.*, p.name as package_name, p.price as price, o.order_sn as order_sn").
		Order("up.id DESC").
		Scan(&joinedResults)
	if err != nil {
		g.Log().Errorf(ctx, "获取用户套餐列表失败: %v", err)
		return nil, err
	}

	// 实时检查并更新套餐状态
	for i := range joinedResults {
		if err = s.CheckAndUpdatePackageStatus(ctx, &joinedResults[i].UserPackages); err != nil {
			g.Log().Errorf(ctx, "检查更新套餐状态失败: %v", err)
		}
	}

	// 转换为API响应格式
	list = make([]userv1.UserMyPackage, 0, len(joinedResults))
	for _, item := range joinedResults {
		var userMyPackage userv1.UserMyPackage

		// 使用工具函数格式化时间和计算剩余时间
		utility.FormatUserMyPackage(&item.UserPackages, &userMyPackage)

		// 设置其他字段
		userMyPackage.PackageName = item.PackageName
		userMyPackage.Price = item.Price
		userMyPackage.OrderSN = item.OrderSN

		list = append(list, userMyPackage)
	}

	return list, nil
}
