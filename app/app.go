package app

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//数据库初始化
var (
	OldDB,
	NewDB *sql.DB
	ClearTB = true  //是否先清理表
	MergeUser = true  //是否合并用户
	ResetPost = false
)

func Init() {
	log.Println(":::正在进入app主程序:::")
	OldDB, NewDB = connDB()

	_, msg := ToPost()
	log.Println(msg)

	_, msg = ToThread()
	log.Println(msg)

	_, msg = ToForum()
	log.Println(msg)

	_, msg = ToUser()
	log.Println(msg)

	/* 更新全部用户帖子数量 */
	if ResetPost {
		_, msg := doUserPosts()
		log.Println(msg)
	}
}

/**
   连接新旧数据库
 */
func connDB() (*sql.DB, *sql.DB) {
	log.Println(":::正在连接数据库:::")

	old := &Hostinfo{
		DBUser: "root",
		DBPassword: "123456",
		DBName: "gxvtc",
	}

	new := &Hostinfo{
		DBUser: "root",
		DBPassword: "123456",
		DBName: "xiuno",
	}

	oldDB, err := connectMysql(old)
	if err != nil {
		log.Println("old db connect err: " + err.Error())
	}

	newDB, err := connectMysql(new)
	if err != nil {
		log.Println("new db connect err: " + err.Error())
	}

	return oldDB, newDB
}

/**
  清理 表
 */
func ClearTable(tbname string) error {
	log.Println(":::正在清理 " + tbname + " 表:::")
	clearSQL := "TRUNCATE TABLE " + tbname

	_, err := NewDB.Exec(clearSQL)
	if err != nil {
		log.Println(":::清理 " + tbname + " 失败: " + err.Error())
	}

	return err
}

/*

func Init3() {
	ToThread()
}

func Init()  {

	//导入版块
	//ToForum()

	//开关
	isRun := false

	isPost :=  false
	//导入帖子
	if isRun == true {
		isPost = ToPost()
	}

	isThread := false
	//导入主题
	if isPost == true {
		isThread = ToThread()
	}
	//ToThread()

	isUser := false
	if isThread == true {
		isUser = ToUser()
		//isUser = true
	}
	//导入用户
	//ToUser()

	isUser,msg := UpdateUser()
	log.Println(msg)

	if isUser == true {
		log.Println("===\n Data Import Success! \n===")
	}

	//ToMyThreads()  //已导入主帖后，导入帖子归属

}
*/