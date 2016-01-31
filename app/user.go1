package app

import (
	"fmt"
	"strings"
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

//按字段分组
type NewUsers struct {
	Uids      []int
	Emails    []string
	UserNames []string
}

//按用户分组
type UserInfo struct {
	Uid int
	Email,
	Username string
}

var userInfos []UserInfo

func ToUser() bool {
	//oldDB, newDB := data.CreateDB()
	//用户主帖和回复统计
	selectTotal := "SELECT (SELECT count(*) FROM `bbs_thread` WHERE uid = ?) AS mythreads, (SELECT COUNT(*) FROM bbs_post WHERE uid = ?) AS myposts"

	selectSQL := "SELECT uid,groupid,email,username,password,regdate FROM " + DxUser + " WHERE uid > 21690"
	Data, _ := OldDB.Query(selectSQL)
	fmt.Println(selectSQL)

	insertData := `INSERT INTO ` + XnUser + ` (uid,gid,email,username,password,create_date,salt,threads,posts) VALUES (?,101,?,?,?,?,'581249',?,?)`

	stmt, err := NewDB.Prepare(insertData)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//初始化用户资料
	newUser := &NewUsers{}
	//按用户分组

	postSQL := "UPDATE " + XnPost + " SET uid = ? WHERE uid = ?"
	threadSQL := "UPDATE " + XnThread + " SET uid = ? WHERE uid = ?"
	myThreadSQL := "UPDATE " + XnMyThread + " SET uid = ? WHERE uid = ?"

	PostPre, _ := NewDB.Prepare(postSQL)
	ThreadPre, _ := NewDB.Prepare(threadSQL)
	MyThreadPre, _ := NewDB.Prepare(myThreadSQL)

	//从数据库查找的开关,出错后中断时使用
	fromUser := true
	if fromUser == true {
		userInfos = selectUserList()
		for _, v := range userInfos {
			newUser.UserNames = append(newUser.UserNames, v.Username)
		}
	}

	for Data.Next() {
		d1 := &DUser{}
		err = Data.Scan(&d1.Uid, &d1.GroupId, &d1.Email, &d1.UserName, &d1.Password, &d1.RegDate)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		sameEmail := false
		//剃除相同邮箱的帐号,更新 post 和 thread形式
		for _, v := range userInfos {
			//转换成小写再对比
			d1.Email = strings.ToLower(d1.Email)
			v.Email = strings.ToLower(v.Email)

			if d1.Email == v.Email {
				sameEmail = true
				//bol := updateAccount(d1.Uid, v.Uid)

				_, p1 := PostPre.Exec(v.Uid, d1.Uid)
				_, p2 := ThreadPre.Exec(v.Uid, d1.Uid)
				_, p3 := MyThreadPre.Exec(v.Uid, d1.Uid)

				//如果更新失败则
				if p1 != nil || p2 != nil || p3 != nil {
					fmt.Println("error.email: ", d1.Email, d1.Uid)
				}

				break
			}
		}

		//去除老论坛中的"_s"后缀,转换成小写
		d1.UserName = strings.ToLower(strings.Replace(d1.UserName, "_s", "", -1))

		//如果不相同的邮箱则添加进来
		if sameEmail == false {
			userInfos = append(userInfos, UserInfo{d1.Uid, d1.Email, d1.UserName})
		} else {
			//相同则跳出以下操作
			continue
		}

		//方式二
		//处理邮箱 - 新用户名的形式
		//email = replaceData(newUser.Emails, email, 1)
		//添加 email 到数组中
		//newUser.Emails = append(newUser.Emails, email)

		//相同的用户名则
		for _, v := range newUser.UserNames {
			if d1.UserName == strings.ToLower(v) {
				d1.UserName = "old." + d1.UserName
			}
		}
		//添加 username 到数组中
		newUser.UserNames = append(newUser.UserNames, d1.UserName)

		var myThreads,myPosts int// = 0,0
		NewDB.QueryRow(selectTotal, d1.Uid, d1.Uid).Scan(&myThreads,&myPosts)

		fmt.Println(d1.Uid, d1.GroupId, d1.Email, d1.UserName, d1.Password, d1.RegDate, myThreads,myPosts)

		_, err = stmt.Exec(d1.Uid, d1.Email, d1.UserName, d1.Password, d1.RegDate, myThreads, myPosts)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}

	return true
}

/**
 * plist 列表, pstr 字符串, ptype 类型(1.邮箱,2.用户名)
 * 重命名方式
 */
func replaceData(plist []string, pstr string, ptype int) string {
	//fmt.Println(pstr)
	for _, v := range plist {

		if pstr == "" {

			//邮箱则改为通用邮箱
			if ptype == 1 {
				pstr = "guest@gxvtc.com"
			}
		}

		//如果存在邮箱则替换
		if pstr == v {
			tmp := "old." + pstr
			pstr = replaceData(plist, tmp, ptype)
		}
	}

	fmt.Println(pstr)
	//plist = append(plist, pstr)
	return pstr
}

func updateAccount(poldUid, pnewUid int) bool {
	fmt.Println(poldUid, pnewUid)
	return false
}

func ToMyThreads() bool {
	selectThreads := "SELECT uid,tid FROM " + XnThread
	mydata, _ := NewDB.Query(selectThreads)

	insertMyThreads := "INSERT INTO " + XnMyThread + " VALUES (?,?)"
	stmt, err := NewDB.Prepare(insertMyThreads)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for mydata.Next() {
		var uid,tid int
		err = mydata.Scan(&uid, &tid)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		_, err = stmt.Exec(uid,tid)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

	}

	return true
}

func selectUserList() []UserInfo {
	userSQL := "SELECT uid,email,username FROM " + XnUser
	data, _ := NewDB.Query(userSQL)
	for data.Next() {
		var uid int
		var email, username string
		data.Scan(&uid,&email,&username)
		userInfos = append(userInfos, UserInfo{uid, email, username})
	}

	return userInfos
}