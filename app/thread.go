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
	selectPost := `SELECT (SELECT userip FROM ` + XnPost + ` WHERE tid = ? AND isfirst = 1 LIMIT 1) as userip, (SELECT pid FROM bbs_post WHERE tid = ? AND isfirst = 1 LIMIT 1) AS minpid, pid, uid FROM ` + XnPost + ` WHERE tid = ? ORDER BY create_date DESC LIMIT 1`
	var userIp,firstPid,lastPid,lastUid int

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

	myThreadsSQL := "INSERT INTO " + XnMyThread + " VALUES (?,?)"
	myTStmt, err := NewDB.Prepare(myThreadsSQL)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("预插入数据-->mythreads", myThreadsSQL)
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

		err = NewDB.QueryRow(selectPost,d1.Tid,d1.Tid,d1.Tid).Scan(&userIp,&firstPid,&lastPid,&lastUid)
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

		_, err = myTStmt.Exec(d1.AuthorId, d1.Tid)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("插入数据出错Mythread>>> tid:", d1.Tid, "\n",myThreadsSQL)
			return false
		}

	}

	return true
}

/**
 修复主帖关联错误的问题
 */
func fixLastPid() {
	selectPost := `SELECT (SELECT userip FROM ` + XnPost + ` WHERE tid = ? AND isfirst = 1 LIMIT 1) as userip, (SELECT pid FROM bbs_post WHERE tid = ? AND isfirst = 1 LIMIT 1) AS minpid, pid, uid FROM ` + XnPost + ` WHERE tid = ? ORDER BY create_date DESC LIMIT 1`
	var userIp, firstPid, lastPid, lastUid int

	selectThread := `SELECT uid,tid FROM ` + XnThread// + ` WHERE tid > 92753`
	threadRun, _ := NewDB.Query(selectThread)

	updateThread := "UPDATE " + XnThread + " SET userip = ?, firstpid = ?, lastpid = ?, lastuid = ? WHERE tid = ?"
	utSmtm, _ := NewDB.Prepare(updateThread)

	for threadRun.Next() {
		var uid, tid int
		threadRun.Scan(&uid, &tid)

		NewDB.QueryRow(selectPost, tid, tid, tid).Scan(&userIp, &firstPid, &lastPid, &lastUid)

		fmt.Println(uid, tid, userIp, firstPid, lastPid, lastUid)
		_, err := utSmtm.Exec(userIp, firstPid, lastPid, lastUid, tid)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("插入数据出错thread>>> tid:", tid)
			return
		}
	}

}