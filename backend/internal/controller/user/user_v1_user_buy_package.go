package user

import (
	"context"

	v1 "backend/api/user/v1"
	"backend/internal/service"
	"backend/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) UserBuyPackage(ctx context.Context, req *v1.UserBuyPackageReq) (res *v1.UserBuyPackageRes, err error) {
	// 1. 从上下文中获取用户ID
	userId, err := utility.GetUserID(ctx)
	if err != nil {
		g.Log().Error(ctx, "购买套餐失败: 未找到用户ID", err)
		return nil, gerror.New("未登录或登录已过期")
	}

	g.Log().Debug(ctx, "购买套餐，用户ID:", userId, "套餐ID:", req.PackageID)

	// 2. 调用服务购买套餐
	orderId, err := service.UserPackageForUser().BuyPackage(ctx, userId, req.PackageID)
	if err != nil {
		g.Log().Error(ctx, "购买套餐失败:", err, "用户ID:", userId, "套餐ID:", req.PackageID)
		return nil, err
	}

	// 3. 构建响应
	res = &v1.UserBuyPackageRes{}
	res.Code = 200
	res.Message = "套餐购买成功，请完成支付"
	res.Data.OrderID = orderId

	g.Log().Info(ctx, "购买套餐成功，用户ID:", userId, "套餐ID:", req.PackageID, "订单ID:", orderId)

	return res, nil
}
