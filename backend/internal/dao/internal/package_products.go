// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PackageProductsDao is the data access object for the table package_products.
type PackageProductsDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  PackageProductsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// PackageProductsColumns defines and stores column names for the table package_products.
type PackageProductsColumns struct {
	PackageId string // 关联畅饮套餐ID
	ProductId string // 关联商品ID
	CreatedAt string // 创建时间
}

// packageProductsColumns holds the columns for the table package_products.
var packageProductsColumns = PackageProductsColumns{
	PackageId: "package_id",
	ProductId: "product_id",
	CreatedAt: "created_at",
}

// NewPackageProductsDao creates and returns a new DAO object for table data access.
func NewPackageProductsDao(handlers ...gdb.ModelHandler) *PackageProductsDao {
	return &PackageProductsDao{
		group:    "default",
		table:    "package_products",
		columns:  packageProductsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PackageProductsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PackageProductsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PackageProductsDao) Columns() PackageProductsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PackageProductsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PackageProductsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PackageProductsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
