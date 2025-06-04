package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminPackageBatchRemoveProducts(ctx context.Context, req *v1.AdminPackageBatchRemoveProductsReq) (res *v1.AdminPackageBatchRemoveProductsRes, err error) {
	// 调用服务批量从套餐中移除商品
	err = service.Package().BatchRemovePackageProducts(ctx, req.PackageID, req.ProductIDs)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminPackageBatchRemoveProductsRes{}
	res.Code = 200
	res.Message = "批量移除套餐商品成功"

	return res, nil
}
