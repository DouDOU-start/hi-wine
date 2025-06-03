// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package qrcode

import (
	"context"

	"backend/api/qrcode/v1"
)

type IQrcodeV1 interface {
	CreateTableQrcode(ctx context.Context, req *v1.CreateTableQrcodeReq) (res *v1.CreateTableQrcodeRes, err error)
	TableQrcodeList(ctx context.Context, req *v1.TableQrcodeListReq) (res *v1.TableQrcodeListRes, err error)
}
