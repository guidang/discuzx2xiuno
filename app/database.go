package app

import (
	"database/sql"
	"log"
)

type Hostinfo struct {
	DBUser,
	DBPassword,
	DBname,
	DBHost,
	DBPort,
	DBChar string
}

/**
  连接数据库
*/
func connectMysql(host *Hostinfo) (*sql.DB, error) {
	if host.DBHost != "" {
		host.DBHost = "tcp(" + host.DBHost + ":" + host.DBPort + ")"
		log.Println(":::连接到 MySQL:" + host.DBHost)
	}

	if host.DBChar == "" {
		host.DBChar = "utf8"
		log.Println(":::MySQL 字符集为:" + host.DBChar)
	}

	db, err := sql.Open("mysql", host.DBUser+":"+host.DBPassword+"@"+host.DBHost+"/"+host.DBname+"?charset="+host.DBChar)
	return db, err
}