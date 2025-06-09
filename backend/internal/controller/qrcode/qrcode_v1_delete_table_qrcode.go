package qrcode

import (
	"context"

	"backend/api/common"
	v1 "backend/api/qrcode/v1"
	"backend/internal/service"
)

func (c *ControllerV1) DeleteTableQrcode(ctx context.Context, req *v1.DeleteTableQrcodeReq) (res *v1.DeleteTableQrcodeRes, err error) {
	// 初始化响应对象
	res = &v1.DeleteTableQrcodeRes{}

	// 调用二维码服务删除桌号二维码
	qrcodeService := service.Qrcode()
	err = qrcodeService.DeleteTableQrcode(ctx, req.ID)
	if err != nil {
		res.Response = common.Response[struct{}]{
			Code:    common.CodeServerError,
			Message: err.Error(),
			Data:    struct{}{},
		}
		return res, nil
	}

	// 设置成功响应
	res.Response = common.Response[struct{}]{
		Code:    common.CodeSuccess,
		Message: "删除桌号二维码成功",
		Data:    struct{}{},
	}

	return res, nil
}
