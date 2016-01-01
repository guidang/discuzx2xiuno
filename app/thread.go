package app

import (
	"fmt"
	"./data"
)

const (
	DxThread = "pre_forum_thread"
	XnThread = "bbs_thread"
)

/**
 xn 主题表
 */
type Thread struct {
	Tid int  //主题 id
	Fid int  //版块 id
	Uid int  //发帖者id
	Subject string  //标题
	CreateDate int  //发帖时间
	LastDate int  //最后回复时间
	Views int  //浏览数
	Posts int  //回复数
	UserIp int  //发帖者 ip
}

/**
 dx 主题表
 */
type DThread struct {
	Tid int  //主题 id
	Fid int  //版块 id
	AuthorId int  //发帖者id
	Subject string  //标题
	Dateline int  //发帖时间
	Lastpost int  //最后回复时间
	Views int  //浏览数
	Replies int  //回复数
}

func ToThread()  {
	oldDB, newDB := data.CreateDB()

	selectSQL := "SELECT tid,fid,authorid,subject,dateline,lastpost,views,replies FROM " + DxThread
	Data, _ := oldDB.Query(selectSQL)

	insertData := `INSERT INTO ` + XnThread + ` (tid,fid,uid,subject,create_date,last_date,views,posts,userip) VALUES (?,?,?,?,?,?,?,?,'2130706433')`
	fmt.Println(insertData)

	stmt, err := newDB.Prepare(insertData)
	if err != nil {
		fmt.Println(err.Error())
	}

	for Data.Next() {
		d1 := &DThread{}
		err = Data.Scan(&d1.Tid, &d1.Fid, &d1.AuthorId, &d1.Subject, &d1.Dateline, &d1.Lastpost, &d1.Views, &d1.Replies)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(d1.Tid,d1.Fid,d1.AuthorId,d1.Subject,d1.Dateline,d1.Lastpost,d1.Views,d1.Replies)

		_, err = stmt.Exec(d1.Tid,d1.Fid,d1.AuthorId,d1.Subject,d1.Dateline,d1.Lastpost,d1.Views,d1.Replies)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}