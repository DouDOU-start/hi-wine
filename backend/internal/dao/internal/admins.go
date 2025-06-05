// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AdminsDao is the data access object for the table admins.
type AdminsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  AdminsColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification
}

// AdminsColumns defines and stores column names for table admins.
type AdminsColumns struct {
	Id          string // 管理员ID
	Username    string // 用户名
	Password    string // 密码（哈希）
	Role        string // 角色
	IsActive    string // 是否激活
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	LastLoginAt string // 最后登录时间
}

// adminsDaoColumns defines and stores column names for table admins.
var adminsDaoColumns = AdminsColumns{
	Id:          "id",
	Username:    "username",
	Password:    "password",
	Role:        "role",
	IsActive:    "is_active",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	LastLoginAt: "last_login_at",
}

// NewAdminsDao creates and returns a new DAO object for table data access.
func NewAdminsDao() *AdminsDao {
	return &AdminsDao{
		group:   "default",
		table:   "admins",
		columns: adminsDaoColumns,
	}
}

// DB retrieves and returns the underlying database connection object.
func (dao *AdminsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Ctx creates and returns the Model for current DAO, which is automatically bound with context.
func (dao *AdminsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
func (dao *AdminsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

// Columns returns all column names of table admins.
func (dao *AdminsDao) Columns() AdminsColumns {
	return dao.columns
}

// Table returns the table name of current dao.
func (dao *AdminsDao) Table() string {
	return dao.table
}

// Group returns the database group name of current dao.
func (dao *AdminsDao) Group() string {
	return dao.group
} 