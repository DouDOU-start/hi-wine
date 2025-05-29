package controller

import (
	"backend/api"
	"backend/internal/service"
	"context"
)

type AdminUploadController struct{}

func NewAdminUpload() *AdminUploadController {
	return &AdminUploadController{}
}

// 上传图片
func (c *AdminUploadController) Image(ctx context.Context, req *api.AdminUploadImageReq) (res *api.AdminUploadImageRes, err error) {
	return service.AdminUpload().Image(ctx, req)
}
