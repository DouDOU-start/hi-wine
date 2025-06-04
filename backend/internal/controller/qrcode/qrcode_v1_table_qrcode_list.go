package qrcode

import (
	"context"

	"backend/api/common"
	v1 "backend/api/qrcode/v1"
	"backend/internal/service"
)

func (c *ControllerV1) TableQrcodeList(ctx context.Context, req *v1.TableQrcodeListReq) (res *v1.TableQrcodeListRes, err error) {
	// 初始化响应对象
	res = &v1.TableQrcodeListRes{}

	// 调用二维码服务获取桌号二维码列表
	qrcodeService := service.Qrcode()
	list, total, err := qrcodeService.GetTableQrcodeList(ctx, req)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = err.Error()
		return res, nil
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "获取桌号二维码列表成功"
	res.Data.List = list
	res.Data.Total = total

	return res, nil
}
