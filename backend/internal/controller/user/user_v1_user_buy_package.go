package user

import (
	"context"

	v1 "backend/api/user/v1"
	"backend/internal/service"
	"backend/internal/utility"
)

func (c *ControllerV1) UserBuyPackage(ctx context.Context, req *v1.UserBuyPackageReq) (res *v1.UserBuyPackageRes, err error) {
	// 1. 从上下文中获取用户ID
	userId, err := utility.GetUserID(ctx)
	// 暂时注释掉用户ID检查逻辑，后续可以根据实际需求决定是否需要
	// if err != nil {
	// 	return nil, err
	// }

	// 2. 调用服务购买套餐
	orderId, err := service.UserPackageForUser().BuyPackage(ctx, userId, req.PackageID)
	if err != nil {
		return nil, err
	}

	// 3. 构建响应
	res = &v1.UserBuyPackageRes{}
	res.Code = 200
	res.Message = "套餐购买成功，请完成支付"
	res.Data.OrderID = orderId

	return res, nil
}
