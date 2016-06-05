package app

import (
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
	Tid,  //主题 id
	Pid,  //帖子 id
	Uid,  //用户 id
	IsFirst,  //主帖,与 thread.firstpid 呼应
	CreateDate,  //发帖时间戳
	UserIp int64  //用户 ip ip2long()
	Sid,  //标识串
	Message string  //内容
}

/**
 dx 帖子表
 */
type DPost struct {
	Tid,  //主题 id
	Pid,  //帖子 id
	AuthorId,  //用户 id
	First,  //主帖
	Dateline,  //发帖时间戳
	Fid int64  //版块 id
	UseIp,  //用户 ip
	Message,  //内容
	Subject  string  //标题
}

/**
  导入 posts 表
 */
func ToPost() string {
	fmt.Println(":::正在导入帖子 posts...")

	selectSQL := fmt.Sprintf("SELECT tid,pid,authorid,first,dateline,useip,message,fid,subject FROM %s", DxPost)// + " LIMIT 100"
	insertSQL := fmt.Sprintf("INSERT INTO %s (tid,pid,uid,isfirst,create_date,userip,sid,message) VALUES (?,?,?,?,?,?,?,?)", XnPost)

	var clearErr error
	if clearErr = ClearTable(XnPost); clearErr != nil {
		return fmt.Sprintf(ClearErrMsg, XnPost, clearErr.Error())
	}

	data, err := OldDB.Query(selectSQL)
	if err != nil {
		return fmt.Sprintf(SelectSQLErr, err.Error(), selectSQL)
	}

	stmt, err := NewDB.Prepare(insertSQL)
	if err != nil {
		return fmt.Sprintf(PreInsSQLErr, err.Error(), insertSQL)
	}

	//插入计数
	var insertCount int
	for data.Next() {
		d1 := &DPost{}
		err = data.Scan(&d1.Tid, &d1.Fid, &d1.AuthorId, &d1.First, &d1.Dateline, &d1.UseIp, &d1.Message, &d1.Fid, &d1.Subject)
		if err != nil {
			return fmt.Sprintf(SelectSQLErr, err.Error(), selectSQL)
		}

		//IP 转整型
		useIp := Ip2long(d1.UseIp)

		//bbcode 转 html
		compiler := bbcode.NewCompiler(true, true)
		d1.Message = compiler.Compile(d1.Message)

		tidStr := fmt.Sprintf(" tid:%s", d1.Tid)
		_, err = stmt.Exec(d1.Tid,d1.Pid,d1.AuthorId,d1.First,d1.Dateline,useIp,Sid,d1.Message)
		if err != nil {
			return fmt.Sprintf(InsertErr, XnPost, err.Error() + tidStr)
		}

		insertCount++
	}

	return fmt.Sprintf(InsertSuccess, XnPost, insertCount)
}