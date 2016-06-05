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
	Fid,  //版块 id
	Threads,  //主题数
	CreateDate int64  //创建时间
	Brief,  //介绍
	Name string  //版块名
}

/**
 dx 版块表
 */
type DForum struct {
	Fid,  //版块 id
	Threads int64  //主题数
	Name string  //版块名称
}

func ToForum() string {
	fmt.Println(":::正在导入 forums...")

	//当前时间
	tmStr1 := time.Now().Unix()
	tmStr := strconv.FormatInt(tmStr1, 10)

	selectSQL := fmt.Sprintf("SELECT fid,name,threads FROM %s", DxForum)
	insertSQL := fmt.Sprintf("INSERT INTO  %s (fid, name, threads, brief, create_date) VALUES (?, ?, ?, '', '%s')", XnForum, tmStr)

	var clearErr error
	if clearErr = ClearTable(XnForum); clearErr != nil {
		return fmt.Sprintf(ClearErrMsg, XnForum, clearErr.Error())
	}

	data, err := OldDB.Query(selectSQL)
	if err != nil {
		return fmt.Sprintf(SelectSQLErr, err.Error(), selectSQL)
	}

	stmt, err := NewDB.Prepare(insertSQL)
	if err != nil {
		return fmt.Sprintf(PreInsSQLErr, err.Error(), insertSQL)
	}

	var insertCount int
	for data.Next() {
		d1 := &DForum{}
		err = data.Scan(&d1.Fid, &d1.Name, &d1.Threads)
		if err != nil {
			return fmt.Sprintf(SelectSQLErr, err.Error(), selectSQL)
		}

		_, err := stmt.Exec(d1.Fid, d1.Name, d1.Threads)
		if err != nil {
			return fmt.Sprintf(InsertErr, XnForum, err.Error())
		}

		insertCount++
	}

	return fmt.Sprintf(InsertSuccess, XnForum, insertCount)
}