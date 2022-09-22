package model

import "gorm.io/gorm"

const TableSysAdminBindRole = "sys_admin_bind_role"

type SysAdminBindRole struct {
	Id        int `gorm:"primary_key"`
	AdminId   int
	RoleId    int
	DeletedAt gorm.DeletedAt

	Role SysRole `gorm:"References:RoleId;foreignKey:Id"`
}
