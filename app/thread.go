package app

import (
	"fmt"
	"log"
)

const (
	DxThread = "pre_forum_thread"
	XnThread = "bbs_thread"
	XnMyThread = "bbs_mythread"
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

func ToThread() (bool, string) {
	log.Println(":::正在导入 threads...")

	//查找最早、最迟的 PostId 及最迟的 PostIP
	selectThread := `SELECT (SELECT userip FROM ` + XnPost + ` WHERE tid = ? AND isfirst = 1 LIMIT 1) as userip, (SELECT pid FROM bbs_post WHERE tid = ? AND isfirst = 1 LIMIT 1) AS minpid, pid, uid FROM ` + XnPost + ` WHERE tid = ? ORDER BY create_date DESC LIMIT 1`
	selectSQL := "SELECT tid,fid,authorid,subject,dateline,lastpost,views,replies FROM " + DxThread + " LIMIT 100"
	insertSQL := `INSERT INTO ` + XnThread + ` (fid,tid,uid,userip,subject,create_date,last_date,views,posts,firstpid,lastuid,lastpid) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`
	myThreadsSQL := "INSERT INTO " + XnMyThread + " VALUES (?,?)"

	var clearErr error
	if clearErr = ClearTable(XnThread); clearErr != nil {
		return false, fmt.Sprintf(ClearErrMsg, XnThread, clearErr)
	}
	if clearErr = ClearTable(XnMyThread); clearErr != nil {
		return false, fmt.Sprintf(ClearErrMsg, XnMyThread, clearErr)
	}

	var userIp,firstPid,lastPid,lastUid int

	data, err := OldDB.Query(selectSQL)
	if err != nil {
		return false, fmt.Sprintf(SelectErr, XnPost, err)
	}

	stmt, err := NewDB.Prepare(insertSQL)
	if err != nil {
		return false, fmt.Sprintf(PreInsertErr, insertSQL, err)
	}

	myTStmt, err := NewDB.Prepare(myThreadsSQL)
	if err != nil {
		return false, fmt.Sprintf(PreInsertErr, myThreadsSQL, err)
	}

	var insertCount int
	for data.Next() {
		d1 := &DThread{}
		err = data.Scan(&d1.Tid, &d1.Fid, &d1.AuthorId, &d1.Subject, &d1.Dateline, &d1.Lastpost, &d1.Views, &d1.Replies)
		if err != nil {
			return false, fmt.Sprintf(SelectErr, selectSQL, err)
		}

		err = NewDB.QueryRow(selectThread,d1.Tid,d1.Tid,d1.Tid).Scan(&userIp,&firstPid,&lastPid,&lastUid)
		if err != nil {
			fmt.Printf(SelectErr + "\n", XnPost, err)
			continue
		}

		_, err = stmt.Exec(d1.Fid,d1.Tid,d1.AuthorId,userIp,d1.Subject,d1.Dateline,d1.Lastpost,d1.Views,d1.Replies,firstPid,lastUid,lastPid)
		if err != nil {
			return false, fmt.Sprintf(InsertErr, insertSQL, err)
		}

		_, err = myTStmt.Exec(d1.AuthorId, d1.Tid)
		if err != nil {
			return false, fmt.Sprintf(InsertErr, myThreadsSQL, err)
		}

		insertCount++

	}

	return true, fmt.Sprintf(InsertSuccess, XnThread,insertCount)
}