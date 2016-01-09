package app
import "fmt"

//数据库初始化
var OldDB,NewDB = CreateDB()

func Init3() {
	ToThread()
}

func Init()  {

	//导入版块
	//ToForum()

	//开关
	isRun := true

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
		//isUser = ToUser()
		isUser = true
	}
	//导入用户
	//ToUser()

	if isUser == true {
		fmt.Println("===\n Data Import Success! \n===")
	}
}