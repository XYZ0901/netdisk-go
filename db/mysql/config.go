package mysql

import "database/sql"

var (
	MysqlDb    *sql.DB
	MysqlDbErr error
)

const (
	USER_NAME = "root"
	PASS_WORD = "root"
	HOST      = "localhost"
	PORT      = "3306"
	DATABASE  = "netdisk"
	CHARSET   = "utf8"
)
