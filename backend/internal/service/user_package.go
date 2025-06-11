package service

import (
	v1 "backend/api/admin/v1"
	userv1 "backend/api/user/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/utility"
	"context"
	"database/sql"
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

	// GetUserPackageFullDetail 获取用户套餐详细信息（包含用户信息、套餐信息和使用情况）
	GetUserPackageFullDetail(ctx context.Context, userPackageID int64) (detail *v1.AdminUserPackageFullDetail, err error)
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
	m = m.LeftJoin("orders o", "o.id = up.order_id")

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

	// 自定义查询结果结构
	type JoinedUserPackage struct {
		entity.UserPackages
		UserName        string  `json:"user_name"`
		Nickname        string  `json:"nickname"`
		Phone           string  `json:"phone"`
		PackageName     string  `json:"package_name"`
		PackagePrice    float64 `json:"package_price"`
		OrderSN         string  `json:"order_sn"`
		DurationMinutes int     `json:"duration_minutes"`
	}

	// 查询数据
	var joinedUserPackages []JoinedUserPackage
	err = m.Fields("up.*, u.nickname, u.phone, p.name as package_name, p.price as package_price, o.order_sn, p.duration_minutes").
		Order("up.id DESC").
		Limit((page-1)*limit, limit).
		Scan(&joinedUserPackages)
	if err != nil {
		g.Log().Errorf(ctx, "查询用户套餐列表失败: %v", err)
		return nil, 0, err
	}

	// 转换为API响应格式
	list = make([]v1.AdminUserPackage, 0, len(joinedUserPackages))
	for _, item := range joinedUserPackages {
		// 实时检查并更新套餐状态
		err = s.CheckAndUpdatePackageStatus(ctx, &item.UserPackages)
		if err != nil {
			g.Log().Errorf(ctx, "检查更新套餐状态失败: %v", err)
		}

		var adminUserPackage v1.AdminUserPackage

		// 基本字段转换
		adminUserPackage.ID = int64(item.Id)
		adminUserPackage.UserID = int64(item.UserId)
		adminUserPackage.PackageID = int64(item.PackageId)
		adminUserPackage.OrderID = int64(item.OrderId)
		adminUserPackage.Status = item.Status

		// 添加用户相关信息
		adminUserPackage.UserName = item.Nickname
		if adminUserPackage.UserName == "" {
			adminUserPackage.UserName = "用户" + gconv.String(item.UserId)
		}
		adminUserPackage.UserPhone = item.Phone

		// 添加套餐相关信息
		adminUserPackage.PackageName = item.PackageName
		adminUserPackage.PackagePrice = item.PackagePrice
		adminUserPackage.OrderSN = item.OrderSN

		// 状态描述
		switch item.Status {
		case consts.PackageStatusActive:
			adminUserPackage.StatusDesc = "已激活"
		case consts.PackageStatusExpired:
			adminUserPackage.StatusDesc = "已过期"
		case consts.PackageStatusPending:
			adminUserPackage.StatusDesc = "待激活"
		case "refunded":
			adminUserPackage.StatusDesc = "已退款"
		default:
			adminUserPackage.StatusDesc = item.Status
		}

		// 时间处理
		if item.StartTime != nil {
			adminUserPackage.StartTime = utility.FormatTimeOrEmpty(item.StartTime)
		}
		if item.EndTime != nil {
			adminUserPackage.EndTime = utility.FormatTimeOrEmpty(item.EndTime)

			// 计算有效期描述
			if item.StartTime != nil && item.Status == consts.PackageStatusActive {
				duration := item.EndTime.Sub(gtime.Now())
				days := int(duration.Hours()) / 24
				hours := int(duration.Hours()) % 24
				if days > 0 {
					adminUserPackage.ValidPeriod = gconv.String(days) + "天"
					if hours > 0 {
						adminUserPackage.ValidPeriod += gconv.String(hours) + "小时"
					}
				} else if hours > 0 {
					adminUserPackage.ValidPeriod = gconv.String(hours) + "小时"
				} else {
					adminUserPackage.ValidPeriod = "即将到期"
				}
			} else if item.Status == consts.PackageStatusExpired {
				adminUserPackage.ValidPeriod = "已过期"
			} else if item.Status == consts.PackageStatusPending {
				adminUserPackage.ValidPeriod = gconv.String(item.DurationMinutes/60) + "小时"
			} else {
				adminUserPackage.ValidPeriod = "-"
			}
		} else {
			adminUserPackage.ValidPeriod = "-"
		}

		adminUserPackage.CreatedAt = utility.FormatTimeOrEmpty(item.CreatedAt)
		adminUserPackage.UpdatedAt = utility.FormatTimeOrEmpty(item.UpdatedAt)

		list = append(list, adminUserPackage)
	}

	return list, count, nil
}

