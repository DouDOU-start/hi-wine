// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserPackages is the golang structure of table user_packages for DAO operations like Where/Data.
type UserPackages struct {
	g.Meta    `orm:"table:user_packages, do:true"`
	Id        interface{} // 记录ID
	UserId    interface{} // 关联用户ID
	PackageId interface{} // 关联畅饮套餐ID
	OrderId   interface{} // 关联购买此套餐的订单ID
	StartTime *gtime.Time // 套餐开始时间（首次使用时激活）
	EndTime   *gtime.Time // 套餐结束时间（根据duration_minutes计算）
	Status    interface{} // 套餐状态
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
