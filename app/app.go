package app

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"bufio"
	"os"
	"fmt"
)

//数据库初始化
var (
	OldDB,
	NewDB *sql.DB
	ClearTB   = true //是否先清理表
	MergeUser = true //是否合并用户
	ResetPost = false
)

func Init() {
	log.Println(":::正在进入app主程序:::")
	//OldDB, NewDB = connDB()
	OldDB, NewDB = InputDatabase()

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
		DBUser:     "root",
		DBPassword: "123456",
		DBname:     "gxvtc",
	}

	new := &Hostinfo{
		DBUser:     "root",
		DBPassword: "123456",
		DBname:     "xiuno",
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

func InputDatabase() (oldDb, newDb *sql.DB) {
	log.Println(":::正在输入数据库:::")

	r := bufio.NewReader(os.Stdin)

	oldhost := &Hostinfo{}
	newhost := &Hostinfo{}

	for {
		if oldhost.DBChar == "" {
			fmt.Println("正在配置discuzx的数据库.....")
			if oldhost.DBHost == "" {
				fmt.Print("配置discuzx的host(默认为127.0.0.1): ")
				b, _, _ := r.ReadLine()
				line := string(b)

				if line == "" {
					line = "127.0.0.1"
				}

				oldhost.DBHost = line
			} else if oldhost.DBUser == "" {
				fmt.Print("配置discuzx的数据库用户: ")
				b, _, _ := r.ReadLine()
				line := string(b)

				oldhost.DBUser = line
			} else if oldhost.DBPassword == "" {
				fmt.Print("配置discuzx的数据库密码(不能为空): ")
				b, _, _ := r.ReadLine()
				line := string(b)

				oldhost.DBPassword = line
			} else if oldhost.DBname == "" {
				fmt.Print("配置discuzx的数据库名: ")
				b, _, _ := r.ReadLine()
				line := string(b)

				oldhost.DBname = line
			} else if oldhost.DBPort == "" {
				fmt.Print("配置discuzx的数据库端口(默认为3306): ")

				b, _, _ := r.ReadLine()
				line := string(b)

				if line == "" {
					line = "3306"
				}

				oldhost.DBPort = line
			} else if oldhost.DBChar == "" {
				fmt.Print("配置discuzx的数据库编码(默认为utf8): ")

				b, _, _ := r.ReadLine()
				line := string(b)

				if line == "" {
					line = "utf8"
				}

				oldhost.DBChar = line
			}
		} else if newhost.DBChar == "" {
			fmt.Println("正在配置xiuno的数据库.....")
			if newhost.DBHost == "" {
				fmt.Print("配置xiuno的host(默认为127.0.0.1): ")
				b, _, _ := r.ReadLine()
				line := string(b)

				if line == "" {
					line = "127.0.0.1"
				}

				newhost.DBHost = line
			} else if newhost.DBUser == "" {
				fmt.Print("配置xiuno的数据库用户: ")
				b, _, _ := r.ReadLine()
				line := string(b)

				newhost.DBUser = line
			} else if newhost.DBPassword == "" {
				fmt.Print("配置xiuno的数据库密码(不能为空): ")
				b, _, _ := r.ReadLine()
				line := string(b)

				newhost.DBPassword = line
			} else if newhost.DBname == "" {
				fmt.Print("配置xiuno的数据库名: ")
				b, _, _ := r.ReadLine()
				line := string(b)

				newhost.DBname = line
			} else if newhost.DBPort == "" {
				fmt.Print("配置xiuno的数据库端口(默认为3306): ")

				b, _, _ := r.ReadLine()
				line := string(b)

				if line == "" {
					line = "3306"
				}

				newhost.DBPort = line
			} else if newhost.DBChar == "" {
				fmt.Print("配置xiuno的数据库编码(默认为utf8): ")

				b, _, _ := r.ReadLine()
				line := string(b)

				if line == "" {
					line = "utf8"
				}

				newhost.DBChar = line
			}
		} else {
			break
		}

	}

	fmt.Println(oldhost, newhost)

	oldDb, err := connectMysql(oldhost)
	if err != nil {
		log.Println("old db connect err: " + err.Error())
	}

	newDb, err = connectMysql(newhost)
	if err != nil {
		log.Println("new db connect err: " + err.Error())
	}

	return
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
