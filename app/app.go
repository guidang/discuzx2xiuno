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
	//AdminUid = 1  //管理员 uid
	AdminUid string
)

func Init() {
	log.Println(":::正在进入app主程序:::")
	//OldDB, NewDB = connDB()

	var err error
	OldDB, NewDB, err = InputDatabase()
	if err != nil {
		log.Fatalln(err)
	}

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

	for {
		fmt.Println("\r\n::: 转换完成, 按 \"回车键\" 退出程序...")
		r := bufio.NewReader(os.Stdin)
		b, _, _ := r.ReadLine()
		inputLen := len(b)
		if inputLen == 0 {
			break
		}
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
		DBName:     "gxvtc",
	}

	new := &Hostinfo{
		DBUser:     "root",
		DBPassword: "123456",
		DBName:     "xiuno",
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

func InputDatabase() (oldDb, newDb *sql.DB, err error) {
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
			} else if oldhost.DBName == "" {
				fmt.Print("配置discuzx的数据库名: ")
				b, _, _ := r.ReadLine()
				line := string(b)

				oldhost.DBName = line
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
			} else if newhost.DBName == "" {
				fmt.Print("配置xiuno的数据库名: ")
				b, _, _ := r.ReadLine()
				line := string(b)

				newhost.DBName = line
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
		} else if(AdminUid == "") {
			fmt.Print("配置xiuno的管理员的uid(默认为uid为1): ")

			b, _, _ := r.ReadLine()
			line := string(b)

			if line == "" {
				line = "1"
			}

			AdminUid = line
		} else {
			break
		}

	}

	fmt.Println("\r\nDiscuz!X数据库:",oldhost, "\r\nXiunoBBS数据库:",newhost,"\r\n")

	oldDb, err = connectMysql(oldhost)
	if err != nil {
		log.Println("old db connect err: " + err.Error())
		return
	}

	newDb, err = connectMysql(newhost)
	if err != nil {
		log.Println("new db connect err: " + err.Error())
		return
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
