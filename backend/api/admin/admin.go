// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package admin

import (
	"context"

	"backend/api/admin/v1"
)

type IAdminV1 interface {
	// 获取当前管理员信息
	AdminProfile(ctx context.Context, req *v1.AdminProfileReq) (res *v1.AdminProfileRes, err error)
	
	// 分类相关接口
	CategoryList(ctx context.Context, req *v1.AdminCategoryListReq) (res *v1.AdminCategoryListRes, err error)
	CreateCategory(ctx context.Context, req *v1.AdminCreateCategoryReq) (res *v1.AdminCreateCategoryRes, err error)
	UpdateCategory(ctx context.Context, req *v1.AdminUpdateCategoryReq) (res *v1.AdminUpdateCategoryRes, err error)
	DeleteCategory(ctx context.Context, req *v1.AdminDeleteCategoryReq) (res *v1.AdminDeleteCategoryRes, err error)
	CategoryDetail(ctx context.Context, req *v1.AdminCategoryDetailReq) (res *v1.AdminCategoryDetailRes, err error)
	
	// 订单相关接口
	AdminOrderList(ctx context.Context, req *v1.AdminOrderListReq) (res *v1.AdminOrderListRes, err error)
	AdminOrderDetail(ctx context.Context, req *v1.AdminOrderDetailReq) (res *v1.AdminOrderDetailRes, err error)
	AdminOrderUpdateStatus(ctx context.Context, req *v1.AdminOrderUpdateStatusReq) (res *v1.AdminOrderUpdateStatusRes, err error)
	
	// 套餐相关接口
	AdminPackageList(ctx context.Context, req *v1.AdminPackageListReq) (res *v1.AdminPackageListRes, err error)
	AdminPackageCreate(ctx context.Context, req *v1.AdminPackageCreateReq) (res *v1.AdminPackageCreateRes, err error)
	AdminPackageUpdate(ctx context.Context, req *v1.AdminPackageUpdateReq) (res *v1.AdminPackageUpdateRes, err error)
	AdminPackageDelete(ctx context.Context, req *v1.AdminPackageDeleteReq) (res *v1.AdminPackageDeleteRes, err error)
	AdminPackageAddProducts(ctx context.Context, req *v1.AdminPackageAddProductsReq) (res *v1.AdminPackageAddProductsRes, err error)
	AdminPackageRemoveProduct(ctx context.Context, req *v1.AdminPackageRemoveProductReq) (res *v1.AdminPackageRemoveProductRes, err error)
	AdminPackageProductList(ctx context.Context, req *v1.AdminPackageProductListReq) (res *v1.AdminPackageProductListRes, err error)
	AdminPackageStats(ctx context.Context, req *v1.AdminPackageStatsReq) (res *v1.AdminPackageStatsRes, err error)
	AdminPackageDetail(ctx context.Context, req *v1.AdminPackageDetailReq) (res *v1.AdminPackageDetailRes, err error)
	AdminPackageFullDetail(ctx context.Context, req *v1.AdminPackageFullDetailReq) (res *v1.AdminPackageFullDetailRes, err error)
	AdminPackageWithProducts(ctx context.Context, req *v1.AdminPackageWithProductsReq) (res *v1.AdminPackageWithProductsRes, err error)
	
	// 商品相关接口
	AdminProductList(ctx context.Context, req *v1.AdminProductListReq) (res *v1.AdminProductListRes, err error)
	AdminProductCreate(ctx context.Context, req *v1.AdminProductCreateReq) (res *v1.AdminProductCreateRes, err error)
	AdminProductUpdate(ctx context.Context, req *v1.AdminProductUpdateReq) (res *v1.AdminProductUpdateRes, err error)
	AdminProductDelete(ctx context.Context, req *v1.AdminProductDeleteReq) (res *v1.AdminProductDeleteRes, err error)
	AdminProductDetail(ctx context.Context, req *v1.AdminProductDetailReq) (res *v1.AdminProductDetailRes, err error)
	
	// 用户套餐相关接口
	AdminUserPackageList(ctx context.Context, req *v1.AdminUserPackageListReq) (res *v1.AdminUserPackageListRes, err error)
	AdminUserPackageDetail(ctx context.Context, req *v1.AdminUserPackageDetailReq) (res *v1.AdminUserPackageDetailRes, err error)
	AdminUserPackageCreate(ctx context.Context, req *v1.AdminUserPackageCreateReq) (res *v1.AdminUserPackageCreateRes, err error)
	AdminUserPackageUpdateStatus(ctx context.Context, req *v1.AdminUserPackageUpdateStatusReq) (res *v1.AdminUserPackageUpdateStatusRes, err error)
	AdminUserActivePackages(ctx context.Context, req *v1.AdminUserActivePackagesReq) (res *v1.AdminUserActivePackagesRes, err error)
	AdminUserPackageFullDetail(ctx context.Context, req *v1.AdminUserPackageFullDetailReq) (res *v1.AdminUserPackageFullDetailRes, err error)

	// 用户管理相关接口
	AdminUserList(ctx context.Context, req *v1.AdminUserListReq) (res *v1.AdminUserListRes, err error)
	AdminUserDetail(ctx context.Context, req *v1.AdminUserDetailReq) (res *v1.AdminUserDetailRes, err error)

	// 上传文件
	UploadFile(ctx context.Context, req *v1.UploadFileReq) (res *v1.UploadFileRes, err error)
}
