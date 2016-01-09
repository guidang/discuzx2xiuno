/**
 版块表 - 作废(因 xiuno 版块限制)
 */

package app

import (
	"fmt"

	"time"
	"strconv"
)

const (
	DxForum = "pre_forum_forum"
	XnForum = "bbs_forum"
)

/**
 xn 版块表
 */
type Forum struct  {
	Fid int  //版块 id
	Name string  //版块名
	Threads int  //主题数
	Brief string  //介绍
	CreateDate int  //创建时间
}

/**
 dx 版块表
 */
type DForum struct {
	Fid int  //版块 id
	Name string  //版块名称
	Threads int  //主题数
}

func ToForum()  {
	oldDB, newDB := CreateDB()

	selectSQL := "SELECT fid,name,threads FROM " + DxForum
	Data, _ := oldDB.Query(selectSQL)

	tmStr1 := time.Now().Unix()
	tmStr := strconv.FormatInt(tmStr1, 10)

	insertData := `INSERT INTO ` + XnForum + `(fid, name, threads, brief, create_date) VALUES (?, ?, ?, '', '` + tmStr + `')`
	fmt.Println(insertData)

	stmt, err := newDB.Prepare(insertData)
	if err != nil {
		fmt.Println(err.Error())
	}

	for Data.Next() {
		d1 := &DForum{}
		err = Data.Scan(&d1.Fid, &d1.Name, &d1.Threads)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(d1.Fid, d1.Name, d1.Threads)

		_, err := stmt.Exec(d1.Fid, d1.Name, d1.Threads)
		if err != nil {
			fmt.Println(err)
		}

	}
}