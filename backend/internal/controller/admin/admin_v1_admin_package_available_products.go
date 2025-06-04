package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminPackageAvailableProducts(ctx context.Context, req *v1.AdminPackageAvailableProductsReq) (res *v1.AdminPackageAvailableProductsRes, err error) {
	// 调用服务获取可添加到套餐的商品列表
	list, total, err := service.Package().GetAvailableProducts(ctx, req.PackageID, req.Keyword, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminPackageAvailableProductsRes{}
	res.Code = 200
	res.Message = "获取可添加到套餐的商品列表成功"
	res.Data.List = list
	res.Data.Total = total

	return res, nil
}
