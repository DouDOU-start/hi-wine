package product

import (
	"context"

	"backend/api/common"
	v1 "backend/api/product/v1"
	"backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) UserProductDetail(ctx context.Context, req *v1.UserProductDetailReq) (res *v1.UserProductDetailRes, err error) {
	// 创建响应对象
	res = &v1.UserProductDetailRes{
		Response: common.Response[v1.UserProductDetail]{
			Code:    common.CodeSuccess,
			Message: "获取商品详情成功",
		},
	}

	// 调用商品服务获取商品详情
	productService := service.Product()
	productDetail, err := productService.GetProductDetail(ctx, req.ProductID)
	if err != nil {
		g.Log().Error(ctx, "获取商品详情失败", err, map[string]interface{}{
			"product_id": req.ProductID,
		})
		res.Code = common.CodeServerError
		res.Message = "获取商品详情失败"
		return res, nil
	}

	// 设置响应数据
	res.Data = *productDetail

	return res, nil
}
