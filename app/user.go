package app

import (
	"fmt"
	"log"
)

const (
	DxUser = "pre_common_member"
	DxUcUser = "pre_ucenter_members"
	DxUserStatus = "pre_common_member_status"
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
	Uid,  //用户 id
	GroupId,  //用户组 id
	RegDate int  //注册时间
	Email,  //邮箱
	UserName,  //用户名
	Password,  //密码
	Salt,  //加密 key
	UcPassword,  //uc中的密码
	Regip,  //注册 ip
	Lastip string  //最后登陆 ip
	Lastvisit int  //最后登陆时间
}

//按字段分组
type NewUsers struct {
	Uids      []int
	Emails    []string
	UserNames []string
}

//用户信息
type UserInfo struct {
	Uid int
	Email,
	Username string
}

//一组用户数据::按用户分组
var userInfos []UserInfo

var (
	selectPostTotalSQL = "SELECT (SELECT count(*) FROM " + XnThread + " WHERE uid = ?) AS mythreads, (SELECT COUNT(*) FROM " + XnPost + " WHERE uid = ?) AS myposts"
	insertPostTotalSQL = "UPDATE " + XnUser + " SET threads = ?, posts = ? WHERE uid = ?"
	selectUserPostSQL = "SELECT uid FROM " + XnUser
)

func ToUser() (bool, string) {
	log.Println(":::正在导入 users...")

	/*
	SELECT m.uid, m.groupid, m.email, m.username, m.password, u.salt, u.password, s.regip, m.regdate, s.lastip, s.lastvisit FROM pre_common_member m LEFT JOIN uc_members u ON u.uid = m.uid LEFT JOIN pre_common_member_status s ON s.uid = m.uid m.uid < 10
	*/

	mField := FieldAddPrev("m", "uid,groupid,email,username,password,regdate")
	uField := FieldAddPrev("u", "salt,password")
	sField := FieldAddPrev("s", "regip,lastip,lastvisit")
	selectSQL := "SELECT " + mField + "," + uField + "," + sField + " FROM " + DxUser + " m LEFT JOIN " + DxUcUser + " u ON u.uid = m.uid LEFT JOIN " + DxUserStatus + " s ON s.uid = m.uid"// WHERE m.uid < 10"
	insertSQL := "INSERT INTO " + XnUser + " (uid,gid,email,username,password,salt,create_ip,create_date,login_ip,login_date) VALUES (?,101,?,?,?,?,?,?,?,?)"

	var clearErr error
	if clearErr = ClearTable(XnUser); clearErr != nil {
		return false, fmt.Sprintf(ClearErrMsg, XnUser, clearErr)
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
		d1 := &DUser{}
		err = data.Scan(&d1.Uid, &d1.GroupId, &d1.Email, &d1.UserName, &d1.Password, &d1.RegDate,&d1.Salt, &d1.UcPassword, &d1.Regip, &d1.Lastip, &d1.Lastvisit)
		if err != nil {
			return false, fmt.Sprintf(SelectErr, selectSQL, err)
		}

		createIp := Ip2long(d1.Regip)
		loginIp := Ip2long(d1.Lastip)
		_, err = stmt.Exec(d1.Uid,d1.Email,d1.UserName,d1.UcPassword,d1.Salt,createIp,d1.RegDate,loginIp,d1.Lastvisit)
		if err != nil {
			return false, fmt.Sprintf(InsertErr, XnUser, err)
		}

		insertCount++
	}

	//用户主帖和回复统计
	return true, fmt.Sprintf(InsertSuccess, XnUser, insertCount)
}

/**
 更新全部用户帖子数量
 */
func doUserPosts() (bool, string) {
	log.Println(":::正在更新 users 帖子统计...")
	data, err := NewDB.Query(selectUserPostSQL)
	if err != nil {
		return false, fmt.Sprintf(SelectErr, selectUserPostSQL, err)
	}

	var insertCount int
	for data.Next() {
		var uid int
		err = data.Scan(&uid)
		if err != nil {
			return false, fmt.Sprintf(SelectErr, selectUserPostSQL, err)
		}

		res, msg := updatePostTotal(uid)
		if res != true {
			log.Println(msg)
			continue
		}
		insertCount++
	}

	return true, fmt.Sprintf(InsertSuccess, XnUser, insertCount)
}
/**
  更新指定用户帖子数量
  uid 用户 id
 */
func updatePostTotal(uid int) (bool, string) {
	var myThreads,myPosts int// = 0,0
	NewDB.QueryRow(selectPostTotalSQL, uid, uid).Scan(&myThreads,&myPosts)
	stmt, err := NewDB.Prepare(insertPostTotalSQL)
	if err != nil {
		return false, fmt.Sprintf(PreInsertErr, insertPostTotalSQL, err)
	}
	_, err = stmt.Exec(myThreads, myPosts, uid)
	if err != nil {
		return false, fmt.Sprintf(InsertErr, XnUser, err)
	}

	return true, fmt.Sprintf(InsertSuccess, uid, 1)
}