package admin

import (
	"context"

	v1 "backend/api/admin/v1"
	"backend/api/common"
	"backend/internal/utility/minio"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) MinioTest(ctx context.Context, req *v1.MinioTestReq) (res *v1.MinioTestRes, err error) {
	// 创建响应对象
	res = &v1.MinioTestRes{}

	// 获取MinIO客户端
	minioClient := minio.GetClient()
	if minioClient == nil {
		res.Code = common.CodeServerError
		res.Message = "MinIO客户端未初始化"
		return res, nil
	}

	// 获取MinIO配置
	endpoint := g.Cfg().MustGet(ctx, "minio.endpoint").String()
	bucketName := g.Cfg().MustGet(ctx, "minio.bucketName").String()
	domain := g.Cfg().MustGet(ctx, "minio.domain").String()

	// 检查存储桶是否存在
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		res.Code = common.CodeServerError
		res.Message = "检查存储桶失败: " + err.Error()
		return res, nil
	}

	// 获取存储桶策略
	policy := ""
	if exists {
		policy, err = minioClient.GetBucketPolicy(ctx, bucketName)
		if err != nil {
			policy = "获取策略失败: " + err.Error()
		}
	}

	// 设置响应数据
	res.Code = common.CodeSuccess
	res.Message = "MinIO连接测试成功"
	res.Data.Status = "连接成功"
	res.Data.BucketExists = exists
	res.Data.BucketPolicy = policy
	res.Data.Endpoint = endpoint
	res.Data.BucketName = bucketName
	res.Data.Domain = domain

	return res, nil
}
