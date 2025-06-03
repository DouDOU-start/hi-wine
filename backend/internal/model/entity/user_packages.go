// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserPackages is the golang structure for table user_packages.
type UserPackages struct {
	Id        int         `json:"id"        orm:"id"         description:"记录ID"`                         // 记录ID
	UserId    int         `json:"userId"    orm:"user_id"    description:"关联用户ID"`                       // 关联用户ID
	PackageId int         `json:"packageId" orm:"package_id" description:"关联畅饮套餐ID"`                     // 关联畅饮套餐ID
	OrderId   int         `json:"orderId"   orm:"order_id"   description:"关联购买此套餐的订单ID"`                 // 关联购买此套餐的订单ID
	StartTime *gtime.Time `json:"startTime" orm:"start_time" description:"套餐开始时间"`                       // 套餐开始时间
	EndTime   *gtime.Time `json:"endTime"   orm:"end_time"   description:"套餐结束时间（根据duration_minutes计算）"` // 套餐结束时间（根据duration_minutes计算）
	Status    string      `json:"status"    orm:"status"     description:"套餐状态"`                         // 套餐状态
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`                         // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`                         // 更新时间
}
