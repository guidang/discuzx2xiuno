package app

import (
	"database/sql"
	"log"
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
func connectMysql(host *Hostinfo) (*sql.DB, error) {
	if host.DBPort == "" {
		host.DBPort = "3306"
	}

	if host.DBHost != "" {
		host.DBHost = "tcp(" + host.DBHost + ":" + host.DBPort + ")"
		log.Println(":::连接到 MySQL:" + host.DBHost)
	}

	if host.DBChar == "" {
		host.DBChar = "utf8"
		log.Println(":::MySQL 字符集为:" + host.DBChar)
	}

	db, err := sql.Open("mysql", host.DBUser+":"+host.DBPassword+"@"+host.DBHost+"/"+host.DBName+"?charset="+host.DBChar)
	return db, err
}

/**
  数据库加前缀
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