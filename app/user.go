package app

import (
	"fmt"
)

const (
	DxUser = "pre_common_member"
	XnUser = "bbs_user"
)

/**
 xn 用户表
 */
type User struct {
	Uid int  //用户 id
	Gid int  //用户组 id
	Email string  //邮箱
	UserName string  //用户名
	Password string  //密码 好像md5
	Threads int  //主题数
	Posts int  //回复数
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

func ToUser() {
	//oldDB, newDB := data.CreateDB()

	selectSQL := "SELECT uid,groupid,email,username,password,regdate FROM " + DxUser + " WHERE uid > 1 and username != 'admin'"
	Data, _ := OldDB.Query(selectSQL)
	fmt.Println(selectSQL)

	insertData := `INSERT INTO ` + XnUser + ` (uid,gid,email,username,password,create_date,salt) VALUES (?,101,?,?,?,?,'581249')`

	stmt, err := NewDB.Prepare(insertData)
	if err != nil {
		fmt.Println(err.Error())
	}

	for Data.Next() {
		d1 := &DUser{}
		err = Data.Scan(&d1.Uid, &d1.GroupId, &d1.Email, &d1.UserName, &d1.Password, &d1.RegDate)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(d1.Uid, d1.GroupId, d1.Email, d1.UserName, d1.Password, d1.RegDate)

		_, err = stmt.Exec(d1.Uid, d1.Email, d1.UserName, d1.Password, d1.RegDate)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}