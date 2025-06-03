// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductsDao is the data access object for the table products.
type ProductsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ProductsColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ProductsColumns defines and stores column names for the table products.
type ProductsColumns struct {
	Id          string // 商品ID
	CategoryId  string // 所属分类ID
	Name        string // 商品名称
	Description string // 商品描述
	Price       string // 商品价格
	ImageUrl    string // 商品图片URL
	Stock       string // 库存数量
	IsActive    string // 是否上架
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
}

// productsColumns holds the columns for the table products.
var productsColumns = ProductsColumns{
	Id:          "id",
	CategoryId:  "category_id",
	Name:        "name",
	Description: "description",
	Price:       "price",
	ImageUrl:    "image_url",
	Stock:       "stock",
	IsActive:    "is_active",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewProductsDao creates and returns a new DAO object for table data access.
func NewProductsDao(handlers ...gdb.ModelHandler) *ProductsDao {
	return &ProductsDao{
		group:    "default",
		table:    "products",
		columns:  productsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ProductsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ProductsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ProductsDao) Columns() ProductsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ProductsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ProductsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ProductsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
