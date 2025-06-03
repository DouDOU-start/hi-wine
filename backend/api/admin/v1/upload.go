package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// 上传文件
type UploadFileReq struct {
	g.Meta `path:"/api/v1/admin/upload" method:"post" mime:"multipart/form-data" tags:"管理端-文件上传" summary:"上传文件"`
	File   *ghttp.UploadFile `json:"file" type:"file" v:"required#请选择文件"`
	Type   string            `json:"type" in:"query" d:"common" description:"文件类型，如: product, avatar, common"`
}

type UploadFileRes struct {
	common.Response[struct {
		URL string `json:"url" description:"文件URL"`
	}] `json:",inline"`
}
