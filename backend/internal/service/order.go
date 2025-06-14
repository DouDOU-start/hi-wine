package service

import (
	"context"
	"fmt"

	adminv1 "backend/api/admin/v1"
	orderv1 "backend/api/order/v1"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/utility"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderService 订单服务接口
type OrderService interface {
	// 用户端接口
	// CreateOrder 创建订单
	CreateOrder(ctx context.Context, req *orderv1.CreateOrderReq, userId int64) (*orderv1.CreateOrderRes, error)

	// GetOrderDetail 获取订单详情
	GetOrderDetail(ctx context.Context, orderId int64, userId int64) (*orderv1.Order, error)

	// GetUserOrderList 获取用户订单列表
	GetUserOrderList(ctx context.Context, req *UserOrderListReq, userId int64) ([]orderv1.Order, int, error)

	// 管理端接口
	// GetOrderList 获取订单列表（分页、筛选、模糊搜索）
	GetOrderList(ctx context.Context, req *adminv1.AdminOrderListReq) ([]adminv1.AdminOrder, int, error)

	// GetOrderDetailAdmin 管理端获取订单详情
	GetOrderDetailAdmin(ctx context.Context, orderId int64) (*adminv1.AdminOrder, error)

	// UpdateOrderStatus 更新订单状态
	UpdateOrderStatus(ctx context.Context, orderId int64, status string, reason string) (*adminv1.AdminOrder, error)
}

// UserOrderListReq 用户订单列表请求参数
type UserOrderListReq struct {
	Status string
	Page   int
	Limit  int
}

// 单例实例
var orderServiceInstance = orderService{}

// Order 获取订单服务实例
func Order() OrderService {
	return &orderServiceInstance
}

// 订单服务实现
type orderService struct{}

// CreateOrder 创建订单
func (s *orderService) CreateOrder(ctx context.Context, req *orderv1.CreateOrderReq, userId int64) (*orderv1.CreateOrderRes, error) {
	// 1. 验证桌号二维码是否存在
	var tableQrcode *entity.TableQrcodes
	err := dao.TableQrcodes.Ctx(ctx).Where(dao.TableQrcodes.Columns().Id, req.TableQrcodeID).Scan(&tableQrcode)
	if err != nil {
		return nil, err
	}
	if tableQrcode == nil {
		return nil, gerror.New("桌号二维码不存在")
	}

	// 2. 验证订单商品
	if len(req.Items) == 0 {
		return nil, gerror.New("订单商品不能为空")
	}

	// 3. 开启事务
	var orderRes *orderv1.CreateOrderRes
	err = dao.Orders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 3.1 生成订单号
		orderSN := utility.GenerateOrderSN()

		// 3.2 创建订单主表记录
		order := &entity.Orders{
			OrderSn:       orderSN,
			UserId:        int(userId),
			TableQrcodeId: int(req.TableQrcodeID),
			TotalAmount:   0, // 先设置为0，后面计算
			PaymentStatus: "pending",
			OrderStatus:   "new",
			PaymentMethod: "",
			TransactionId: "",
			CreatedAt:     gtime.Now(),
			UpdatedAt:     gtime.Now(),
		}

		// 3.3 插入订单
		result, err := tx.Model(dao.Orders.Table()).Data(order).Insert()
		if err != nil {
			return err
		}

		// 3.4 获取订单ID
		orderId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// 3.5 检查用户是否有有效的畅饮套餐
		// 注意：只有状态为 active 的套餐才能被用于畅饮，pending 状态的套餐（未支付）不能使用
		// 用户必须先完成支付，套餐状态变为 active 后才能享受畅饮
		var userPackage *entity.UserPackages
		err = tx.Model(dao.UserPackages.Table()).
			Where(dao.UserPackages.Columns().UserId, userId).
			Where(dao.UserPackages.Columns().Status, "active").            // 只查询已激活的套餐
			Where(dao.UserPackages.Columns().EndTime+" > ?", gtime.Now()). // 确保套餐未过期
			Scan(&userPackage)
		if err != nil {
			return err
		}

		// 3.6 处理订单商品
		var totalAmount float64 = 0
		for _, item := range req.Items {
			// 3.6.1 查询商品信息
			var product *entity.Products
			err = tx.Model(dao.Products.Table()).
				Where(dao.Products.Columns().Id, item.ProductID).
				Where(dao.Products.Columns().IsActive, 1).
				Scan(&product)
			if err != nil {
				return err
			}
			if product == nil {
				return gerror.Newf("商品ID %d 不存在或已下架", item.ProductID)
			}

			// 3.6.2 检查库存
			if product.Stock < item.Quantity {
				return gerror.Newf("商品 %s 库存不足", product.Name)
			}

			// 3.6.3 判断是否为套餐内商品
			isPackageItem := 0
			var userPackageId interface{} = nil // 使用nil而不是0
			itemPrice := product.Price

			if userPackage != nil {
				// 查询该商品是否在用户当前有效套餐内
				var packageProduct *entity.PackageProducts
				err = tx.Model(dao.PackageProducts.Table()).
					Where(dao.PackageProducts.Columns().PackageId, userPackage.PackageId).
					Where(dao.PackageProducts.Columns().ProductId, product.Id).
					Scan(&packageProduct)
				if err != nil {
					return err
				}

				// 如果是套餐内商品，则价格为0
				if packageProduct != nil {
					isPackageItem = 1
					userPackageId = userPackage.Id
					itemPrice = 0
				}
			}

			// 3.6.4 创建订单商品记录
			orderItemData := g.Map{
				dao.OrderItems.Columns().OrderId:       int(orderId),
				dao.OrderItems.Columns().ProductId:     product.Id,
				dao.OrderItems.Columns().ProductName:   product.Name,
				dao.OrderItems.Columns().Price:         product.Price,
				dao.OrderItems.Columns().Quantity:      item.Quantity,
				dao.OrderItems.Columns().Subtotal:      product.Price * float64(item.Quantity),
				dao.OrderItems.Columns().IsPackageItem: isPackageItem,
				dao.OrderItems.Columns().UserPackageId: userPackageId,
				dao.OrderItems.Columns().ItemPrice:     itemPrice * float64(item.Quantity),
				dao.OrderItems.Columns().CreatedAt:     gtime.Now(),
				dao.OrderItems.Columns().UpdatedAt:     gtime.Now(),
			}

			// 插入订单商品
			_, err = tx.Model(dao.OrderItems.Table()).Data(orderItemData).Insert()
			if err != nil {
				return err
			}

			// 计算订单总金额（实际支付金额）
			totalAmount += itemPrice * float64(item.Quantity)

			// 3.6.5 扣减库存
			_, err = tx.Model(dao.Products.Table()).
				Where(dao.Products.Columns().Id, product.Id).
				Data(g.Map{
					dao.Products.Columns().Stock: gdb.Raw(fmt.Sprintf("%s - %d", dao.Products.Columns().Stock, item.Quantity)),
				}).
				Update()
			if err != nil {
				return err
			}
		}

		// 3.7 更新订单总金额
		_, err = tx.Model(dao.Orders.Table()).
			Where(dao.Orders.Columns().Id, orderId).
			Data(g.Map{
				dao.Orders.Columns().TotalAmount: totalAmount,
			}).
			Update()
		if err != nil {
			return err
		}

		// 3.8 设置返回结果
		orderRes = &orderv1.CreateOrderRes{
			OrderID:     orderId,
			OrderSN:     orderSN,
			TotalAmount: totalAmount,
			PrepayID:    "", // 微信支付的预支付ID，需要调用支付接口获取
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return orderRes, nil
}

// GetOrderDetail 获取订单详情
func (s *orderService) GetOrderDetail(ctx context.Context, orderId int64, userId int64) (*orderv1.Order, error) {
	// 1. 查询订单基本信息
	var order *entity.Orders
	err := dao.Orders.Ctx(ctx).
		Where(dao.Orders.Columns().Id, orderId).
		Where(dao.Orders.Columns().UserId, userId).
		Scan(&order)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.New("订单不存在")
	}

	// 2. 查询订单商品
	var orderItems []*entity.OrderItems
	err = dao.OrderItems.Ctx(ctx).
		Where(dao.OrderItems.Columns().OrderId, orderId).
		Scan(&orderItems)
	if err != nil {
		return nil, err
	}

	// 3. 转换为API响应格式
	result := &orderv1.Order{
		ID:            int64(order.Id),
		OrderSN:       order.OrderSn,
		Status:        fmt.Sprintf("%s_%s", order.PaymentStatus, order.OrderStatus),
		TotalAmount:   order.TotalAmount,
		CreatedAt:     order.CreatedAt.String(),
		TableQrcodeID: int64(order.TableQrcodeId),
		TotalNotes:    "", // 订单备注字段，如果需要可以添加到数据库表中
		Items:         make([]orderv1.OrderItem, len(orderItems)),
	}

	// 4. 转换订单商品
	for i, item := range orderItems {
		// 查询商品图片
		var product *entity.Products
		imageURL := ""
		if item.ProductId > 0 {
			err = dao.Products.Ctx(ctx).
				Where(dao.Products.Columns().Id, item.ProductId).
				Fields("image_url").
				Scan(&product)
			if err == nil && product != nil {
				imageURL = product.ImageUrl
			}
		}

		result.Items[i] = orderv1.OrderItem{
			ProductID:     int64(item.ProductId),
			Name:          item.ProductName,
			ImageURL:      imageURL,
			Quantity:      item.Quantity,
			ItemPrice:     item.ItemPrice,
			IsPackageItem: item.IsPackageItem == 1,
			UserPackageID: int64(item.UserPackageId),
			Notes:         "", // 商品备注字段，如果需要可以添加到数据库表中
		}

		// 如果是套餐商品，查询套餐详情
		if item.IsPackageItem == 1 && item.UserPackageId > 0 {
			// 查询用户套餐详情
			packageDetail, err := UserPackage().GetUserPackageDetail(ctx, int64(item.UserPackageId))
			if err == nil && packageDetail != nil {
				// 添加套餐状态和有效期信息
				result.Items[i].PackageInfo = &orderv1.PackageInfo{
					Status:    packageDetail.Status,
					StartTime: packageDetail.StartTime,
					EndTime:   packageDetail.EndTime,
				}
			}
		}
	}

	// 5. 如果订单商品为空，但是订单是套餐订单，则查询关联的套餐信息
	if len(orderItems) == 0 {
		// 查询是否有关联的套餐
		var userPackage entity.UserPackages
		err = g.DB().Model("user_packages").
			Where("order_id", orderId).
			Scan(&userPackage)

		if err == nil && userPackage.Id > 0 {
			// 查询套餐详情
			packageInfo, err := UserPackage().GetUserPackageDetail(ctx, int64(userPackage.Id))
			if err == nil && packageInfo != nil {
				// 查询套餐商品信息
				var packageDetail entity.DrinkAllYouCanPackages
				err = g.DB().Model("drink_all_you_can_packages").
					Where("id", userPackage.PackageId).
					Scan(&packageDetail)

				if err == nil && packageDetail.Id > 0 {
					// 创建一个套餐商品项
					packageItem := orderv1.OrderItem{
						ProductID:     int64(packageDetail.Id),
						Name:          packageDetail.Name,
						ImageURL:      "", // 套餐暂无图片，设置为空字符串
						Quantity:      1,
						ItemPrice:     packageDetail.Price,
						IsPackageItem: true,
						UserPackageID: int64(userPackage.Id),
						Notes:         "",
						PackageInfo: &orderv1.PackageInfo{
							Status:    packageInfo.Status,
							StartTime: packageInfo.StartTime,
							EndTime:   packageInfo.EndTime,
						},
					}

					// 添加到订单商品列表
					result.Items = append(result.Items, packageItem)
				}
			}
		}
	}

	return result, nil
}

// GetUserOrderList 获取用户订单列表
func (s *orderService) GetUserOrderList(ctx context.Context, req *UserOrderListReq, userId int64) ([]orderv1.Order, int, error) {
	// 1. 构建查询条件
	m := dao.Orders.Ctx(ctx).Where(dao.Orders.Columns().UserId, userId)

	// 1.1 订单状态筛选
	if req.Status != "" {
		// 状态格式为 payment_status_order_status
		statuses := utility.SplitOrderStatus(req.Status)
		if statuses.PaymentStatus != "" {
			m = m.Where(dao.Orders.Columns().PaymentStatus, statuses.PaymentStatus)
		}
		if statuses.OrderStatus != "" {
			m = m.Where(dao.Orders.Columns().OrderStatus, statuses.OrderStatus)
		}
	}

	// 2. 查询总数
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	// 3. 分页参数处理
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}

	// 4. 查询订单数据
	var orders []*entity.Orders
	err = m.Page(page, limit).
		Order(dao.Orders.Columns().Id + " DESC").
		Scan(&orders)
	if err != nil {
		return nil, 0, err
	}

	// 5. 获取订单商品
	result := make([]orderv1.Order, len(orders))
	for i, order := range orders {
		// 5.1 查询订单商品
		var orderItems []*entity.OrderItems
		err = dao.OrderItems.Ctx(ctx).
			Where(dao.OrderItems.Columns().OrderId, order.Id).
			Scan(&orderItems)
		if err != nil {
			return nil, 0, err
		}

		// 5.2 转换为API响应格式
		result[i] = orderv1.Order{
			ID:            int64(order.Id),
			OrderSN:       order.OrderSn,
			Status:        fmt.Sprintf("%s_%s", order.PaymentStatus, order.OrderStatus),
			TotalAmount:   order.TotalAmount,
			CreatedAt:     order.CreatedAt.String(),
			TableQrcodeID: int64(order.TableQrcodeId),
			TotalNotes:    "", // 订单备注字段，如果需要可以添加到数据库表中
			Items:         make([]orderv1.OrderItem, len(orderItems)),
		}

		// 5.3 转换订单商品
		for j, item := range orderItems {
			// 查询商品图片
			var product *entity.Products
			imageURL := ""
			if item.ProductId > 0 {
				err = dao.Products.Ctx(ctx).
					Where(dao.Products.Columns().Id, item.ProductId).
					Fields("image_url").
					Scan(&product)
				if err == nil && product != nil {
					imageURL = product.ImageUrl
				}
			}

			result[i].Items[j] = orderv1.OrderItem{
				ProductID:     int64(item.ProductId),
				Name:          item.ProductName,
				ImageURL:      imageURL,
				Quantity:      item.Quantity,
				ItemPrice:     item.ItemPrice,
				IsPackageItem: item.IsPackageItem == 1,
				UserPackageID: int64(item.UserPackageId),
				Notes:         "", // 商品备注字段，如果需要可以添加到数据库表中
			}

			// 如果是套餐商品，查询套餐详情
			if item.IsPackageItem == 1 && item.UserPackageId > 0 {
				// 查询用户套餐详情
				packageDetail, err := UserPackage().GetUserPackageDetail(ctx, int64(item.UserPackageId))
				if err == nil && packageDetail != nil {
					// 添加套餐状态和有效期信息
					result[i].Items[j].PackageInfo = &orderv1.PackageInfo{
						Status:    packageDetail.Status,
						StartTime: packageDetail.StartTime,
						EndTime:   packageDetail.EndTime,
					}
				}
			}
		}

		// 5.4 如果订单项为空，但是订单是套餐订单，则查询关联的套餐信息
		if len(orderItems) == 0 {
			// 查询是否有关联的套餐
			var userPackage entity.UserPackages
			err = g.DB().Model("user_packages").
				Where("order_id", order.Id).
				Scan(&userPackage)

			if err == nil && userPackage.Id > 0 {
				// 查询套餐详情
				packageInfo, err := UserPackage().GetUserPackageDetail(ctx, int64(userPackage.Id))
				if err == nil && packageInfo != nil {
					// 查询套餐商品信息
					var packageDetail entity.DrinkAllYouCanPackages
					err = g.DB().Model("drink_all_you_can_packages").
						Where("id", userPackage.PackageId).
						Scan(&packageDetail)

					if err == nil && packageDetail.Id > 0 {
						// 创建一个套餐商品项
						packageItem := orderv1.OrderItem{
							ProductID:     int64(packageDetail.Id),
							Name:          packageDetail.Name,
							ImageURL:      "", // 套餐暂无图片，设置为空字符串
							Quantity:      1,
							ItemPrice:     packageDetail.Price,
							IsPackageItem: true,
							UserPackageID: int64(userPackage.Id),
							Notes:         "",
							PackageInfo: &orderv1.PackageInfo{
								Status:    packageInfo.Status,
								StartTime: packageInfo.StartTime,
								EndTime:   packageInfo.EndTime,
							},
						}

						// 添加到订单商品列表
						result[i].Items = append(result[i].Items, packageItem)
					}
				}
			}
		}
	}

	return result, total, nil
}

