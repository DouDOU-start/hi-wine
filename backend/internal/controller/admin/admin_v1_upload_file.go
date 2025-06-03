package admin

import (
	"context"
	"path"
	"strings"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/utility/minio"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) UploadFile(ctx context.Context, req *v1.UploadFileReq) (res *v1.UploadFileRes, err error) {
	// 创建响应对象
	res = &v1.UploadFileRes{}

	// 获取上传的文件
	file := req.File
	if file == nil {
		res.Code = common.CodeParamError
		res.Message = "请选择文件"
		return res, nil
	}

	// 检查文件类型
	fileType := req.Type
	if fileType == "" {
		fileType = "common"
	}

	g.Log().Info(ctx, "开始上传文件: ", file.Filename, ", 类型: ", fileType, ", 大小: ", file.Size)

	// 检查文件扩展名
	ext := strings.ToLower(path.Ext(file.Filename))
	allowExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	isValidExt := false
	for _, allowExt := range allowExts {
		if ext == allowExt {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		res.Code = common.CodeParamError
		res.Message = "不支持的文件类型，仅支持jpg、jpeg、png、gif、webp格式图片"
		return res, nil
	}

	// 检查文件大小（限制为5MB）
	maxSize := int64(5 * 1024 * 1024)
	if file.Size > maxSize {
		res.Code = common.CodeParamError
		res.Message = "文件大小不能超过5MB"
		return res, nil
	}

	// 获取MinIO客户端
	minioClient := minio.GetClient()
	if minioClient == nil {
		res.Code = common.CodeServerError
		res.Message = "文件存储服务未初始化"
		return res, nil
	}

	// 上传文件到MinIO
	fileURL, err := minioClient.UploadFile(ctx, file.FileHeader, fileType)
	if err != nil {
		g.Log().Error(ctx, "文件上传失败: ", err)
		res.Code = common.CodeServerError
		res.Message = "文件上传失败: " + err.Error()
		return res, nil
	}

	g.Log().Info(ctx, "文件上传成功，URL: ", fileURL)

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "文件上传成功"
	res.Data.URL = fileURL

	return res, nil
}
