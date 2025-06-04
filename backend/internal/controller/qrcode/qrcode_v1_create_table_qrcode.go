package qrcode

import (
	"context"

	"backend/api/common"
	v1 "backend/api/qrcode/v1"
	"backend/internal/service"
)

func (c *ControllerV1) CreateTableQrcode(ctx context.Context, req *v1.CreateTableQrcodeReq) (res *v1.CreateTableQrcodeRes, err error) {
	// 初始化响应对象
	res = &v1.CreateTableQrcodeRes{}

	// 调用二维码服务生成桌号二维码
	qrcodeService := service.Qrcode()
	tableQrcode, err := qrcodeService.CreateTableQrcode(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "生成桌号二维码成功"
	res.Data = *tableQrcode

	return res, nil
}