// GetOrderList 获取订单列表（分页、筛选、模糊搜索）
func (s *orderService) GetOrderList(ctx context.Context, req *adminv1.AdminOrderListReq) ([]adminv1.AdminOrder, int, error) {
	// 1. 构建查询条件
	m := dao.Orders.Ctx(ctx)

	// 1.1 订单状态筛选
	if req.Status != "" {
		// 状态格式为 payment_status_order_status
		statuses := utility.SplitOrderStatus(req.Status)
		if statuses.PaymentStatus != "" {
			m = m.Where(dao.Orders.Columns().PaymentStatus, statuses.PaymentStatus)
		}
		if statuses.OrderStatus != "" {
			m = m.Where(dao.Orders.Columns().OrderStatus, statuses.OrderStatus)
		}
	}

	// 1.2 订单号搜索
	if req.OrderSN != "" {
		m = m.WhereLike(dao.Orders.Columns().OrderSn, "%"+req.OrderSN+"%")
	}

	// 1.3 日期范围筛选
	if req.StartDate != "" {
		m = m.WhereGTE(dao.Orders.Columns().CreatedAt, req.StartDate+" 00:00:00")
	}
	if req.EndDate != "" {
		m = m.WhereLTE(dao.Orders.Columns().CreatedAt, req.EndDate+" 23:59:59")
	}

	// 1.4 桌号搜索
	if req.TableNumber != "" {
		// 先查询桌号对应的二维码ID
		var tableQrcodes []*entity.TableQrcodes
		err := dao.TableQrcodes.Ctx(ctx).
			WhereLike(dao.TableQrcodes.Columns().TableNumber, "%"+req.TableNumber+"%").
			Scan(&tableQrcodes)
		if err != nil {
			return nil, 0, err
		}

		if len(tableQrcodes) > 0 {
			var tableQrcodeIds []int
			for _, qrcode := range tableQrcodes {
				tableQrcodeIds = append(tableQrcodeIds, qrcode.Id)
			}
			m = m.WhereIn(dao.Orders.Columns().TableQrcodeId, tableQrcodeIds)
		} else {
			// 如果没有找到匹配的桌号，则返回空结果
			return []adminv1.AdminOrder{}, 0, nil
		}
	}

	// 1.5 用户ID筛选
	if req.UserID > 0 {
		m = m.Where(dao.Orders.Columns().UserId, req.UserID)
	}

	// 2. 查询总数
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	// 3. 分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// 4. 查询数据
	var orders []*entity.Orders
	err = m.Page(page, limit).
		Order(dao.Orders.Columns().Id + " DESC").
		Scan(&orders)
	if err != nil {
		return nil, 0, err
	}

	// 5. 转换为API响应格式
	result := make([]adminv1.AdminOrder, len(orders))
	for i, order := range orders {
		// 查询桌号信息
		tableNumber := ""
		if order.TableQrcodeId > 0 {
			var tableQrcode *entity.TableQrcodes
			err = dao.TableQrcodes.Ctx(ctx).
				Where(dao.TableQrcodes.Columns().Id, order.TableQrcodeId).
				Scan(&tableQrcode)
			if err == nil && tableQrcode != nil {
				tableNumber = tableQrcode.TableNumber
			}
		}

		// 查询用户信息
		var user *entity.Users
		userName := ""
		userNickname := ""
		userPhone := ""
		if order.UserId > 0 {
			err = dao.Users.Ctx(ctx).
				Where(dao.Users.Columns().Id, order.UserId).
				Scan(&user)
			if err == nil && user != nil {
				userName = user.Openid // 使用openid作为用户名
				userNickname = user.Nickname
				userPhone = user.Phone
			}
		}

		// 查询订单项
		var orderItems []*entity.OrderItems
		err = dao.OrderItems.Ctx(ctx).
			Where(dao.OrderItems.Columns().OrderId, order.Id).
			Scan(&orderItems)

		// 计算商品总数量
		itemCount := 0

		// 转换订单项
		items := make([]orderv1.OrderItem, 0)
		if err == nil {
			for _, item := range orderItems {
				// 查询商品图片
				var product *entity.Products
				imageURL := ""
				if item.ProductId > 0 {
					err = dao.Products.Ctx(ctx).
						Where(dao.Products.Columns().Id, item.ProductId).
						Fields("image_url").
						Scan(&product)
					if err == nil && product != nil {
						imageURL = product.ImageUrl
					}
				}

				orderItem := orderv1.OrderItem{
					ProductID:     int64(item.ProductId),
					Name:          item.ProductName,
					ImageURL:      imageURL,
					Quantity:      item.Quantity,
					ItemPrice:     item.ItemPrice,
					IsPackageItem: item.IsPackageItem == 1,
					UserPackageID: int64(item.UserPackageId),
				}

				// 累加商品数量
				itemCount += item.Quantity

				// 如果是套餐商品，查询套餐详情
				if item.IsPackageItem == 1 && item.UserPackageId > 0 {
					// 查询用户套餐详情
					packageDetail, err := UserPackage().GetUserPackageDetail(ctx, int64(item.UserPackageId))
					if err == nil && packageDetail != nil {
						// 添加套餐状态和有效期信息
						orderItem.PackageInfo = &orderv1.PackageInfo{
							Status:    packageDetail.Status,
							StartTime: packageDetail.StartTime,
							EndTime:   packageDetail.EndTime,
						}
					}
				}

				items = append(items, orderItem)
			}
		}

		// 如果订单项为空，但是订单是套餐订单，则查询关联的套餐信息
		if len(items) == 0 {
			// 查询是否有关联的套餐
			var userPackage entity.UserPackages
			err = g.DB().Model("user_packages").
				Where("order_id", order.Id).
				Scan(&userPackage)

			if err == nil && userPackage.Id > 0 {
				// 查询套餐详情
				packageInfo, err := UserPackage().GetUserPackageDetail(ctx, int64(userPackage.Id))
				if err == nil && packageInfo != nil {
					// 查询套餐商品信息
					var packageDetail entity.DrinkAllYouCanPackages
					err = g.DB().Model("drink_all_you_can_packages").
						Where("id", userPackage.PackageId).
						Scan(&packageDetail)

					if err == nil && packageDetail.Id > 0 {
						// 创建一个套餐商品项
						packageItem := orderv1.OrderItem{
							ProductID:     int64(packageDetail.Id),
							Name:          packageDetail.Name,
							ImageURL:      "", // 套餐暂无图片，设置为空字符串
							Quantity:      1,
							ItemPrice:     packageDetail.Price,
							IsPackageItem: true,
							UserPackageID: int64(userPackage.Id),
							Notes:         "",
							PackageInfo: &orderv1.PackageInfo{
								Status:    packageInfo.Status,
								StartTime: packageInfo.StartTime,
								EndTime:   packageInfo.EndTime,
							},
						}

						// 添加到订单商品列表
						items = append(items, packageItem)

						// 套餐也算一个商品
						itemCount = 1
					}
				}
			}
		}

		result[i] = adminv1.AdminOrder{
			ID:            int64(order.Id),
			OrderSN:       order.OrderSn,
			UserID:        int64(order.UserId),
			UserName:      userName,
			UserNickname:  userNickname,
			UserPhone:     userPhone,
			TableQrcodeID: int64(order.TableQrcodeId),
			TableNumber:   tableNumber,
			TotalAmount:   order.TotalAmount,
			ItemCount:     itemCount,
			PaymentStatus: order.PaymentStatus,
			OrderStatus:   order.OrderStatus,
			CreatedAt:     order.CreatedAt.String(),
			UpdatedAt:     order.UpdatedAt.String(),
			Items:         items,
		}

		// 添加支付时间（如果有）
		if order.PaidAt != nil {
			result[i].PaidAt = order.PaidAt.String()
		}
	}

	return result, total, nil
}

