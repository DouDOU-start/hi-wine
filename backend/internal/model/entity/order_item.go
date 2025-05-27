// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// OrderItem is the golang structure for table order_item.
type OrderItem struct {
	Id          int64   `json:"id"          orm:"id"           description:"主键"`   // 主键
	OrderId     int64   `json:"orderId"     orm:"order_id"     description:"订单ID"` // 订单ID
	ProductId   int64   `json:"productId"   orm:"product_id"   description:"商品ID"` // 商品ID
	ProductName string  `json:"productName" orm:"product_name" description:"商品名称"` // 商品名称
	Price       float64 `json:"price"       orm:"price"        description:"商品单价"` // 商品单价
	Quantity    int     `json:"quantity"    orm:"quantity"     description:"数量"`   // 数量
}
