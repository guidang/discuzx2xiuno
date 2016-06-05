package app

import (
	"fmt"
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
	Fid,  //版块 id
	Tid,  //主题 id
	Uid,  //发帖者id
	UserIp,  //发帖者 ip
	CreateDate,  //发帖时间
	LastDate,  //最后回复时间
	Views,  //浏览数
	Posts,  //回复数
	FirstPId,  //主题帖 id
	LastUid,  //最后回复人
	LastPid int64  //最后回复 id
	Subject string  //标题
}

/**
 dx 主题表
 */
type DThread struct {
	Tid,  //主题 id
	Fid,  //版块 id
	AuthorId,  //发帖者id
	Dateline,  //发帖时间
	Lastpost,  //最后回复时间
	Views,  //浏览数
	Replies int64  //回复数
	Subject string  //标题
}

/**
	转换主题
 */
func ToThread() string {
	fmt.Println(":::正在导入主题 threads...")

	//查找最早、最迟的 PostId 及最迟的 PostIP
	selectThread := fmt.Sprintf(`SELECT
(SELECT userip FROM %s WHERE tid = ? AND isfirst = 1 LIMIT 1) AS userip,
(SELECT pid FROM bbs_post WHERE tid = ? AND isfirst = 1 LIMIT 1) AS minpid,
pid, uid
FROM %s
WHERE tid = ?
ORDER BY create_date DESC
LIMIT 1`, XnPost, XnPost)

	selectSQL := fmt.Sprintf("SELECT tid,fid,authorid,subject,dateline,lastpost,views,replies FROM %s", DxThread)  // + " LIMIT 100"
	insertSQL := fmt.Sprintf(`INSERT INTO %s (fid,tid,uid,userip,subject,create_date,last_date,views,posts,firstpid,lastuid,lastpid) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, XnThread)
	myThreadsSQL := fmt.Sprintf("INSERT INTO %s VALUES (?,?)", XnMyThread)

	var clearErr error
	if clearErr = ClearTable(XnThread); clearErr != nil {
		return fmt.Sprintf(ClearErrMsg, XnThread, clearErr.Error())
	}

	if clearErr = ClearTable(XnMyThread); clearErr != nil {
		return fmt.Sprintf(ClearErrMsg, XnMyThread, clearErr.Error())
	}

	var userIp,firstPid,lastPid,lastUid int

	data, err := OldDB.Query(selectSQL)
	if err != nil {
		return fmt.Sprintf(SelectSQLErr, err.Error(), selectSQL)
	}

	stmt, err := NewDB.Prepare(insertSQL)
	if err != nil {
		return fmt.Sprintf(PreInsSQLErr, err.Error(), insertSQL)
	}

	myTStmt, err := NewDB.Prepare(myThreadsSQL)
	if err != nil {
		return fmt.Sprintf(PreInsSQLErr, err.Error(), myThreadsSQL)
	}

	var insertCount int
	for data.Next() {
		d1 := &DThread{}
		err = data.Scan(&d1.Tid, &d1.Fid, &d1.AuthorId, &d1.Subject, &d1.Dateline, &d1.Lastpost, &d1.Views, &d1.Replies)
		if err != nil {
			return fmt.Sprintf(SelectSQLErr, err.Error(), selectSQL)
		}

		tidStr := fmt.Sprintf(" tid:%d", d1.Tid)

		err = NewDB.QueryRow(selectThread,d1.Tid,d1.Tid,d1.Tid).Scan(&userIp,&firstPid,&lastPid,&lastUid)
		if err != nil {
			fmt.Printf(SelectSQLErr + "\n", err.Error() + tidStr, XnPost)
			continue
		}

		_, err = stmt.Exec(d1.Fid,d1.Tid,d1.AuthorId,userIp,d1.Subject,d1.Dateline,d1.Lastpost,d1.Views,d1.Replies,firstPid,lastUid,lastPid)
		if err != nil {
			return fmt.Sprintf(InsertErr, insertSQL, err.Error() + tidStr)
		}

		_, err = myTStmt.Exec(d1.AuthorId, d1.Tid)
		if err != nil {
			return fmt.Sprintf(InsertErr, myThreadsSQL, err.Error() + tidStr)
		}

		insertCount++

	}

	return fmt.Sprintf(InsertSuccess, XnThread, insertCount)
}