// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserPackagesDao is the data access object for the table user_packages.
type UserPackagesDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  UserPackagesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// UserPackagesColumns defines and stores column names for the table user_packages.
type UserPackagesColumns struct {
	Id        string // 记录ID
	UserId    string // 关联用户ID
	PackageId string // 关联畅饮套餐ID
	OrderId   string // 关联购买此套餐的订单ID
	StartTime string // 套餐开始时间
	EndTime   string // 套餐结束时间（根据duration_minutes计算）
	Status    string // 套餐状态
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// userPackagesColumns holds the columns for the table user_packages.
var userPackagesColumns = UserPackagesColumns{
	Id:        "id",
	UserId:    "user_id",
	PackageId: "package_id",
	OrderId:   "order_id",
	StartTime: "start_time",
	EndTime:   "end_time",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewUserPackagesDao creates and returns a new DAO object for table data access.
func NewUserPackagesDao(handlers ...gdb.ModelHandler) *UserPackagesDao {
	return &UserPackagesDao{
		group:    "default",
		table:    "user_packages",
		columns:  userPackagesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserPackagesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserPackagesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserPackagesDao) Columns() UserPackagesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserPackagesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserPackagesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserPackagesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
