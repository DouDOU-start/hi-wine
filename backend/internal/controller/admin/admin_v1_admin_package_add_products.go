package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminPackageAddProducts(ctx context.Context, req *v1.AdminPackageAddProductsReq) (res *v1.AdminPackageAddProductsRes, err error) {
	// 调用服务为套餐添加商品
	err = service.Package().AddPackageProducts(ctx, req.PackageID, req.ProductIDs)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminPackageAddProductsRes{}
	res.Code = 200
	res.Message = "添加套餐商品成功"

	return res, nil
}