// GetOrderDetailAdmin 管理端获取订单详情
func (s *orderService) GetOrderDetailAdmin(ctx context.Context, orderId int64) (*adminv1.AdminOrder, error) {
	// 1. 查询订单基本信息
	var order *entity.Orders
	err := dao.Orders.Ctx(ctx).
		Where(dao.Orders.Columns().Id, orderId).
		Scan(&order)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.New("订单不存在")
	}

	// 2. 查询订单项
	var orderItems []*entity.OrderItems
	err = dao.OrderItems.Ctx(ctx).
		Where(dao.OrderItems.Columns().OrderId, order.Id).
		Scan(&orderItems)
	if err != nil {
		return nil, err
	}

	// 3. 查询桌号信息
	var tableQrcode *entity.TableQrcodes
	tableNumber := ""
	if order.TableQrcodeId > 0 {
		err = dao.TableQrcodes.Ctx(ctx).
			Where(dao.TableQrcodes.Columns().Id, order.TableQrcodeId).
			Scan(&tableQrcode)
		if err != nil {
			return nil, err
		}
		if tableQrcode != nil {
			tableNumber = tableQrcode.TableNumber
		}
	}

	// 3.1 查询用户信息
	var user *entity.Users
	userName := ""
	userNickname := ""
	userPhone := ""
	if order.UserId > 0 {
		err = dao.Users.Ctx(ctx).
			Where(dao.Users.Columns().Id, order.UserId).
			Scan(&user)
		if err == nil && user != nil {
			userName = user.Openid // 使用openid作为用户名
			userNickname = user.Nickname
			userPhone = user.Phone
		}
	}

	// 4. 转换订单项
	items := make([]orderv1.OrderItem, len(orderItems))

	// 计算商品总数量
	itemCount := 0

	for i, item := range orderItems {
		// 查询商品图片
		var product *entity.Products
		imageURL := ""
		if item.ProductId > 0 {
			err = dao.Products.Ctx(ctx).
				Where(dao.Products.Columns().Id, item.ProductId).
				Fields("image_url").
				Scan(&product)
			if err == nil && product != nil {
				imageURL = product.ImageUrl
			}
		}

		items[i] = orderv1.OrderItem{
			ProductID:     int64(item.ProductId),
			Name:          item.ProductName,
			ImageURL:      imageURL,
			Quantity:      item.Quantity,
			ItemPrice:     item.ItemPrice,
			IsPackageItem: item.IsPackageItem == 1,
			UserPackageID: int64(item.UserPackageId),
			Notes:         "", // 商品备注字段，如果需要可以添加到数据库表中
		}

		// 累加商品数量
		itemCount += item.Quantity

		// 如果是套餐商品，查询套餐详情
		if item.IsPackageItem == 1 && item.UserPackageId > 0 {
			// 查询用户套餐详情
			packageDetail, err := UserPackage().GetUserPackageDetail(ctx, int64(item.UserPackageId))
			if err == nil && packageDetail != nil {
				// 添加套餐状态和有效期信息
				items[i].PackageInfo = &orderv1.PackageInfo{
					Status:    packageDetail.Status,
					StartTime: packageDetail.StartTime,
					EndTime:   packageDetail.EndTime,
				}
			}
		}
	}

	// 5. 如果订单项为空，但是订单是套餐订单，则查询关联的套餐信息
	if len(orderItems) == 0 {
		// 查询是否有关联的套餐
		var userPackage entity.UserPackages
		err = g.DB().Model("user_packages").
			Where("order_id", orderId).
			Scan(&userPackage)

		if err == nil && userPackage.Id > 0 {
			// 查询套餐详情
			packageInfo, err := UserPackage().GetUserPackageDetail(ctx, int64(userPackage.Id))
			if err == nil && packageInfo != nil {
				// 查询套餐商品信息
				var packageDetail entity.DrinkAllYouCanPackages
				err = g.DB().Model("drink_all_you_can_packages").
					Where("id", userPackage.PackageId).
					Scan(&packageDetail)

				if err == nil && packageDetail.Id > 0 {
					// 创建一个套餐商品项
					packageItem := orderv1.OrderItem{
						ProductID:     int64(packageDetail.Id),
						Name:          packageDetail.Name,
						ImageURL:      "", // 套餐暂无图片，设置为空字符串
						Quantity:      1,
						ItemPrice:     packageDetail.Price,
						IsPackageItem: true,
						UserPackageID: int64(userPackage.Id),
						Notes:         "",
						PackageInfo: &orderv1.PackageInfo{
							Status:    packageInfo.Status,
							StartTime: packageInfo.StartTime,
							EndTime:   packageInfo.EndTime,
						},
					}

					// 添加到订单商品列表
					items = append(items, packageItem)

					// 套餐也算一个商品
					itemCount = 1
				}
			}
		}
	}

	// 6. 转换为API响应格式
	result := &adminv1.AdminOrder{
		ID:            int64(order.Id),
		OrderSN:       order.OrderSn,
		UserID:        int64(order.UserId),
		UserName:      userName,
		UserNickname:  userNickname,
		UserPhone:     userPhone,
		TableQrcodeID: int64(order.TableQrcodeId),
		TableNumber:   tableNumber,
		TotalAmount:   order.TotalAmount,
		ItemCount:     itemCount,
		PaymentStatus: order.PaymentStatus,
		OrderStatus:   order.OrderStatus,
		CreatedAt:     order.CreatedAt.String(),
		UpdatedAt:     order.UpdatedAt.String(),
		Items:         items,
		TotalNotes:    "", // 订单备注字段，如果需要可以添加到数据库表中
	}

	// 添加支付时间（如果有）
	if order.PaidAt != nil {
		result.PaidAt = order.PaidAt.String()
	}

	return result, nil
}

