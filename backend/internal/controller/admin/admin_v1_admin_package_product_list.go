package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminPackageProductList(ctx context.Context, req *v1.AdminPackageProductListReq) (res *v1.AdminPackageProductListRes, err error) {
	// 调用服务获取套餐包含的商品列表
	list, err := service.Package().GetPackageProducts(ctx, req.PackageID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminPackageProductListRes{}
	res.Code = 200
	res.Message = "获取套餐商品列表成功"
	res.Data.List = list

	return res, nil
}
