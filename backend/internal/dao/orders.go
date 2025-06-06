// =================================================================================
// This file is auto-generated by the GoFrame CLI tool. You may modify it as needed.
// =================================================================================

package dao

import (
	"backend/internal/dao/internal"
)

// ordersDao is the data access object for the table orders.
// You can define custom methods on it to extend its functionality as needed.
type ordersDao struct {
	*internal.OrdersDao
}

var (
	// Orders is a globally accessible object for table orders operations.
	Orders = ordersDao{internal.NewOrdersDao()}
)

// Add your custom methods and functionality below.