// GetUserPackageDetail 获取用户套餐详情
func (s *userPackageService) GetUserPackageDetail(ctx context.Context, id int64) (detail *v1.AdminUserPackage, err error) {
	// 查询用户套餐（关联用户和套餐信息）
	type JoinedUserPackage struct {
		entity.UserPackages
		Nickname        string  `json:"nickname"`
		Phone           string  `json:"phone"`
		PackageName     string  `json:"package_name"`
		PackagePrice    float64 `json:"package_price"`
		OrderSN         string  `json:"order_sn"`
		DurationMinutes int     `json:"duration_minutes"`
	}

	var joinedUserPackage JoinedUserPackage
	err = g.DB().Model("user_packages").As("up").
		LeftJoin("users u", "u.id = up.user_id").
		LeftJoin("drink_all_you_can_packages p", "p.id = up.package_id").
		LeftJoin("orders o", "o.id = up.order_id").
		Fields("up.*, u.nickname, u.phone, p.name as package_name, p.price as package_price, o.order_sn, p.duration_minutes").
		Where("up.id", id).
		Scan(&joinedUserPackage)
	if err != nil {
		g.Log().Errorf(ctx, "查询用户套餐失败, ID:%d, 错误:%v", id, err)
		return nil, err
	}
	if joinedUserPackage.Id == 0 {
		return nil, gerror.New("用户套餐不存在")
	}

	// 实时检查并更新套餐状态
	if err = s.CheckAndUpdatePackageStatus(ctx, &joinedUserPackage.UserPackages); err != nil {
		g.Log().Errorf(ctx, "检查更新套餐状态失败: %v", err)
	}

	// 转换为API响应格式
	detail = &v1.AdminUserPackage{
		ID:           int64(joinedUserPackage.Id),
		UserID:       int64(joinedUserPackage.UserId),
		PackageID:    int64(joinedUserPackage.PackageId),
		OrderID:      int64(joinedUserPackage.OrderId),
		Status:       joinedUserPackage.Status,
		UserName:     joinedUserPackage.Nickname,
		UserPhone:    joinedUserPackage.Phone,
		PackageName:  joinedUserPackage.PackageName,
		PackagePrice: joinedUserPackage.PackagePrice,
		OrderSN:      joinedUserPackage.OrderSN,
	}

	// 如果用户名为空，使用默认名称
	if detail.UserName == "" {
		detail.UserName = "用户" + gconv.String(joinedUserPackage.UserId)
	}

	// 状态描述
	switch joinedUserPackage.Status {
	case consts.PackageStatusActive:
		detail.StatusDesc = "已激活"
	case consts.PackageStatusExpired:
		detail.StatusDesc = "已过期"
	case consts.PackageStatusPending:
		detail.StatusDesc = "待激活"
	case "refunded":
		detail.StatusDesc = "已退款"
	default:
		detail.StatusDesc = joinedUserPackage.Status
	}

	// 时间处理
	if joinedUserPackage.StartTime != nil {
		detail.StartTime = utility.FormatTimeOrEmpty(joinedUserPackage.StartTime)
	}
	if joinedUserPackage.EndTime != nil {
		detail.EndTime = utility.FormatTimeOrEmpty(joinedUserPackage.EndTime)

		// 计算有效期描述
		if joinedUserPackage.StartTime != nil && joinedUserPackage.Status == consts.PackageStatusActive {
			duration := joinedUserPackage.EndTime.Sub(gtime.Now())
			days := int(duration.Hours()) / 24
			hours := int(duration.Hours()) % 24
			if days > 0 {
				detail.ValidPeriod = gconv.String(days) + "天"
				if hours > 0 {
					detail.ValidPeriod += gconv.String(hours) + "小时"
				}
			} else if hours > 0 {
				detail.ValidPeriod = gconv.String(hours) + "小时"
			} else {
				detail.ValidPeriod = "即将到期"
			}
		} else if joinedUserPackage.Status == consts.PackageStatusExpired {
			detail.ValidPeriod = "已过期"
		} else if joinedUserPackage.Status == consts.PackageStatusPending {
			detail.ValidPeriod = gconv.String(joinedUserPackage.DurationMinutes/60) + "小时"
		} else {
			detail.ValidPeriod = "-"
		}
	} else {
		detail.ValidPeriod = "-"
	}

	detail.CreatedAt = utility.FormatTimeOrEmpty(joinedUserPackage.CreatedAt)
	detail.UpdatedAt = utility.FormatTimeOrEmpty(joinedUserPackage.UpdatedAt)

	return detail, nil
}

