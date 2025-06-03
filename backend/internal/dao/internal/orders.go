// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrdersDao is the data access object for the table orders.
type OrdersDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  OrdersColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// OrdersColumns defines and stores column names for the table orders.
type OrdersColumns struct {
	Id            string // 订单ID
	OrderSn       string // 订单号，唯一
	UserId        string // 用户ID
	TableQrcodeId string // 关联的桌号二维码ID
	TotalAmount   string // 订单总金额
	PaymentStatus string // 支付状态
	OrderStatus   string // 订单状态
	PaymentMethod string // 支付方式（例如：wechat_pay）
	TransactionId string // 微信支付交易ID
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
	PaidAt        string // 支付时间
}

// ordersColumns holds the columns for the table orders.
var ordersColumns = OrdersColumns{
	Id:            "id",
	OrderSn:       "order_sn",
	UserId:        "user_id",
	TableQrcodeId: "table_qrcode_id",
	TotalAmount:   "total_amount",
	PaymentStatus: "payment_status",
	OrderStatus:   "order_status",
	PaymentMethod: "payment_method",
	TransactionId: "transaction_id",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	PaidAt:        "paid_at",
}

// NewOrdersDao creates and returns a new DAO object for table data access.
func NewOrdersDao(handlers ...gdb.ModelHandler) *OrdersDao {
	return &OrdersDao{
		group:    "default",
		table:    "orders",
		columns:  ordersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OrdersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OrdersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OrdersDao) Columns() OrdersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OrdersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OrdersDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *OrdersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
