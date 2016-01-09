package app

import (
	"fmt"
)

const (
	DxPost = "pre_forum_post"
	XnPost = "bbs_post"
	Sid = "56912b6948e57"
)

/**
 xn 帖子表
 */
type Post struct {
	Tid int  //主题 id
	Pid int  //帖子 id
	Uid int  //用户 id
	IsFirst int  //主帖,与 thread.firstpid 呼应
	CreateDate int  //发帖时间戳
	UserIp int  //用户 ip ip2long()
	Sid string  //标识串
	Message string  //内容
}

/**
 dx 帖子表
 */
type DPost struct {
	Tid int  //主题 id
	Pid int  //帖子 id
	AuthorId int  //用户 id
	First int  //主帖
	Dateline int  //发帖时间戳
	UseIp string  //用户 ip
	Message string  //内容
	Fid int  //版块 id
	Subject  string  //标题
}

func ToPost() bool {

	//oldDB, newDB := CreateDB()

	selectSQL := "SELECT tid,pid,authorid,first,dateline,useip,message,fid,subject FROM " + DxPost// + " limit 1"
	Data, _ := OldDB.Query(selectSQL)

	insertData := `INSERT INTO ` + XnPost + ` (tid,pid,uid,isfirst,create_date,userip,sid,message) VALUES (?,?,?,?,?,?,?,?)`

	stmt, err := NewDB.Prepare(insertData)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("预插入数据", insertData)
		return false
	}

	for Data.Next() {
		d1 := &DPost{}
		err = Data.Scan(&d1.Tid, &d1.Fid, &d1.AuthorId, &d1.First, &d1.Dateline, &d1.UseIp, &d1.Message, &d1.Fid, &d1.Subject)

		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("取数据出错", selectSQL)
			return false
		}
		//fmt.Println(d1.Tid,d1.Fid,d1.AuthorId,d1.First,d1.Dateline,d1.UseIp,d1.Message,d1.Fid,d1.Subject)

		useIp := Ip2long(d1.UseIp)

		_, err = stmt.Exec(d1.Tid,d1.Pid,d1.AuthorId,d1.First,d1.Dateline,useIp,Sid,d1.Message)

		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("插入数据出错", insertData)
			return false
		}
	}

	return true;
}