package app
import "fmt"

const (
	DxUcUser = "uc_members"
)

type ucUser struct {
	uid int
	salt,password,regip string
}

func UpdateUser() (bool,string) {
	ucUserSQL := "UPDATE " + XnUser + " SET password = ?, salt = ?, create_ip = ? WHERE uid = ?"

	stmt, err := NewDB.Prepare(ucUserSQL)
	if err != nil {
		fmt.Println(err.Error())
		return false, "updateUser.预处理失败"
	}

	selectMyUserSQL := "SELECT uid FROM " + XnUser
	uidData, err := NewDB.Query(selectMyUserSQL)

	selectUcUserSQL := "SELECT uid,password,regip,salt FROM " + DxUcUser
	ucUserData, err := OldDB.Query(selectUcUserSQL)

	//从 uc_member 取出用户数据
	myUser := make(map[int]ucUser)
	for ucUserData.Next() {
		u := &ucUser{}
		err = ucUserData.Scan(&u.uid,&u.password,&u.regip,&u.salt)
		if err != nil {
			fmt.Println(err.Error())
			return false, "selectUcUser.失败"
		}
		myUser[u.uid] = ucUser{u.uid,u.salt,u.password,u.regip}

	}
	//fmt.Println(myUser)

	//从 bbs_user 取用户信息
	var uid int
	for uidData.Next() {
		err = uidData.Scan(&uid)
		user := myUser[uid]

		//已经从 uc_member 删除的用户
		if user.uid == 0 {
			continue
		}

		userRegip := Ip2long(user.regip)
		_, err = stmt.Exec(user.password, user.salt, userRegip, uid)
		//fmt.Println(user.uid, user.password)
		if err != nil {
			fmt.Println(err.Error())
			return false, "selectMyUser.查找用户失败:" + string(uid)
		}
	}

	return true, "User password,salt,regip Update Success"
}