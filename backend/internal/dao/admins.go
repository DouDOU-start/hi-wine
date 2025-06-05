// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package dao

import (
	"backend/internal/dao/internal"
)

// adminsDao is the manager for logic model data accessing and custom defined data operations functions.
type adminsDao struct {
	*internal.AdminsDao
}

var (
	// Admins is the global data access object.
	Admins = adminsDao{
		internal.NewAdminsDao(),
	}
) 