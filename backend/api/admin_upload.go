package api

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// 管理员上传图片请求
type AdminUploadImageReq struct {
	g.Meta `path:"/upload/image" method:"post" summary:"上传图片(管理员)"`
	File   *ghttp.UploadFile `json:"file" type:"file" description:"图片文件" v:"required"`
}

// 管理员上传图片响应
type AdminUploadImageRes struct {
	Url string `json:"url" description:"图片URL"`
}
