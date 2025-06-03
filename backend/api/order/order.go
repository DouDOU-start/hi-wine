// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package order

import (
	"context"

	"backend/api/order/v1"
)

type IOrderV1 interface {
	CreateOrder(ctx context.Context, req *v1.CreateOrderReq) (res *v1.CreateOrderRes, err error)
	OrderDetail(ctx context.Context, req *v1.OrderDetailReq) (res *v1.OrderDetailRes, err error)
}
