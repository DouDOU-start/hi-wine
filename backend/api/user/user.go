// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"backend/api/user/v1"
)

type IUserV1 interface {
	UserOrderList(ctx context.Context, req *v1.UserOrderListReq) (res *v1.UserOrderListRes, err error)
	UserOrderDetail(ctx context.Context, req *v1.UserOrderDetailReq) (res *v1.UserOrderDetailRes, err error)
	UserPackageList(ctx context.Context, req *v1.UserPackageListReq) (res *v1.UserPackageListRes, err error)
	UserPackageDetail(ctx context.Context, req *v1.UserPackageDetailReq) (res *v1.UserPackageDetailRes, err error)
	UserBuyPackage(ctx context.Context, req *v1.UserBuyPackageReq) (res *v1.UserBuyPackageRes, err error)
	WechatLogin(ctx context.Context, req *v1.WechatLoginReq) (res *v1.WechatLoginRes, err error)
	UserProfile(ctx context.Context, req *v1.UserProfileReq) (res *v1.UserProfileRes, err error)
	UpdateUserProfile(ctx context.Context, req *v1.UpdateUserProfileReq) (res *v1.UpdateUserProfileRes, err error)
}
