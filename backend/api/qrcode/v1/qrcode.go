package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 二维码分组

type TableQrcode struct {
	ID          int64  `json:"id"`
	TableNumber string `json:"table_number"`
	QrcodeURL   string `json:"qrcode_url"`
	CreatedAt   string `json:"created_at"`
}

// 生成桌号二维码
type CreateTableQrcodeReq struct {
	g.Meta      `path:"/table-qrcodes" method:"post" tags:"管理端-二维码" summary:"生成桌号二维码"`
	TableNumber string `json:"table_number" v:"required#桌号必填"`
}
type CreateTableQrcodeRes struct {
	common.Response[TableQrcode] `json:",inline"`
}

// 获取桌号二维码列表
type TableQrcodeListReq struct {
	g.Meta      `path:"/table-qrcodes" method:"get" tags:"管理端-二维码" summary:"获取所有桌号二维码列表"`
	TableNumber string `json:"table_number" in:"query" description:"桌号模糊搜索"`
	Page        int    `json:"page" in:"query" description:"页码，默认1"`
	Limit       int    `json:"limit" in:"query" description:"每页数量，默认10"`
}
type TableQrcodeListRes struct {
	common.Response[struct {
		List  []TableQrcode `json:"list"`
		Total int           `json:"total"`
	}] `json:",inline"`
}

// 删除桌号二维码
type DeleteTableQrcodeReq struct {
	g.Meta `path:"/table-qrcodes/{id}" method:"delete" tags:"管理端-二维码" summary:"删除桌号二维码"`
	ID     int64 `json:"id" in:"path" v:"required#桌号ID必填"`
}
type DeleteTableQrcodeRes struct {
	common.Response[struct{}] `json:",inline"`
}
