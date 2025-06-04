package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/internal/service"
)

func (c *ControllerV1) AdminPackageDelete(ctx context.Context, req *v1.AdminPackageDeleteReq) (res *v1.AdminPackageDeleteRes, err error) {
	// 调用服务删除套餐
	err = service.Package().DeletePackage(ctx, req.PackageID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.AdminPackageDeleteRes{}
	res.Code = 200
	res.Message = "删除套餐成功"

	return res, nil
}
