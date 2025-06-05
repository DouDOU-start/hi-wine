package v1

import (
	"backend/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// 打印分组

// 打印订单
type PrintOrderReq struct {
	g.Meta    `path:"/print/order" method:"post" tags:"管理端-打印" summary:"打印订单"`
	OrderID   int64  `json:"order_id" v:"required#订单ID必填"`
	PrintType string `json:"print_type,omitempty" description:"打印类型"`
	Notes     string `json:"notes,omitempty" description:"打印备注"`
}

type PrintOrderRes struct {
	common.Response[struct {
		Status     string `json:"status"`
		PrintJobID string `json:"print_job_id,omitempty"`
	}] `json:",inline"`
}

// 打印任务结果
type PrintOrderResult struct {
	Status     string `json:"status"`
	PrintJobID string `json:"print_job_id,omitempty"`
}