// UpdateOrderStatus 更新订单状态
func (s *orderService) UpdateOrderStatus(ctx context.Context, orderId int64, status string, reason string) (*adminv1.AdminOrder, error) {
	// 1. 验证订单是否存在
	var order *entity.Orders
	err := dao.Orders.Ctx(ctx).
		Where(dao.Orders.Columns().Id, orderId).
		Scan(&order)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.New("订单不存在")
	}

	// 2. 解析状态
	statuses := utility.SplitOrderStatus(status)
	if statuses.PaymentStatus == "" && statuses.OrderStatus == "" {
		return nil, gerror.New("无效的订单状态")
	}

	// 3. 开启事务
	var result *adminv1.AdminOrder
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 3.1 构建更新数据
		data := g.Map{
			dao.Orders.Columns().UpdatedAt: gtime.Now(),
		}

		// 记录是否更新为已支付状态，用于后续处理套餐激活
		updateToPaid := false

		if statuses.PaymentStatus != "" {
			// 确保使用正确的枚举值
			switch statuses.PaymentStatus {
			case "pending_payment":
				data[dao.Orders.Columns().PaymentStatus] = "pending"
			case "paid", "cancelled":
				data[dao.Orders.Columns().PaymentStatus] = statuses.PaymentStatus
			default:
				return gerror.New("无效的支付状态")
			}

			// 如果是已支付状态，记录支付时间
			if statuses.PaymentStatus == "paid" && order.PaymentStatus != "paid" {
				data[dao.Orders.Columns().PaidAt] = gtime.Now()
				updateToPaid = true
			}
		}

		if statuses.OrderStatus != "" {
			// 确保使用正确的枚举值
			switch statuses.OrderStatus {
			case "new", "processing", "completed", "cancelled":
				data[dao.Orders.Columns().OrderStatus] = statuses.OrderStatus
			default:
				return gerror.New("无效的订单状态")
			}
		}

		// 3.2 更新订单状态
		_, err := dao.Orders.Ctx(ctx).TX(tx).
			Where(dao.Orders.Columns().Id, orderId).
			Update(data)
		if err != nil {
			return err
		}

		// 3.3 如果订单状态更新为已支付，并且是套餐订单，则激活套餐
		if updateToPaid {
			// 查询是否有关联的套餐
			var userPackage *entity.UserPackages
			err = g.DB().Model("user_packages").
				Where("order_id", orderId).
				Scan(&userPackage)

			if err == nil && userPackage != nil && userPackage.Id > 0 {
				// 激活套餐
				err = UserPackage().ActivateUserPackageAfterPayment(ctx, orderId)
				if err != nil {
					return err
				}
			}
		}

		// 3.4 获取更新后的订单详情
		result, err = s.GetOrderDetailAdmin(ctx, orderId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
