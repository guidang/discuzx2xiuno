package main

import (
	"log"

	//"github.com/skiy/discuzx-xiuno/app"
	"./app"
)

func main() {
	log.Println(`
	:::欢迎使用DiscuzX3.2 To XiunoBBS3.0 转换程序:::
	:::作者: Skiychan <dev@skiy.net>
	:::网站: www.zzzzy.com
	::: QQ : 1005043848
	:::本程序已开源托管至GitHub: https://github.com/skiy/discuzx-xiuno
	`)

	//转换版块
	app.Init()
}