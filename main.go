package main

import (
	"fmt"

	"adc-sdk-go/initialize"
	"adc-sdk-go/request"
)

func adcApiT() {
	req := request.RequestAdc{
		Logout:  false,
		Authkey: "433ad43189a344c5e0f85c5a4457dd",
		Print:   true,
	}
	// 初始化配置
	req.GetConf()
	initialize.InitRequestsPool(&req) // 初始化HttpClient
	// 获取session
	req.GetAuthkey()
	// v2 api 测试

	// node
	//列表
	// req.SlbNodeList()
	// // common分区
	// req.SlbNodeListCommon()
	//获取

	// req.SlbNodeGet("3.1.1.1")
	// 添加
	// node := response.Node{Name: "3.1.1.1", Host: "3.1.1.1"}
	// req.SlbNodeAdd(node)
	// 更新
	// node := response.Node{Name: "3.1.1.1", Host: "3.1.1.1", Conn_limit: 555}
	// req.SlbNodeEdit(node)
	// 删除
	// req.SlbNodeDel("3.1.1.1")
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ //

	// ----------------------------------------- //
	fmt.Println("Authkey  ===>  ", req.Authkey)
	fmt.Println("ResultPtr  ===>  ", req.ResultPtr)
	fmt.Println("req  ===>  ", req)
	req.LogoutF()
}

func main() {
	adcApiT()
}
