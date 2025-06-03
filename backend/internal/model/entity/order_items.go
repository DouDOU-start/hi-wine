// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderItems is the golang structure for table order_items.
type OrderItems struct {
	Id            int         `json:"id"            orm:"id"              description:"订单项ID"`                  // 订单项ID
	OrderId       int         `json:"orderId"       orm:"order_id"        description:"订单ID"`                   // 订单ID
	ProductId     int         `json:"productId"     orm:"product_id"      description:"商品ID"`                   // 商品ID
	ProductName   string      `json:"productName"   orm:"product_name"    description:"冗余商品名称"`                 // 冗余商品名称
	Price         float64     `json:"price"         orm:"price"           description:"下单时商品单价"`                // 下单时商品单价
	Quantity      int         `json:"quantity"      orm:"quantity"        description:"购买数量"`                   // 购买数量
	Subtotal      float64     `json:"subtotal"      orm:"subtotal"        description:"小计"`                     // 小计
	IsPackageItem int         `json:"isPackageItem" orm:"is_package_item" description:"是否为畅饮套餐内商品"`             // 是否为畅饮套餐内商品
	UserPackageId int         `json:"userPackageId" orm:"user_package_id" description:"关联的用户套餐购买记录ID（如果为套餐商品）"` // 关联的用户套餐购买记录ID（如果为套餐商品）
	ItemPrice     float64     `json:"itemPrice"     orm:"item_price"      description:"该订单项的实际结算价格（畅饮则为0）"`     // 该订单项的实际结算价格（畅饮则为0）
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:"创建时间"`                   // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      description:"更新时间"`                   // 更新时间
}
