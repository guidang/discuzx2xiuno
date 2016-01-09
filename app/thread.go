package app

import (
	"fmt"
)

const (
	DxThread = "pre_forum_thread"
	XnThread = "bbs_thread"
)

/**
 xn 主题表
 */
type Thread struct {
	Fid int  //版块 id
	Tid int  //主题 id
	Uid int  //发帖者id
	UserIp int  //发帖者 ip
	Subject string  //标题
	CreateDate int  //发帖时间
	LastDate int  //最后回复时间
	Views int  //浏览数
	Posts int  //回复数
	FirstPId int  //主题帖 id
	LastUid int  //最后回复人
	LastPid int  //最后回复 id
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

func ToThread() bool {
	//oldDB, newDB := CreateDB()
	//查找最早、最迟的 PostId 及最迟的 PostIP
	selectPost := `SELECT (SELECT userip FROM ` + XnPost + ` WHERE tid = ? ORDER BY pid ASC LIMIT 1) AS userip, MIN(pid), (SELECT uid FROM ` + XnPost + ` WHERE tid = ? ORDER BY pid ASC LIMIT 1) AS lastip, MAX(pid) FROM ` + XnPost + ` WHERE tid = ?` //AND tid > 5488`
	var userIp,firstPid,lastUid,lastPid int

	//NewDB.QueryRow(selectPost,1,1,1).Scan(&userIp,&firstPid,&lastUid,&lastPid)
	//fmt.Println(userIp,firstPid,lastUid,lastPid)
	//return false

	selectSQL := "SELECT tid,fid,authorid,subject,dateline,lastpost,views,replies FROM " + DxThread //+ " WHERE tid > 5489"
	Data, _ := OldDB.Query(selectSQL)

	insertData := `INSERT INTO ` + XnThread + ` (fid,tid,uid,userip,subject,create_date,last_date,views,posts,firstpid,lastuid,lastpid) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`
	fmt.Println(insertData)

	stmt, err := NewDB.Prepare(insertData)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("预插入数据", insertData)
		return false
	}

	for Data.Next() {
		d1 := &DThread{}
		err = Data.Scan(&d1.Tid, &d1.Fid, &d1.AuthorId, &d1.Subject, &d1.Dateline, &d1.Lastpost, &d1.Views, &d1.Replies)

		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("取数据出错", selectSQL)
			return false
		}
		fmt.Println(d1.Tid,d1.Fid,d1.AuthorId,d1.Subject,d1.Dateline,d1.Lastpost,d1.Views,d1.Replies)

		err = NewDB.QueryRow(selectPost,d1.Tid,d1.Tid,d1.Tid).Scan(&userIp,&firstPid,&lastUid,&lastPid)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("取数据出错Post>>> tid:", d1.Tid, "\n",selectPost)
			continue
			//return false
		}

		d1.Fid = 1
		_, err = stmt.Exec(d1.Fid,d1.Tid,d1.AuthorId,userIp,d1.Subject,d1.Dateline,d1.Lastpost,d1.Views,d1.Replies,firstPid,lastUid,lastPid)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("插入数据出错Insert>>> tid:", d1.Tid, "\n",insertData)
			return false
		}
	}

	return true
}