package app

/**
 xn 版块表
 */
type  Forum struct  {
	Fid int  //版块 id
	Name string  //版块名
	Threads int  //主题数
	Brief string  //介绍
	CreateDate int  //创建时间
}

/**
 dx 版块表
 */
type DForum struct {
	Fid int  //版块 id
	Name string  //版块名称
	Threads int  //主题数
}