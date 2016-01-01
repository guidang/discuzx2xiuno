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

	selectSQL := "SELECT fid,name,threads FROM " + DxThread
	Data, _ := oldDB.Query(selectSQL)
}