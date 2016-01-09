package app

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

func CreateDB() (olddb, newdb *sql.DB) {
	olddb, err := connMysql("gxvtc")

	if err != nil {
		fmt.Println(err)
	}

	newdb, err = connMysql("xiuno")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("connectDb")
	return olddb, newdb
}

/**
 连接数据库
 */
func connMysql(dbs string) (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:123456@/" + dbs)
	return db, err
}