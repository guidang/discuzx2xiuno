package app

import (
	"fmt"
)

const (
	DxUser = "pre_common_member"
	DxUcUser = "pre_ucenter_members"
	DzUcUser = "uc_members"
	DxUserStatus = "pre_common_member_status"
	XnUser = "bbs_user"
)

/**
 xn 用户表
 */
type User struct {
	Uid,  //用户 id
	Gid,  //用户组 id
	Threads,  //主题数
	Posts,  //回复数
	Salt,  //加密码
	CreateIp,  //创建 ip
	CreateDate,  //创建时间
	LoginIp,  //登陆 ip
	LoginDate int64  //登陆日期
	Email,  //邮箱
	UserName,  //用户名
	Password string  //密码 好像md5
}

/**
 pre_common_member
 dx 用户表
 */
type DUser struct {
	Uid,  //用户 id
	GroupId,  //用户组 id
	RegDate,  //注册时间
	Lastvisit int64  //最后登陆时间
	Email,  //邮箱
	UserName,  //用户名
	Password,  //密码
	Salt,  //加密 key
	UcPassword,  //uc中的密码
	Regip,  //注册 ip
	Lastip string  //最后登陆 ip
}

//按字段分组
type NewUsers struct {
	Uids      []int64
	Emails    []string
	UserNames []string
}

//用户信息
type UserInfo struct {
	Uid int64
	Email,
	Username string
}

var (
	selectPostTotalSQL = "SELECT (SELECT count(*) FROM " + XnThread + " WHERE uid = ?) AS mythreads, (SELECT COUNT(*) FROM " + XnPost + " WHERE uid = ?) AS myposts"
	insertPostTotalSQL = "UPDATE " + XnUser + " SET threads = ?, posts = ? WHERE uid = ?"
	selectUserPostSQL = "SELECT uid FROM " + XnUser
)

/**
	转换用户
 */
func ToUser() string {
	fmt.Println(":::正在导入 users...")

	/*
	SELECT m.uid, m.groupid, m.email, m.username, m.password, u.salt, u.password, s.regip, m.regdate, s.lastip, s.lastvisit FROM pre_common_member m LEFT JOIN uc_members u ON u.uid = m.uid LEFT JOIN pre_common_member_status s ON s.uid = m.uid m.uid < 10
	*/

	oldUsers := DxUcUser

	if Exts.UpdateFromDz {
		oldUsers = DzUcUser
	}

	mField := FieldAddPrev("m", "uid,groupid,email,username,password,regdate")
	uField := FieldAddPrev("u", "salt,password")
	sField := FieldAddPrev("s", "regip,lastip,lastvisit")
	selectSQL := "SELECT " + mField + "," + uField + "," + sField + " FROM " + DxUser + " m LEFT JOIN " + oldUsers + " u ON u.uid = m.uid LEFT JOIN " + DxUserStatus + " s ON s.uid = m.uid ORDER BY m.uid ASC"// WHERE m.uid < 10"
	insertSQL := "INSERT INTO " + XnUser + " (uid,gid,email,username,password,salt,create_ip,create_date,login_ip,login_date) VALUES (?,101,?,?,?,?,?,?,?,?)"

	var clearErr error
	if clearErr = ClearTable(XnUser); clearErr != nil {
		return fmt.Sprintf(ClearErrMsg, XnUser, clearErr)
	}

	data, err := OldDB.Query(selectSQL)
	if err != nil {
		return fmt.Sprintf(SelectSQLErr, err.Error(), selectSQL)
	}

	stmt, err := NewDB.Prepare(insertSQL)
	if err != nil {
		return fmt.Sprintf(PreInsSQLErr, err.Error(), insertSQL)
	}

	var insertCount int
	for data.Next() {
		d1 := &DUser{}
		var salt, password []byte
		err = data.Scan(&d1.Uid, &d1.GroupId, &d1.Email, &d1.UserName, &d1.Password, &d1.RegDate,&salt,&password, &d1.Regip, &d1.Lastip, &d1.Lastvisit)
		if err != nil {
			return fmt.Sprintf(SelectSQLErr, err.Error(), selectSQL)
		}

		if salt == nil {
			d1.Salt = "111111"
		} else {
			d1.Salt = string(salt)	
		}

		if password == nil {
			d1.UcPassword = d1.Password
		} else {
			d1.UcPassword = string(password)
		}

		createIp := Ip2long(d1.Regip)
		loginIp := Ip2long(d1.Lastip)
		_, err = stmt.Exec(d1.Uid,d1.Email,d1.UserName,d1.UcPassword,d1.Salt,createIp,d1.RegDate,loginIp,d1.Lastvisit)
		if err != nil {
			return fmt.Sprintf(InsertErr, XnUser, err.Error())
		}

		insertCount++
	}

	//用户主帖和回复统计
	return fmt.Sprintf(InsertSuccess, XnUser, insertCount)
}

/**
	更新全部用户帖子数量
 */
func doUserPosts() string {
	fmt.Println(":::正在更新 users 帖子统计...")
	data, err := NewDB.Query(selectUserPostSQL)
	if err != nil {
		return fmt.Sprintf(SelectSQLErr, err.Error(), selectUserPostSQL)
	}

	var insertCount int
	for data.Next() {
		var uid int
		err = data.Scan(&uid)
		if err != nil {
			return fmt.Sprintf(SelectSQLErr, err.Error(), selectUserPostSQL)
		}

		res, msg := updatePostTotal(uid)
		if !res {
			fmt.Println(msg)
			continue
		}
		insertCount++
	}

	return fmt.Sprintf(InsertSuccess, XnUser, insertCount)
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
		return false, fmt.Sprintf(PreInsSQLErr, err.Error(), insertPostTotalSQL)
	}
	_, err = stmt.Exec(myThreads, myPosts, uid)
	if err != nil {
		return false, fmt.Sprintf(InsertErr, XnUser, err.Error())
	}

	return true, fmt.Sprintf(InsertSuccess, uid, 1)
}

/**
	更新管理员帐号
 */
func updateAdminUser() string {
	adminSQL := fmt.Sprintf("UPDATE %s SET gid = 1 WHERE uid = ?", XnUser)
	stmt, err := NewDB.Prepare(adminSQL)
	if err != nil {
		return fmt.Sprintf(PreInsSQLErr, err.Error(), adminSQL)
	}
	_, err = stmt.Exec(Exts.AdminUid)
	if err != nil {
		return fmt.Sprintf(InsertErr, XnUser, err.Error())
	}

	return fmt.Sprintf(UpAdminSuccess, "管理员", Exts.AdminUid)
}