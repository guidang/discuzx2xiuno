package app

import (
	"database/sql"
	"log"

	"bufio"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

//数据库初始化
var (
	OldDB,
	NewDB *sql.DB
	ResetPost = false
	Exts      Extenter
	ClearTB   = true //是否先清理表
	MergeUser = true //是否合并用户
)

type Extenter struct {
	AdminUid     string //管理员 uid
	UpdateFromDz bool   //是否从旧版的 dz 升级到 dx 的
}

/**
初始化程序
*/
func Init() {
	fmt.Println("::: 正在进入app主程序...")
	//OldDB, NewDB = connDB()

	var err error
	OldDB, NewDB, err = InputDatabase()
	if err != nil {
		log.Fatalln(err.Error())
	}

	//设置 uid 和 判断升级版本
	inputOtherInfo()

	var msg string

	//转换回复帖子
	msg = ToPost()
	fmt.Println(msg, "\r\n")

	//转换主题帖子
	msg = ToThread()
	fmt.Println(msg, "\r\n")

	//转换版块
	msg = ToForum()
	fmt.Println(msg, "\r\n")

	msg = ToUser()
	fmt.Println(msg, "\r\n")

	msg = updateAdminUser()
	fmt.Println(msg, "\r\n")

	/* 更新全部用户帖子数量 */
	if ResetPost {
		msg := doUserPosts()
		fmt.Println(msg, "\r\n")
	}

	for {
		fmt.Println("\r\n::: 转换完成, 按 \"回车键\" 退出程序...")
		r := bufio.NewReader(os.Stdin)
		b, _, _ := r.ReadLine()
		inputLen := len(b)
		if inputLen == 0 {
			os.Exit(0)
		}
	}
}

/**
	手动连接
  	连接新旧数据库
*/
func connDB() (oldDB, newDB *sql.DB) {
	fmt.Println("::: 正在输入数据库信息...")

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

	var err error
	oldDB, err = connectMysql(old)
	if err != nil {
		log.Fatalf("\r\n%s connect err: %s\r\n", "Discuz", err.Error())
	}

	newDB, err = connectMysql(new)
	if err != nil {
		log.Fatalf("\r\n%s connect err: %s\r\n", "Xiuno", err.Error())
	}

	return
}

/**
配置数据库信息
*/
func InputDatabase() (oldDb, newDb *sql.DB, err error) {
	fmt.Println(`
::: 正在输入数据库信息...
::: 执行过程中按"Q", 再按"回车键"退出本程序...
`)

	r := bufio.NewReader(os.Stdin)

	oldhost := &Hostinfo{}
	newhost := &Hostinfo{}

	o_flag := "Discuz!X"
	n_flag := "XiunoBBS"

	inputDataInfo(r, oldhost, o_flag)
	inputDataInfo(r, newhost, n_flag)

	fmt.Printf("\r\n%s: %s \r\n%s: %s \r\n\r\n", o_flag, oldhost, n_flag, newhost)

	oldDb, err = connectMysql(oldhost)
	if err != nil {
		log.Fatalf("\r\n%s connect err: %s\r\n", o_flag, err.Error())
	}

	newDb, err = connectMysql(newhost)
	if err != nil {
		log.Fatalf("\r\n%s connect err: %s\r\n", n_flag, err.Error())
	}

	return
}

/**
清理数据表
*/
func ClearTable(tbname string) error {
	//若不清理表则
	if !ClearTB {
		return nil
	}

	fmt.Printf(":::正在清理 %s 表\r\n", tbname)

	clearSQL := fmt.Sprintf("TRUNCATE TABLE %s", tbname)

	_, err := NewDB.Exec(clearSQL)
	if err != nil {
		fmt.Printf(":::清理 %s 表失败: %s\r\n", tbname, err.Error())
	}

	return err
}

/**
输入数据库信息
*/
func inputDataInfo(r *bufio.Reader, h *Hostinfo, t string) {
	fmt.Printf("\r\n正在配置 %s 的数据库信息.....", t)

	var flag int
	for {
		switch flag {
		case 0:
			fmt.Printf("\r\n配置 %s 的host(默认为 127.0.0.1): ", t)
			s := inputData(r)

			if s == "" {
				s = "127.0.0.1"
			}
			h.DBHost = s
			flag++

		case 1:
			fmt.Printf("\r\n配置 %s 的数据库用户(默认为 root):", t)
			s := inputData(r)
			if s == "" {
				s = "root"
			}
			h.DBUser = s
			flag++

		case 2:
			fmt.Printf("\r\n配置 %s 的数据库密码:", t)
			s := inputData(r)
			h.DBPassword = s
			flag++

		case 3:
			fmt.Printf("\r\n配置 %s 的数据库名:", t)
			s := inputData(r)
			if s != "" {
				h.DBName = s
				flag++
			}

		case 4:
			fmt.Printf("\r\n配置 %s 的数据库端口(默认为3306):", t)
			s := inputData(r)
			if s == "" {
				s = "3306"
			}
			h.DBPort = s
			flag++

		case 5:
			fmt.Printf("\r\n配置 %s 的数据库编码(默认为utf8):", t)
			s := inputData(r)
			if s == "" {
				s = "utf8"
			}
			h.DBChar = s
			flag = 99
			break
		}

		if flag == 99 {
			break
		}
	}
}

/**
输入扩展信息
*/
func inputOtherInfo() {
	r := bufio.NewReader(os.Stdin)

	//配置 uid
	fmt.Printf("\r\n配置管理员的uid(默认为1):")
	s := inputData(r)
	if s == "" {
		s = "1"
	}
	Exts.AdminUid = s

	//是否从 dz 升级上来的
	fmt.Printf(`
数据是否从老版本Discuz!系列升级到Discuz!X系列?
请输入 "Y" or "N"(默认为否N):
`)

	s = inputData(r)
	if s == "Y" || s == "y" {
		Exts.UpdateFromDz = true
	} else {
		Exts.UpdateFromDz = false
	}
}

/**
键盘输入数据
*/
func inputData(r *bufio.Reader) string {
	b, _, _ := r.ReadLine()
	s := string(b)

	if s == "q" || s == "Q" {
		os.Exit(0)
	}

	return s
}
