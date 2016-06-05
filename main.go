package main

import (

	"github.com/skiy/discuzx-xiuno/app"
	//"./app"
	"fmt"
)

func main() {
	fmt.Println(`
::: 欢迎使用Discuz!X3.2 To XiunoBBS3.0 转换程序
::: 作者: Skiychan <dev@skiy.net>
::: 网站: https://www.zzzzy.com
::: QQ : 1005043848
::: 本程序已开源托管至GitHub: https://github.com/skiy/discuzx-xiuno
`)

	//进入主程序
	app.Init()
}