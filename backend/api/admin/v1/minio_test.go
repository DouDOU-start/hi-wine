package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// MinIO测试接口
type MinioTestReq struct {
	g.Meta `path:"/api/v1/admin/minio-test" method:"get" tags:"管理端-文件上传" summary:"测试MinIO连接"`
}

type MinioTestRes struct {
	common.Response[struct {
		Status       string `json:"status" description:"连接状态"`
		BucketExists bool   `json:"bucket_exists" description:"存储桶是否存在"`
		BucketPolicy string `json:"bucket_policy" description:"存储桶策略"`
		Endpoint     string `json:"endpoint" description:"MinIO终端节点"`
		BucketName   string `json:"bucket_name" description:"存储桶名称"`
		Domain       string `json:"domain" description:"域名"`
	}] `json:",inline"`
}
