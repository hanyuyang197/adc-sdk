package request

import (
	"adc-sdk-go/response"
	"adc-sdk-go/utils"
	"fmt"
)

// 节点(node)列表
func (req *RequestAdc) SlbNodeList() (*RequestAdc, error) {
	req.slbNodeList(false)
	return req, nil
}

// 单独查common分区
func (req *RequestAdc) SlbNodeListCommon() (*RequestAdc, error) {
	req.slbNodeList(true)
	return req, nil
}

// list
func (req *RequestAdc) slbNodeList(common bool) (*RequestAdc, error) {
	// 将err参数置空
	req.Err = nil
	req.Errcode = 0
	action := "slb.node.list"
	if common {
		action += ".withcommon"
	}
	result := []response.Node{}
	req.URL = fmt.Sprintf("https://%s:%d/%s/?action=%s&authkey=%s", req.BaseInfo.Ip, req.BaseInfo.Port, req.BaseInfo.AdcVersion, action, req.Authkey)
	req.ResultPtr = &result

	req.Get()
	req.ErrLog(1, fmt.Errorf("测试打印"))
	// if req.Errcode != 0 {
	// 	return req, fmt.Errorf(req.Err.Error())
	// }
	return req, nil
}

// 获取
func (req *RequestAdc) SlbNodeGet(name string) error {
	// 将err参数置空
	req.Err = nil
	req.Errcode = 0
	action := "slb.node.get"
	result := response.RespResult{}

	bodyData := response.NodeRequest{Name: name}
	req.URL = fmt.Sprintf("https://%s:%d/%s/?action=%s&authkey=%s", req.BaseInfo.Ip, req.BaseInfo.Port, req.BaseInfo.AdcVersion, action, req.Authkey)

	req.Data = bodyData
	req.ResultPtr = &result
	req.Print = true
	req.Post()
	if req.Errcode != 0 {
		fmt.Println("Errcode error  ===>  ", result.Errmsg)
	}
	if result.Result != "success" {
		fmt.Println("Result error  ===>  ", result.Errmsg)

	}

	return req.Err
}

// 添加
func (req *RequestAdc) SlbNodeAdd(node response.Node) error {
	// 将err参数置空
	req.Err = nil
	req.Errcode = 0
	action := "slb.node.add"
	result := response.RespResult{}

	nodeBody, _ := utils.ToMap(node, "json")
	bodyData := response.NodeRequest{}
	bodyData.Node = nodeBody
	fmt.Printf("(node)  bodyData  ===>  %+v", nodeBody)
	req.URL = fmt.Sprintf("https://%s:%d/%s/?action=%s&authkey=%s", req.BaseInfo.Ip, req.BaseInfo.Port, req.BaseInfo.AdcVersion, action, req.Authkey)

	req.Data = bodyData
	req.ResultPtr = &result
	req.Print = true

	req.Post()

	if req.Errcode != 0 {
		fmt.Println("Errcode error  ===>  ", result.Errmsg)
	}
	if result.Result != "success" {
		fmt.Println("Result error  ===>  ", result.Errmsg)

	}

	return req.Err
}

// 更新
func (req *RequestAdc) SlbNodeEdit(node response.Node) error {
	// 将err参数置空
	req.Err = nil
	req.Errcode = 0
	action := "slb.node.edit"
	result := response.RespResult{}

	nodeBody, _ := utils.ToMap(node, "json")
	bodyData := response.NodeRequest{}
	bodyData.Node = nodeBody

	fmt.Printf("(node)  bodyData  ===>  %+v", nodeBody)
	req.URL = fmt.Sprintf("https://%s:%d/%s/?action=%s&authkey=%s", req.BaseInfo.Ip, req.BaseInfo.Port, req.BaseInfo.AdcVersion, action, req.Authkey)

	req.Data = bodyData
	req.ResultPtr = &result
	req.Print = true

	req.Post()

	if req.Errcode != 0 {
		fmt.Println("Errcode error  ===>  ", result.Errmsg)
	}
	if result.Result != "success" {
		fmt.Println("Result error  ===>  ", result.Errmsg)

	}

	return req.Err
}

// 删除
func (req *RequestAdc) SlbNodeDel(name string) error {
	// 将err参数置空
	req.Err = nil
	req.Errcode = 0
	action := "slb.node.del"
	result := response.RespResult{}

	bodyData := response.NodeRequest{Name: name}
	req.URL = fmt.Sprintf("https://%s:%d/%s/?action=%s&authkey=%s", req.BaseInfo.Ip, req.BaseInfo.Port, req.BaseInfo.AdcVersion, action, req.Authkey)

	req.Data = bodyData
	req.ResultPtr = &result
	req.Print = true
	req.Post()
	if req.Errcode != 0 {
		fmt.Println("Errcode error  ===>  ", result.Errmsg)
	}
	if result.Result != "success" {
		fmt.Println("Result error  ===>  ", result.Errmsg)

	}

	return req.Err
}
