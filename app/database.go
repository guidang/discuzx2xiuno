package app

import (
	"database/sql"
	"fmt"
	"strings"
)

type Hostinfo struct {
	DBUser,
	DBPassword,
	DBName,
	DBHost,
	DBPort,
	DBChar string
}

/**
  连接数据库
*/
func connectMysql(host *Hostinfo) (db *sql.DB, err error) {
	if host.DBPort == "" {
		host.DBPort = "3306"
	}

	if host.DBChar == "" {
		host.DBChar = "utf8"
	}

	if host.DBHost != "" {
		host.DBHost = "tcp(" + host.DBHost + ":" + host.DBPort + ")"
	}

	dbStr := fmt.Sprintf("%s:%s@%s/%s?%s",
		host.DBUser,
		host.DBPassword,
		host.DBHost,
		host.DBName,
		host.DBChar,
	)

	db, err = sql.Open("mysql", dbStr)
	return
}

/**
  数据库字段批量加前缀
*/
func FieldAddPrev(prev, fieldStr string) string {
	fieldArr := strings.Split(fieldStr, ",")

	prev = prev + "."
	var newFieldArr []string
	for _, v := range fieldArr {
		newFieldArr = append(newFieldArr, prev+v)
	}
	newFieldStr := strings.Join(newFieldArr, ",")

	return newFieldStr
}
