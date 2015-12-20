package app

/**
 xn 用户表
 */
type User struct {
	Uid int  //用户 id
	Gid int  //用户组 id
	Email string  //邮箱
	UserName string  //用户名
	Password string  //密码 好像md5
	Salt int  //加密码
	CreateIp int  //创建 ip
	CreateDate int  //创建时间
	LoginIp int  //登陆 ip
	LoginDate int  //登陆日期
}

/**
 pre_common_member
 dx 用户表
 */
type DUser struct {
	Uid int  //用户 id
	GroupId int  //用户组 id
	Email string  //邮箱
	UserName string  //用户名
	Password string  //密码
	RegDate int  //注册时间
}