package service

import (
	"context"

	v1 "backend/api/user/v1"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/utility"
	"backend/internal/utility/jwt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserService 用户服务接口
type UserService interface {
	// 用户登录 (通过微信授权)
	LoginByWechat(ctx context.Context, code string, nickname string, avatarURL string) (user *v1.UserProfile, token string, err error)

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
func (s *userService) LoginByWechat(ctx context.Context, code string, nickname string, avatarURL string) (user *v1.UserProfile, token string, err error) {
	// 测试模式：当code为test_code时，返回测试用户
	if code == "test_code" {
		g.Log().Info(ctx, "使用测试模式登录")

		// 查询或创建测试用户
		var testUser *entity.Users
		err = dao.Users.Ctx(ctx).Where("openid", "test_openid").Scan(&testUser)
		if err != nil {
			return nil, "", gerror.New("查询测试用户失败: " + err.Error())
		}

		if testUser == nil {
			// 创建测试用户
			testNickname := "测试用户"
			if nickname != "" {
				testNickname = nickname
			}

			testAvatar := "https://example.com/avatar.png"
			if avatarURL != "" {
				testAvatar = avatarURL
			}

			testUser = &entity.Users{
				Openid:    "test_openid",
				Nickname:  testNickname,
				AvatarUrl: testAvatar,
				CreatedAt: gtime.Now(),
				UpdatedAt: gtime.Now(),
			}

			result, err := dao.Users.Ctx(ctx).Insert(testUser)
			if err != nil {
				return nil, "", gerror.New("创建测试用户失败: " + err.Error())
			}

			lastInsertId, err := result.LastInsertId()
			if err != nil {
				return nil, "", gerror.New("获取测试用户ID失败: " + err.Error())
			}
			testUser.Id = int(lastInsertId)

			g.Log().Info(ctx, "创建测试用户成功", g.Map{
				"userId": testUser.Id,
				"openid": testUser.Openid,
			})
		}

		// 生成token
		token, err = jwt.GenerateToken(g.Map{
			"userId": testUser.Id,
			"openid": testUser.Openid,
		})
		if err != nil {
			return nil, "", gerror.New("生成token失败: " + err.Error())
		}

		// 构建返回用户信息
		user = &v1.UserProfile{
			ID:        int64(testUser.Id),
			Openid:    testUser.Openid,
			Nickname:  testUser.Nickname,
			AvatarURL: testUser.AvatarUrl,
			Phone:     testUser.Phone,
		}

		g.Log().Debug(ctx, "测试登录成功", g.Map{
			"userId": user.ID,
			"token":  token,
		})

		return user, token, nil
	}

	// 正常模式：调用微信API获取用户openid
	// 微信小程序登录API: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
	appID := g.Cfg().MustGet(ctx, "wechat.appid").String()
	appSecret := g.Cfg().MustGet(ctx, "wechat.secret").String()

	g.Log().Debug(ctx, "微信登录参数：", g.Map{
		"code":      code,
		"appID":     appID,
		"appSecret": appSecret[:4] + "****", // 仅记录前几位，保护密钥安全
	})

	if appID == "" || appSecret == "" {
		return nil, "", gerror.New("微信小程序配置缺失")
	}

	// 构建请求URL
	url := "https://api.weixin.qq.com/sns/jscode2session"
	params := g.Map{
		"appid":      appID,
		"secret":     appSecret,
		"js_code":    code,
		"grant_type": "authorization_code",
	}

	// 发送请求到微信API
	response, err := g.Client().Get(ctx, url, params)
	if err != nil {
		g.Log().Error(ctx, "调用微信API失败:", err)
		return nil, "", gerror.New("调用微信API失败: " + err.Error())
	}
	defer response.Close()

	// 读取响应内容
	responseBody := response.ReadAllString()
	g.Log().Debug(ctx, "微信API响应:", responseBody)

	// 解析响应为JSON
	j, err := gjson.DecodeToJson(responseBody)
	if err != nil {
		g.Log().Error(ctx, "解析微信API响应失败:", err)
		return nil, "", gerror.New("解析微信API响应失败: " + err.Error())
	}
	wxResp := j.Map()

	// 检查是否有错误
	if wxResp["errcode"] != nil && wxResp["errcode"].(float64) != 0 {
		errMsg := g.Map{
			"errcode": wxResp["errcode"],
			"errmsg":  wxResp["errmsg"],
		}
		g.Log().Error(ctx, "微信登录返回错误:", errMsg)
		return nil, "", gerror.Newf("微信登录失败: %s", wxResp["errmsg"])
	}

	// 获取openid
	openid := wxResp["openid"].(string)
	if openid == "" {
		return nil, "", gerror.New("获取openid失败")
	}

	// 查询用户是否存在
	var userEntity *entity.Users
	err = dao.Users.Ctx(ctx).Where("openid", openid).Scan(&userEntity)
	if err != nil {
		return nil, "", gerror.New("查询用户失败: " + err.Error())
	}

	// 用户不存在，创建新用户
	if userEntity == nil {
		// 创建新用户，保存昵称和头像信息
		userEntity = &entity.Users{
			Openid:    openid,
			Nickname:  nickname,
			AvatarUrl: avatarURL,
			CreatedAt: gtime.Now(),
			UpdatedAt: gtime.Now(),
		}

		// 如果微信返回了unionid，也可以保存
		if unionid, ok := wxResp["unionid"].(string); ok && unionid != "" {
			// 如果Users表中有unionid字段，可以保存
			// userEntity.Unionid = unionid
		}

		result, err := dao.Users.Ctx(ctx).Insert(userEntity)
		if err != nil {
			return nil, "", gerror.New("创建用户失败: " + err.Error())
		}

		// 获取新创建的用户ID
		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return nil, "", gerror.New("获取用户ID失败: " + err.Error())
		}
		userEntity.Id = int(lastInsertId)

		g.Log().Info(ctx, "新用户注册成功", g.Map{
			"userId":    userEntity.Id,
			"nickname":  nickname,
			"hasAvatar": avatarURL != "",
		})
	} else if nickname != "" || avatarURL != "" {
		// 用户已存在，但传入了昵称或头像信息，可能是首次授权
		updateData := g.Map{
			"updated_at": gtime.Now(),
		}

		needUpdate := false

		// 只有当数据库中昵称为空且传入了昵称时，才更新昵称
		if nickname != "" && userEntity.Nickname == "" {
			updateData["nickname"] = nickname
			userEntity.Nickname = nickname
			needUpdate = true
		}

		// 只有当数据库中头像为空且传入了头像时，才更新头像
		if avatarURL != "" && userEntity.AvatarUrl == "" {
			updateData["avatar_url"] = avatarURL
			userEntity.AvatarUrl = avatarURL
			needUpdate = true
		}

		if needUpdate {
			_, err = dao.Users.Ctx(ctx).Where("id", userEntity.Id).Data(updateData).Update()
			if err != nil {
				g.Log().Error(ctx, "更新用户信息失败:", err)
				// 这里我们不返回错误，因为主要功能是登录，更新信息失败不应影响登录流程
			} else {
				g.Log().Info(ctx, "更新用户信息成功", g.Map{
					"userId":   userEntity.Id,
					"nickname": userEntity.Nickname != "",
					"avatar":   userEntity.AvatarUrl != "",
				})
			}
		}
	} else {
		g.Log().Debug(ctx, "用户登录成功，无需更新信息", g.Map{
			"userId":      userEntity.Id,
			"hasNickname": userEntity.Nickname != "",
			"hasAvatar":   userEntity.AvatarUrl != "",
		})
	}

	// 生成token
	token, err = jwt.GenerateToken(g.Map{
		"userId": userEntity.Id,
		"openid": userEntity.Openid,
	})
	if err != nil {
		return nil, "", gerror.New("生成token失败: " + err.Error())
	}

	// 构建返回用户信息
	user = &v1.UserProfile{
		ID:        int64(userEntity.Id),
		Openid:    userEntity.Openid,
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
		g.Log().Error(ctx, "查询用户信息失败:", err)
		return nil, err
	}
	if userEntity == nil {
		return nil, gerror.New("用户不存在")
	}

	user = &v1.UserProfile{
		ID:        int64(userEntity.Id),
		Openid:    userEntity.Openid,
		Nickname:  userEntity.Nickname,
		AvatarURL: userEntity.AvatarUrl,
		Phone:     userEntity.Phone,
	}

	g.Log().Debug(ctx, "获取用户信息:", g.Map{
		"userId":   user.ID,
		"openid":   user.Openid,
		"nickname": user.Nickname,
		"avatar":   user.AvatarURL != "",
		"phone":    user.Phone != "",
	})

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
