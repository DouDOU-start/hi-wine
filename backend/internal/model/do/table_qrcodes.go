// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TableQrcodes is the golang structure of table table_qrcodes for DAO operations like Where/Data.
type TableQrcodes struct {
	g.Meta      `orm:"table:table_qrcodes, do:true"`
	Id          interface{} // ID
	TableNumber interface{} // 桌号，例如“A1”, “吧台”
	QrcodeUrl   interface{} // 生成的二维码图片URL
	Status      interface{} // 桌位状态
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
