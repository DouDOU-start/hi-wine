package service

import (
	"context"

	v1 "backend/api/user/v1"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/utility"
	"backend/internal/utility/jwt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserService 用户服务接口
type UserService interface {
	// 用户登录 (通过微信授权)
	LoginByWechat(ctx context.Context, code string) (user *v1.UserProfile, token string, err error)

	// 获取用户信息
	GetUserInfo(ctx context.Context, userId int64) (user *v1.UserProfile, err error)

	// 更新用户信息
	UpdateUserInfo(ctx context.Context, userId int64, req *v1.UpdateUserProfileReq) (user *v1.UserProfile, err error)
}

// UserPackageServiceForUser 用户端套餐服务接口
type UserPackageServiceForUser interface {
	// 获取可购买套餐列表
	GetPackageList(ctx context.Context, req *v1.UserPackageListReq) (list []v1.UserPackage, total int, err error)

	// 获取套餐详情
	GetPackageDetail(ctx context.Context, packageId int64) (*v1.UserPackageDetail, error)

	// 购买套餐
	// 创建套餐购买订单，套餐状态为待支付(pending)，支付成功后才会激活套餐
	// 在套餐时效内，用户不可重复购买套餐；如有待支付套餐，也不能购买新套餐
	BuyPackage(ctx context.Context, userId int64, packageId int64) (orderId int64, err error)
}

// userService 用户服务实现
type userService struct{}

// userPackageServiceForUser 用户端套餐服务实现
type userPackageServiceForUser struct{}

// 确保实现了接口
var _ UserService = (*userService)(nil)
var _ UserPackageServiceForUser = (*userPackageServiceForUser)(nil)

// 单例实例
var (
	userServiceInstance        = userService{}
	userPackageForUserInstance = userPackageServiceForUser{}
)

// User 获取用户服务实例
func User() UserService {
	return &userServiceInstance
}

// UserPackageForUser 获取用户套餐服务实例
func UserPackageForUser() UserPackageServiceForUser {
	return &userPackageForUserInstance
}

// LoginByWechat 通过微信登录
func (s *userService) LoginByWechat(ctx context.Context, code string) (user *v1.UserProfile, token string, err error) {
	// 这里实现微信登录逻辑
	// 1. 调用微信API获取用户openid
	// 2. 检查用户是否存在，不存在则创建
	// 3. 生成JWT token
	// 此处为示例代码，实际项目中需要替换为真实的微信API调用
	openid := "test_openid_" + code // 模拟openid

	// 查询用户是否存在
	var userEntity *entity.Users
	err = dao.Users.Ctx(ctx).Where("openid", openid).Scan(&userEntity)
	if err != nil {
		return nil, "", err
	}

	// 用户不存在，创建新用户
	if userEntity == nil {
		userEntity = &entity.Users{
			Openid:    openid,
			CreatedAt: gtime.Now(),
			UpdatedAt: gtime.Now(),
		}
		_, err = dao.Users.Ctx(ctx).Insert(userEntity)
		if err != nil {
			return nil, "", err
		}

		// 获取新创建的用户ID
		err = dao.Users.Ctx(ctx).Where("openid", openid).Scan(&userEntity)
		if err != nil {
			return nil, "", err
		}
	}

	// 生成token
	token, err = jwt.GenerateToken(int64(userEntity.Id))
	if err != nil {
		return nil, "", err
	}

	// 构建返回用户信息
	user = &v1.UserProfile{
		ID:        int64(userEntity.Id),
		Nickname:  userEntity.Nickname,
		AvatarURL: userEntity.AvatarUrl,
		Phone:     userEntity.Phone,
	}

	return user, token, nil
}

// GetUserInfo 获取用户信息
func (s *userService) GetUserInfo(ctx context.Context, userId int64) (user *v1.UserProfile, err error) {
	var userEntity *entity.Users
	err = dao.Users.Ctx(ctx).Where("id", userId).Scan(&userEntity)
	if err != nil {
		return nil, err
	}
	if userEntity == nil {
		return nil, gerror.New("用户不存在")
	}

	user = &v1.UserProfile{
		ID:        int64(userEntity.Id),
		Nickname:  userEntity.Nickname,
		AvatarURL: userEntity.AvatarUrl,
		Phone:     userEntity.Phone,
	}

	return user, nil
}

// UpdateUserInfo 更新用户信息
func (s *userService) UpdateUserInfo(ctx context.Context, userId int64, req *v1.UpdateUserProfileReq) (user *v1.UserProfile, err error) {
	// 构建更新数据
	updateData := g.Map{
		"updated_at": gtime.Now(),
	}

	if req.Nickname != "" {
		updateData["nickname"] = req.Nickname
	}
	if req.AvatarURL != "" {
		updateData["avatar_url"] = req.AvatarURL
	}
	if req.Phone != "" {
		updateData["phone"] = req.Phone
	}

	// 更新用户信息
	_, err = dao.Users.Ctx(ctx).Where("id", userId).Data(updateData).Update()
	if err != nil {
		return nil, err
	}

	// 获取更新后的用户信息
	return s.GetUserInfo(ctx, userId)
}

// GetPackageList 获取可购买套餐列表
func (s *userPackageServiceForUser) GetPackageList(ctx context.Context, req *v1.UserPackageListReq) (list []v1.UserPackage, total int, err error) {
	// 构建查询条件
	m := g.DB().Model(dao.DrinkAllYouCanPackages.Table()).
		Where(dao.DrinkAllYouCanPackages.Columns().IsActive, 1) // 只查询激活的套餐

	// 添加名称模糊搜索
	if req.Name != "" {
		m = m.WhereLike(dao.DrinkAllYouCanPackages.Columns().Name, "%"+req.Name+"%")
	}

	// 查询总数
	total, err = m.Count()
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
	err = m.Order(dao.DrinkAllYouCanPackages.Columns().Id+" DESC").
		Limit((page-1)*limit, limit).
		Scan(&packages)
	if err != nil {
		return nil, 0, err
	}

	// 转换为API响应格式
	list = make([]v1.UserPackage, 0, len(packages))
	for _, item := range packages {
		list = append(list, v1.UserPackage{
			ID:            int64(item.Id),
			Name:          item.Name,
			Description:   item.Description,
			Price:         item.Price,
			DurationHours: item.DurationMinutes / 60, // 转换分钟为小时
			IsActive:      item.IsActive == 1,
		})
	}

	return list, total, nil
}

// GetPackageDetail 获取套餐详情
func (s *userPackageServiceForUser) GetPackageDetail(ctx context.Context, packageId int64) (*v1.UserPackageDetail, error) {
	// 查询套餐信息
	var packageInfo entity.DrinkAllYouCanPackages
	err := g.DB().Model(dao.DrinkAllYouCanPackages.Table()).
		Where(dao.DrinkAllYouCanPackages.Columns().Id, packageId).
		Where(dao.DrinkAllYouCanPackages.Columns().IsActive, 1). // 只查询激活的套餐
		Scan(&packageInfo)
	if err != nil {
		return nil, err
	}
	if packageInfo.Id == 0 {
		return nil, gerror.New("套餐不存在或已下架")
	}

	// 构建返回结果
	detail := &v1.UserPackageDetail{
		UserPackage: v1.UserPackage{
			ID:            int64(packageInfo.Id),
			Name:          packageInfo.Name,
			Description:   packageInfo.Description,
			Price:         packageInfo.Price,
			DurationHours: packageInfo.DurationMinutes / 60, // 转换分钟为小时
			IsActive:      packageInfo.IsActive == 1,
		},
	}

	return detail, nil
}

// BuyPackage 购买套餐
func (s *userPackageServiceForUser) BuyPackage(ctx context.Context, userId int64, packageId int64) (orderId int64, err error) {
	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 0.1 检查用户是否已有有效套餐
		var activePackageCount int
		activePackageCount, err = tx.Model(dao.UserPackages.Table()).
			Where(dao.UserPackages.Columns().UserId, userId).
			Where(dao.UserPackages.Columns().Status, "active").
			Where(dao.UserPackages.Columns().EndTime+" > ?", gtime.Now()).
			Count()
		if err != nil {
			return err
		}

		// 如果用户已有有效套餐，则不允许重复购买
		if activePackageCount > 0 {
			return gerror.New("您已有有效套餐，不能重复购买")
		}

		// 0.2 检查用户是否有待支付的套餐
		var pendingPackageCount int
		pendingPackageCount, err = tx.Model(dao.UserPackages.Table()).
			Where(dao.UserPackages.Columns().UserId, userId).
			Where(dao.UserPackages.Columns().Status, "pending").
			Count()
		if err != nil {
			return err
		}

		// 如果用户有待支付的套餐，则提示用户先完成支付
		if pendingPackageCount > 0 {
			return gerror.New("您有待支付的套餐订单，请先完成支付")
		}

		// 1. 查询套餐信息
		var packageInfo entity.DrinkAllYouCanPackages
		err = tx.Model(dao.DrinkAllYouCanPackages.Table()).
			Where(dao.DrinkAllYouCanPackages.Columns().Id, packageId).
			Where(dao.DrinkAllYouCanPackages.Columns().IsActive, 1). // 只允许购买激活的套餐
			Scan(&packageInfo)
		if err != nil {
			return err
		}
		if packageInfo.Id == 0 {
			return gerror.New("套餐不存在或已下架")
		}

		// 2. 验证用户是否存在
		var userCount int
		userCount, err = tx.Model(dao.Users.Table()).Where(dao.Users.Columns().Id, userId).Count()
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
		// 这里调用已创建的CreateUserPackageAfterPurchase方法
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
