// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TableQrcodes is the golang structure for table table_qrcodes.
type TableQrcodes struct {
	Id          int         `json:"id"          orm:"id"           description:"ID"`              // ID
	TableNumber string      `json:"tableNumber" orm:"table_number" description:"桌号，例如“A1”, “吧台”"` // 桌号，例如“A1”, “吧台”
	QrcodeUrl   string      `json:"qrcodeUrl"   orm:"qrcode_url"   description:"生成的二维码图片URL"`     // 生成的二维码图片URL
	Status      string      `json:"status"      orm:"status"       description:"桌位状态"`            // 桌位状态
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"`            // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:"更新时间"`            // 更新时间
}
