// 此文件已废弃，管理端分类列表接口请使用 /api/v1/admin/categories 路由及相关实现。

package product

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "backend/api/product/v1"
)

func (c *ControllerV1) CategoryList(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
