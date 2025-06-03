// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TableQrcodesDao is the data access object for the table table_qrcodes.
type TableQrcodesDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  TableQrcodesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// TableQrcodesColumns defines and stores column names for the table table_qrcodes.
type TableQrcodesColumns struct {
	Id          string // ID
	TableNumber string // 桌号，例如“A1”, “吧台”
	QrcodeUrl   string // 生成的二维码图片URL
	Status      string // 桌位状态
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
}

// tableQrcodesColumns holds the columns for the table table_qrcodes.
var tableQrcodesColumns = TableQrcodesColumns{
	Id:          "id",
	TableNumber: "table_number",
	QrcodeUrl:   "qrcode_url",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewTableQrcodesDao creates and returns a new DAO object for table data access.
func NewTableQrcodesDao(handlers ...gdb.ModelHandler) *TableQrcodesDao {
	return &TableQrcodesDao{
		group:    "default",
		table:    "table_qrcodes",
		columns:  tableQrcodesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TableQrcodesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TableQrcodesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TableQrcodesDao) Columns() TableQrcodesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TableQrcodesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TableQrcodesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TableQrcodesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
