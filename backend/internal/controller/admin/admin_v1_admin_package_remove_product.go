package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminPackageRemoveProduct(ctx context.Context, req *v1.AdminPackageRemoveProductReq) (res *v1.AdminPackageRemoveProductRes, err error) {
	// 调用服务从套餐中移除商品
	err = service.Package().RemovePackageProduct(ctx, req.PackageID, req.ProductID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminPackageRemoveProductRes{}
	res.Code = 200
	res.Message = "移除套餐商品成功"

	return res, nil
}
