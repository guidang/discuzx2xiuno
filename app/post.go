package app

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
	UserIp string  //用户 ip
	Message string  //内容
	Fid int  //版块 id
	Subject  string  //标题
}