package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 文件代理接口
type FileProxyReq struct {
	g.Meta   `path:"/file/*file_path" method:"get" tags:"文件服务" summary:"文件代理访问"`
	FilePath string `json:"file_path" in:"path" v:"required#文件路径不能为空" description:"文件路径，例如：common/xxx.jpg"`
}

// 该接口直接返回文件内容，不使用标准响应格式
