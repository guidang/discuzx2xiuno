package app

import (
	"log"
	"fmt"
	"github.com/frustra/bbcode"
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

/**
  导入 posts 表
 */
func ToPost() (bool, string) {
	log.Println(":::正在导入 posts...")

	selectSQL := "SELECT tid,pid,authorid,first,dateline,useip,message,fid,subject FROM " + DxPost + " LIMIT 100"
	insertSQL := "INSERT INTO " + XnPost + " (tid,pid,uid,isfirst,create_date,userip,sid,message) VALUES (?,?,?,?,?,?,?,?)"

	var clearErr error
	if clearErr = ClearTable(XnPost); clearErr != nil {
		return false, fmt.Sprintf(ClearErrMsg, XnPost, clearErr)
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
		d1 := &DPost{}
		err = data.Scan(&d1.Tid, &d1.Fid, &d1.AuthorId, &d1.First, &d1.Dateline, &d1.UseIp, &d1.Message, &d1.Fid, &d1.Subject)
		if err != nil {
			return false, fmt.Sprintf(SelectErr, selectSQL, err)
		}

		useIp := Ip2long(d1.UseIp)

		compiler := bbcode.NewCompiler(true, true)
		d1.Message = compiler.Compile(d1.Message)

		_, err = stmt.Exec(d1.Tid,d1.Pid,d1.AuthorId,d1.First,d1.Dateline,useIp,Sid,d1.Message)
		if err != nil {
			return false, fmt.Sprintf(InsertErr, XnPost, err)
		}

		insertCount++
	}

	return true, fmt.Sprintf(InsertSuccess, XnPost, insertCount)
}