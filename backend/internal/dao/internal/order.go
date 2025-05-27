// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrderDao is the data access object for the table order.
type OrderDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  OrderColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// OrderColumns defines and stores column names for the table order.
type OrderColumns struct {
	Id          string // 主键
	UserId      string // 下单用户ID
	TableId     string // 桌台ID
	TotalAmount string // 总金额
	Status      string // 订单状态（0待支付，1已支付，2已完成，3已取消）
	PayTime     string // 支付时间
	CreateTime  string // 创建时间
	UpdateTime  string // 更新时间
}

// orderColumns holds the columns for the table order.
var orderColumns = OrderColumns{
	Id:          "id",
	UserId:      "user_id",
	TableId:     "table_id",
	TotalAmount: "total_amount",
	Status:      "status",
	PayTime:     "pay_time",
	CreateTime:  "create_time",
	UpdateTime:  "update_time",
}

// NewOrderDao creates and returns a new DAO object for table data access.
func NewOrderDao(handlers ...gdb.ModelHandler) *OrderDao {
	return &OrderDao{
		group:    "default",
		table:    "order",
		columns:  orderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OrderDao) Columns() OrderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OrderDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
