package qrcode

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"backend/api/qrcode/v1"
)

func (c *ControllerV1) CreateTableQrcode(ctx context.Context, req *v1.CreateTableQrcodeReq) (res *v1.CreateTableQrcodeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
