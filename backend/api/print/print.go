// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package print

import (
	"context"

	"backend/api/print/v1"
)

type IPrintV1 interface {
	PrintOrder(ctx context.Context, req *v1.PrintOrderReq) (res *v1.PrintOrderRes, err error)
}
