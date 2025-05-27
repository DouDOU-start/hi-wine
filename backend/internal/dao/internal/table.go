// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TableDao is the data access object for the table table.
type TableDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TableColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TableColumns defines and stores column names for the table table.
type TableColumns struct {
	Id         string // 主键
	Name       string // 桌台名称/编号
	Status     string // 状态（0空闲，1占用）
	CreateTime string // 创建时间
	UpdateTime string // 更新时间
}

// tableColumns holds the columns for the table table.
var tableColumns = TableColumns{
	Id:         "id",
	Name:       "name",
	Status:     "status",
	CreateTime: "create_time",
	UpdateTime: "update_time",
}

// NewTableDao creates and returns a new DAO object for table data access.
func NewTableDao(handlers ...gdb.ModelHandler) *TableDao {
	return &TableDao{
		group:    "default",
		table:    "table",
		columns:  tableColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TableDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TableDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TableDao) Columns() TableColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TableDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TableDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TableDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
