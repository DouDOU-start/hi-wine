// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DrinkAllYouCanPackagesDao is the data access object for the table drink_all_you_can_packages.
type DrinkAllYouCanPackagesDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  DrinkAllYouCanPackagesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// DrinkAllYouCanPackagesColumns defines and stores column names for the table drink_all_you_can_packages.
type DrinkAllYouCanPackagesColumns struct {
	Id              string // 套餐ID
	Name            string // 套餐名称
	Description     string // 套餐描述
	Price           string // 套餐价格
	DurationMinutes string // 有效时长（分钟），0表示无时间限制
	IsActive        string // 是否激活（是否可售）
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// drinkAllYouCanPackagesColumns holds the columns for the table drink_all_you_can_packages.
var drinkAllYouCanPackagesColumns = DrinkAllYouCanPackagesColumns{
	Id:              "id",
	Name:            "name",
	Description:     "description",
	Price:           "price",
	DurationMinutes: "duration_minutes",
	IsActive:        "is_active",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewDrinkAllYouCanPackagesDao creates and returns a new DAO object for table data access.
func NewDrinkAllYouCanPackagesDao(handlers ...gdb.ModelHandler) *DrinkAllYouCanPackagesDao {
	return &DrinkAllYouCanPackagesDao{
		group:    "default",
		table:    "drink_all_you_can_packages",
		columns:  drinkAllYouCanPackagesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DrinkAllYouCanPackagesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DrinkAllYouCanPackagesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DrinkAllYouCanPackagesDao) Columns() DrinkAllYouCanPackagesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DrinkAllYouCanPackagesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DrinkAllYouCanPackagesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DrinkAllYouCanPackagesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
