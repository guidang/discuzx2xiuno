package app

import (
	"fmt"
	"encoding/binary"

	"./data"
	"net"
)

const (
	DxPost = "pre_forum_post"
	XnPost = "bbs_post"
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

func ToPost() {
	oldDB, newDB := data.CreateDB()

	selectSQL := "SELECT tid,pid,authorid,first,dateline,useip,message,fid,subject FROM " + DxPost// + " limit 1"
	Data, _ := oldDB.Query(selectSQL)

	insertData := `INSERT INTO ` + XnPost + ` (tid,pid,uid,isfirst,create_date,userip,message) VALUES (?,?,?,?,?,?,?)`

	stmt, err := newDB.Prepare(insertData)
	if err != nil {
		fmt.Println(err.Error())
	}

	for Data.Next() {
		d1 := &DPost{}
		err = Data.Scan(&d1.Tid, &d1.Fid, &d1.AuthorId, &d1.First, &d1.Dateline, &d1.UseIp, &d1.Message, &d1.Fid, &d1.Subject)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(d1.Tid,d1.Fid,d1.AuthorId,d1.First,d1.Dateline,d1.UseIp,d1.Message,d1.Fid,d1.Subject)

		useIp := ip2long(d1.UseIp)
		_, err = stmt.Exec(d1.Tid,d1.Pid,d1.AuthorId,d1.First,d1.Dateline,useIp,d1.Message)

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func ip2long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

func long2ip(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}