// CreateUserPackage 创建用户套餐
func (s *userPackageService) CreateUserPackage(ctx context.Context, req *v1.AdminUserPackageCreateReq) (detail *v1.AdminUserPackage, err error) {
	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 检查用户是否存在
		var userCount int
		userCount, err = tx.Model("users").Where("id", req.UserID).Count()
		if err != nil {
			return err
		}
		if userCount == 0 {
			return gerror.New("用户不存在")
		}

		// 检查套餐是否存在
		var packageInfo entity.DrinkAllYouCanPackages
		err = tx.Model("drink_all_you_can_packages").Where("id", req.PackageID).Scan(&packageInfo)
		if err != nil {
			return err
		}
		if packageInfo.Id == 0 {
			return gerror.New("套餐不存在")
		}

		// 创建套餐记录
		now := gtime.Now()
		var startTime *gtime.Time
		var endTime *gtime.Time

		// 如果状态为激活，则设置开始时间和结束时间
		if req.Status == consts.PackageStatusActive {
			if req.StartTime != "" {
				// 解析指定的开始时间
				startTimeObj, err := gtime.StrToTime(req.StartTime)
				if err != nil {
					return gerror.New("开始时间格式错误")
				}
				startTime = startTimeObj
			} else {
				// 默认为当前时间
				startTime = now
			}

			// 计算结束时间 = 开始时间 + 套餐时长
			endTime = gtime.NewFromTime(startTime.Time.Add(time.Minute * time.Duration(packageInfo.DurationMinutes)))
		}

		// 创建用户套餐记录
		userPackageData := entity.UserPackages{
			UserId:    int(req.UserID),
			PackageId: int(req.PackageID),
			OrderId:   int(req.OrderID),
			StartTime: startTime,
			EndTime:   endTime,
			Status:    req.Status,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// 插入记录
		result, err := tx.Model("user_packages").Data(userPackageData).Insert()
		if err != nil {
			return err
		}

		// 获取插入的ID
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// 查询创建的套餐详细信息（包含关联数据）
		type JoinedUserPackage struct {
			entity.UserPackages
			Nickname        string  `json:"nickname"`
			Phone           string  `json:"phone"`
			PackageName     string  `json:"package_name"`
			PackagePrice    float64 `json:"package_price"`
			OrderSN         string  `json:"order_sn"`
			DurationMinutes int     `json:"duration_minutes"`
		}

		var joinedUserPackage JoinedUserPackage
		err = tx.Model("user_packages").As("up").
			LeftJoin("users u", "u.id = up.user_id").
			LeftJoin("drink_all_you_can_packages p", "p.id = up.package_id").
			LeftJoin("orders o", "o.id = up.order_id").
			Fields("up.*, u.nickname, u.phone, p.name as package_name, p.price as package_price, o.order_sn, p.duration_minutes").
			Where("up.id", lastInsertID).
			Scan(&joinedUserPackage)
		if err != nil {
			return err
		}

		// 转换为API响应格式
		detail = &v1.AdminUserPackage{
			ID:           int64(joinedUserPackage.Id),
			UserID:       int64(joinedUserPackage.UserId),
			PackageID:    int64(joinedUserPackage.PackageId),
			OrderID:      int64(joinedUserPackage.OrderId),
			Status:       joinedUserPackage.Status,
			UserName:     joinedUserPackage.Nickname,
			UserPhone:    joinedUserPackage.Phone,
			PackageName:  joinedUserPackage.PackageName,
			PackagePrice: joinedUserPackage.PackagePrice,
			OrderSN:      joinedUserPackage.OrderSN,
		}

		// 如果用户名为空，使用默认名称
		if detail.UserName == "" {
			detail.UserName = "用户" + gconv.String(joinedUserPackage.UserId)
		}

		// 状态描述
		switch joinedUserPackage.Status {
		case consts.PackageStatusActive:
			detail.StatusDesc = "已激活"
		case consts.PackageStatusExpired:
			detail.StatusDesc = "已过期"
		case consts.PackageStatusPending:
			detail.StatusDesc = "待激活"
		case "refunded":
			detail.StatusDesc = "已退款"
		default:
			detail.StatusDesc = joinedUserPackage.Status
		}

		// 时间处理
		if joinedUserPackage.StartTime != nil {
			detail.StartTime = utility.FormatTimeOrEmpty(joinedUserPackage.StartTime)
		}
		if joinedUserPackage.EndTime != nil {
			detail.EndTime = utility.FormatTimeOrEmpty(joinedUserPackage.EndTime)

			// 计算有效期描述
			if joinedUserPackage.StartTime != nil && joinedUserPackage.Status == consts.PackageStatusActive {
				duration := joinedUserPackage.EndTime.Sub(gtime.Now())
				days := int(duration.Hours()) / 24
				hours := int(duration.Hours()) % 24
				if days > 0 {
					detail.ValidPeriod = gconv.String(days) + "天"
					if hours > 0 {
						detail.ValidPeriod += gconv.String(hours) + "小时"
					}
				} else if hours > 0 {
					detail.ValidPeriod = gconv.String(hours) + "小时"
				} else {
					detail.ValidPeriod = "即将到期"
				}
			} else if joinedUserPackage.Status == consts.PackageStatusExpired {
				detail.ValidPeriod = "已过期"
			} else if joinedUserPackage.Status == consts.PackageStatusPending {
				detail.ValidPeriod = gconv.String(joinedUserPackage.DurationMinutes/60) + "小时"
			} else {
				detail.ValidPeriod = "-"
			}
		} else {
			detail.ValidPeriod = "-"
		}

		detail.CreatedAt = utility.FormatTimeOrEmpty(joinedUserPackage.CreatedAt)
		detail.UpdatedAt = utility.FormatTimeOrEmpty(joinedUserPackage.UpdatedAt)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return detail, nil
}

// UpdateUserPackageStatus 更新用户套餐状态
func (s *userPackageService) UpdateUserPackageStatus(ctx context.Context, req *v1.AdminUserPackageUpdateStatusReq) (*v1.AdminUserPackage, error) {
	// 1. 开启事务
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 2. 查询并锁定用户套餐记录
		var userPackage entity.UserPackages
		sqlTx := tx.GetSqlTX()
		err := sqlTx.QueryRowContext(ctx, "SELECT * FROM `user_packages` WHERE `id` = ? FOR UPDATE", req.UserPackageID).Scan(
			&userPackage.Id, &userPackage.UserId, &userPackage.PackageId, &userPackage.OrderId, &userPackage.StartTime, &userPackage.EndTime, &userPackage.Status, &userPackage.CreatedAt, &userPackage.UpdatedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return gerror.New("用户套餐不存在")
			}
			g.Log().Errorf(ctx, "查询用户套餐失败: %v", err)
			return err
		}

		// 3. 检查状态是否允许变更 (可以根据业务逻辑添加更复杂的校验)
		// 例如：已过期的套餐不能再改为激活状态
		if userPackage.Status == consts.PackageStatusExpired && req.Status == consts.PackageStatusActive {
			return gerror.New("已过期的套餐不能被激活")
		}
		// 如果状态没有变化，则直接返回
		if userPackage.Status == req.Status {
			return nil // 不需要更新
		}

		// 4. 更新状态
		updateData := g.Map{
			"status":     req.Status,
			"updated_at": gtime.Now(),
		}

		// 根据新的状态，可能需要调整开始/结束时间
		switch req.Status {
		case consts.PackageStatusActive:
			// 如果是从pending激活，需要设置开始和结束时间
			if userPackage.Status == consts.PackageStatusPending {
				var pkg entity.DrinkAllYouCanPackages
				err := tx.Model(dao.DrinkAllYouCanPackages.Table()).Where("id", userPackage.PackageId).Scan(&pkg)
				if err != nil {
					return gerror.Wrap(err, "获取套餐模板信息失败")
				}
				now := gtime.Now()
				updateData["start_time"] = now
				updateData["end_time"] = now.Add(time.Minute * time.Duration(pkg.DurationMinutes))
			}
		case consts.PackageStatusPending:
			// 如果改为待激活，清除开始和结束时间
			updateData["start_time"] = nil
			updateData["end_time"] = nil
		case "refunded":
			// 如果是退款，也应该清除时间，并可能记录原因（如果表有字段）
			updateData["start_time"] = nil
			updateData["end_time"] = nil
		}

		_, err = tx.Model(dao.UserPackages.Table()).
			Data(updateData).
			Where("id", userPackage.Id).
			Update()

		if err != nil {
			g.Log().Errorf(ctx, "更新用户套餐状态失败: %v", err)
			return err
		}

		// 可以在这里记录状态变更日志（如果需要）
		g.Log().Infof(ctx, "用户套餐ID %d 状态已从 %s 更新为 %s，原因: %s", userPackage.Id, userPackage.Status, req.Status, req.Reason)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 5. 返回更新后的完整信息
	return s.GetUserPackageDetail(ctx, req.UserPackageID)
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

// GetUserPackageFullDetail 获取用户套餐详细信息（包含用户信息、套餐信息和使用情况）
func (s *userPackageService) GetUserPackageFullDetail(ctx context.Context, userPackageID int64) (detail *v1.AdminUserPackageFullDetail, err error) {
	detail = &v1.AdminUserPackageFullDetail{}

	// 定义一个结构体来接收JOIN查询的结果
	type FullDetailResult struct {
		entity.UserPackages
		UserNickname    string      `gorm:"column:user_nickname"`
		UserAvatarUrl   string      `gorm:"column:user_avatar_url"`
		UserPhone       string      `gorm:"column:user_phone"`
		PackageName     string      `gorm:"column:package_name"`
		PackageDesc     string      `gorm:"column:package_desc"`
		PackagePrice    float64     `gorm:"column:package_price"`
		PackageDuration int         `gorm:"column:package_duration"`
		OrderSN         string      `gorm:"column:order_sn"`
		OrderTotalFee   float64     `gorm:"column:order_total_fee"`
		OrderPayStatus  string      `gorm:"column:order_pay_status"`
		OrderPaidAt     *gtime.Time `gorm:"column:order_paid_at"`
	}

	var result FullDetailResult

	// 优化1：使用LEFT JOIN一次性查询关联的所有信息
	err = g.DB().Model("user_packages up").
		Fields(
			"up.*",
			"u.nickname as user_nickname, u.avatar_url as user_avatar_url, u.phone as user_phone",
			"p.name as package_name, p.description as package_desc, p.price as package_price, p.duration_minutes as package_duration",
			"o.order_sn as order_sn, o.total_amount as order_total_fee, o.payment_status as order_pay_status, o.paid_at as order_paid_at",
		).
		LeftJoin("users u", "u.id = up.user_id").
		LeftJoin("drink_all_you_can_packages p", "p.id = up.package_id").
		LeftJoin("orders o", "o.id = up.order_id").
		Where("up.id", userPackageID).
		Scan(&result)

	if err != nil {
		g.Log().Errorf(ctx, "查询用户套餐完整信息失败, ID:%d, 错误:%v", userPackageID, err)
		return nil, err
	}

	if result.Id == 0 {
		return nil, gerror.New("用户套餐不存在")
	}

	// 实时检查并更新套餐状态
	if err = s.CheckAndUpdatePackageStatus(ctx, &result.UserPackages); err != nil {
		// 即使状态更新失败，也继续返回数据，只记录日志
		g.Log().Errorf(ctx, "检查更新套餐状态失败: %v", err)
	}

	// 优化2：将查询结果直接映射到返回的DTO结构中
	// 转换基本信息
	detail.ID = int64(result.Id)
	detail.UserID = int64(result.UserId)
	detail.PackageID = int64(result.PackageId)
	detail.OrderID = int64(result.OrderId)
	detail.Status = result.Status
	detail.CreatedAt = utility.FormatTimeOrEmpty(result.CreatedAt)
	detail.UpdatedAt = utility.FormatTimeOrEmpty(result.UpdatedAt)
	if result.StartTime != nil {
		detail.StartTime = utility.FormatTimeOrEmpty(result.StartTime)
	}
	if result.EndTime != nil {
		detail.EndTime = utility.FormatTimeOrEmpty(result.EndTime)
	}

	// 填充用户信息
	detail.User.ID = int64(result.UserId)
	detail.User.Nickname = result.UserNickname
	detail.User.AvatarUrl = result.UserAvatarUrl
	detail.User.Phone = result.UserPhone
	detail.UserName = result.UserNickname
	if detail.UserName == "" {
		detail.UserName = "用户" + gconv.String(result.UserId)
	}

	// 填充套餐信息
	detail.Package.ID = int64(result.PackageId)
	detail.Package.Name = result.PackageName
	detail.Package.Description = result.PackageDesc
	detail.Package.Price = result.PackagePrice
	detail.Package.DurationMinutes = result.PackageDuration
	detail.PackageName = result.PackageName

	// 填充订单信息
	if result.OrderId > 0 {
		detail.Order.ID = int64(result.OrderId)
		detail.Order.OrderSN = result.OrderSN
		detail.OrderSN = result.OrderSN
		detail.Order.TotalFee = result.OrderTotalFee
		detail.Order.PayStatus = result.OrderPayStatus
		detail.Order.PayTime = utility.FormatTimeOrEmpty(result.OrderPaidAt)
	}

	// 设置状态描述
	switch result.Status {
	case consts.PackageStatusActive:
		detail.StatusDesc = "已激活"
	case consts.PackageStatusExpired:
		detail.StatusDesc = "已过期"
	case consts.PackageStatusPending:
		detail.StatusDesc = "待激活"
	case "refunded":
		detail.StatusDesc = "已退款"
	default:
		detail.StatusDesc = result.Status
	}

	// 计算有效期描述
	if result.EndTime != nil {
		if result.StartTime != nil && result.Status == consts.PackageStatusActive {
			duration := result.EndTime.Sub(gtime.Now())
			days := int(duration.Hours()) / 24
			hours := int(duration.Hours()) % 24
			if days > 0 {
				detail.ValidPeriod = gconv.String(days) + "天"
				if hours > 0 {
					detail.ValidPeriod += gconv.String(hours) + "小时"
				}
			} else if hours > 0 {
				detail.ValidPeriod = gconv.String(hours) + "小时"
			} else {
				detail.ValidPeriod = "即将到期"
			}
		} else if result.Status == consts.PackageStatusExpired {
			detail.ValidPeriod = "已过期"
		} else if result.Status == consts.PackageStatusPending {
			if result.PackageDuration > 0 {
				detail.ValidPeriod = gconv.String(result.PackageDuration/60) + "小时"
			} else {
				detail.ValidPeriod = "-"
			}
		} else {
			detail.ValidPeriod = "-"
		}
	} else {
		detail.ValidPeriod = "-"
	}

	// 查询使用情况统计
	// 这里假设有一个表记录了套餐的使用情况，如果没有可以根据实际情况调整
	type UsageRecord struct {
		TotalUsedTimes int         `json:"total_used_times"`
		LastUsedTime   *gtime.Time `json:"last_used_time"`
	}

	var usageRecord UsageRecord
	// 这里需要根据实际的数据库结构查询使用记录
	// 优化3：保持这个独立的聚合查询，因为它与主查询的粒度不同
	err = g.DB().Model("order_items").
		Where("user_package_id", userPackageID).
		Fields("COUNT(1) as total_used_times, MAX(created_at) as last_used_time").
		Scan(&usageRecord)
	if err != nil {
		// 查询使用记录失败不应阻塞主流程
		g.Log().Warningf(ctx, "查询套餐使用记录失败, UserPackageID: %d, Error: %v", userPackageID, err)
	} else {
		detail.Usage.TotalUsedTimes = usageRecord.TotalUsedTimes
		detail.Usage.LastUsedTime = utility.FormatTimeOrEmpty(usageRecord.LastUsedTime)
	}

	return detail, nil
}
