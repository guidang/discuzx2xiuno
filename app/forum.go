package app

import (
	"fmt"
	"log"

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

func ToForum() (bool, string) {
	log.Println(":::正在导入 forums...")

	//当前时间
	tmStr1 := time.Now().Unix()
	tmStr := strconv.FormatInt(tmStr1, 10)

	selectSQL := "SELECT fid,name,threads FROM " + DxForum
	insertSQL := "INSERT INTO " + XnForum + "(fid, name, threads, brief, create_date) VALUES (?, ?, ?, '', '" + tmStr + "')"

	var clearErr error
	if clearErr = ClearTable(XnForum); clearErr != nil {
		return false, fmt.Sprintf(ClearErrMsg, XnForum, clearErr)
	}

	data, err := OldDB.Query(selectSQL)
	if err != nil {
		return false, fmt.Sprintf(SelectErr, selectSQL, err)
	}

	stmt, err := NewDB.Prepare(insertSQL)
	if err != nil {
		return false, fmt.Sprintf(PreInsertErr, insertSQL, err)
	}

	var insertCount int
	for data.Next() {
		d1 := &DForum{}
		err = data.Scan(&d1.Fid, &d1.Name, &d1.Threads)
		if err != nil {
			return false, fmt.Sprintf(SelectErr, selectSQL, err)
		}

		_, err := stmt.Exec(d1.Fid, d1.Name, d1.Threads)
		if err != nil {
			return false, fmt.Sprintf(InsertErr, XnForum, err)
		}

		insertCount++
	}

	return true, fmt.Sprintf(InsertSuccess, XnForum, insertCount)
}