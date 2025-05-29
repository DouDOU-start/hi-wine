package controller

import (
	"backend/api"
	"backend/internal/service"
	"context"
)

type AdminStatsController struct{}

func NewAdminStats() *AdminStatsController {
	return &AdminStatsController{}
}

// 获取仪表盘统计数据
func (c *AdminStatsController) Dashboard(ctx context.Context, req *api.AdminDashboardStatsReq) (res *api.AdminDashboardStatsRes, err error) {
	return service.AdminStats().Dashboard(ctx, req)
}

// 获取销售统计数据
func (c *AdminStatsController) Sales(ctx context.Context, req *api.AdminSalesStatsReq) (res *api.AdminSalesStatsRes, err error) {
	return service.AdminStats().Sales(ctx, req)
}

// 获取商品销量排行
func (c *AdminStatsController) ProductRanking(ctx context.Context, req *api.AdminProductRankingReq) (res *api.AdminProductRankingRes, err error) {
	return service.AdminStats().ProductRanking(ctx, req)
}

// 获取用户消费排行
func (c *AdminStatsController) UserRanking(ctx context.Context, req *api.AdminUserRankingReq) (res *api.AdminUserRankingRes, err error) {
	return service.AdminStats().UserRanking(ctx, req)
}

// 获取分类销售统计
func (c *AdminStatsController) CategorySales(ctx context.Context, req *api.AdminCategorySalesReq) (res *api.AdminCategorySalesRes, err error) {
	return service.AdminStats().CategorySales(ctx, req)
}
