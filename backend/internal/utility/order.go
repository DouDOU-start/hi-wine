package utility

import (
	"strings"
	"time"

	"fmt"

	"github.com/gogf/gf/v2/util/grand"
)

// OrderStatus 订单状态结构体
type OrderStatus struct {
	PaymentStatus string // 支付状态
	OrderStatus   string // 订单状态
}

// SplitOrderStatus 拆分订单状态
// 格式：payment_status_order_status
// 例如：pending_new, paid_processing, cancelled_cancelled
func SplitOrderStatus(status string) OrderStatus {
	result := OrderStatus{}

	if status == "" {
		return result
	}

	// 拆分状态
	parts := strings.Split(status, "_")
	if len(parts) >= 1 {
		// 处理支付状态
		switch parts[0] {
		case "pending", "pending_payment":
			result.PaymentStatus = "pending"
		case "paid", "cancelled":
			result.PaymentStatus = parts[0]
		}
	}

	if len(parts) >= 2 {
		// 处理订单状态
		switch parts[1] {
		case "new", "processing", "completed", "cancelled":
			result.OrderStatus = parts[1]
		}
	}

	return result
}

// GenerateOrderSN 生成订单号
// 格式：年月日时分秒 + 4位随机数
func GenerateOrderSN() string {
	// 当前时间，格式：20060102150405
	timeStr := time.Now().Format("20060102150405")

	// 4位随机数
	randStr := grand.Digits(4)

	// 组合成订单号
	return fmt.Sprintf("%s%s", timeStr, randStr)
}
