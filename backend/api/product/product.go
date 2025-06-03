// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package product

import (
	"context"

	"backend/api/product/v1"
)

type IProductV1 interface {
	UserCategoryList(ctx context.Context, req *v1.UserCategoryListReq) (res *v1.UserCategoryListRes, err error)
	CategoryList(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error)
	ProductList(ctx context.Context, req *v1.ProductListReq) (res *v1.ProductListRes, err error)
	ProductDetail(ctx context.Context, req *v1.ProductDetailReq) (res *v1.ProductDetailRes, err error)
	UserProductList(ctx context.Context, req *v1.UserProductListReq) (res *v1.UserProductListRes, err error)
	UserProductDetail(ctx context.Context, req *v1.UserProductDetailReq) (res *v1.UserProductDetailRes, err error)
}
