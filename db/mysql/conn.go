package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func init() {
	dnDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		USER_NAME, PASS_WORD, HOST, PORT, DATABASE, CHARSET)
	MysqlDb, MysqlDbErr = sql.Open("mysql", dnDSN)
	if MysqlDbErr != nil {
		log.Println("dbDSN: " + dnDSN)
		log.Fatal("数据源配置不正确: " + MysqlDbErr.Error())
	}
	MysqlDb.SetConnMaxLifetime(100)
	MysqlDb.SetMaxIdleConns(20)
	MysqlDb.SetConnMaxLifetime(10 * time.Second)
	if MysqlDbErr = MysqlDb.Ping(); nil != MysqlDbErr {
		log.Fatal("数据库连接失败: " + MysqlDbErr.Error())
	}
}
