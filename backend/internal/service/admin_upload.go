package service

import (
	"backend/api"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

type IAdminUploadService interface {
	Image(ctx context.Context, req *api.AdminUploadImageReq) (res *api.AdminUploadImageRes, err error)
}

var adminUploadService = localAdminUploadService{}

type localAdminUploadService struct{}

func AdminUpload() IAdminUploadService {
	return &adminUploadService
}

// 上传图片
func (s *localAdminUploadService) Image(ctx context.Context, req *api.AdminUploadImageReq) (res *api.AdminUploadImageRes, err error) {
	res = &api.AdminUploadImageRes{}

	// 检查文件类型
	fileType := req.File.Header.Get("Content-Type")
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
		return nil, gerror.New("只支持上传jpg、png、gif格式的图片")
	}

	// 检查文件大小
	if req.File.Size > 5*1024*1024 { // 5MB
		return nil, gerror.New("图片大小不能超过5MB")
	}

	// 打开文件
	fileReader, err := req.File.Open()
	if err != nil {
		return nil, gerror.Wrap(err, "打开文件失败")
	}
	defer fileReader.Close()

	// 使用MinIO服务上传文件
	objectName, err := Minio().Upload(ctx, fileReader, req.File.Size, req.File.Filename, fileType)
	if err != nil {
		return nil, gerror.Wrap(err, "上传文件失败")
	}

	// 获取文件URL
	fileUrl, err := Minio().GetURL(ctx, objectName)
	if err != nil {
		return nil, gerror.Wrap(err, "获取文件URL失败")
	}

	res.Url = fileUrl

	return res, nil
}
