package request

import (
	"adc-sdk-go/global"
	"adc-sdk-go/response"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Errcode uint

const (
	Fail Errcode = iota + 1
	ParseErr
	URLParseErr
	DataParseErr
	FileParseErr
	FileReadErr
	NotDevErr
	AuthkeyErr
	LoginErr
	ResultTypeErr
	RespReadErr
	JSONParseErr
	StatusCodeErr
	SessionErr
	KeepAliveErr
	GetFileErr

	OtherErr Errcode = 30
)

type HttpRequest struct {
	LogLevel              string `mapstructure:"log-level" json:"log-level" yaml:"log-level"`                                           // adc 日志等级
	Timeout               int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`                                                 // 超时时间
	MaxIdleConnsPerHost   int    `mapstructure:"max-idle-per-host" json:"max-idle-per-host" yaml:"max-idle-per-host"`                   // 最大空闲连接数/Host(Host=IP+Port)
	MaxConnsPerHost       int    `mapstructure:"max-conns-per-host" json:"max-conns-per-host" yaml:"max-conns-per-host"`                // 最大连接数/Host(Host=IP+Port)
	IdleConnTimeout       int    `mapstructure:"idle-conn-timeout" json:"idle-conn-timeout" yaml:"idle-conn-timeout"`                   // 空闲连接超时时间
	TLSHandshakeTimeout   int    `mapstructure:"tls-handshake-timeout" json:"tls-handshake-timeout" yaml:"tls-handshake-timeout"`       // tls 握手 超时时间
	ExpectContinueTimeout int    `mapstructure:"expect-continue-timeout" json:"expect-continue-timeout" yaml:"expect-continue-timeout"` //
	InsecureSkipVerify    bool   `mapstructure:"insecure-skip-verify" json:"insecure-skip-verify" yaml:"insecure-skip-verify"`          // 跳过证书验证
}

type RequestAdc struct {
	URL       string
	BaseInfo  BaseInfo // 自动填充IP port
	IP        string
	Port      int
	Authkey   string
	Params    map[string]string
	Data      interface{}
	Files     []File
	ResultPtr interface{}    // 请使用指针, 结果自动填充
	Resp      *http.Response // 调用方法后自动填充
	Err       error          // 调用方法后自动填充
	Errcode   uint           // 调用方法后自动填充
	Mehtod    string         // 调用方法后自动填充
	Logout    bool           // 登出
	Print     bool           // 打印
	GetFile   *os.File       // 获取文件, 指针为 *os.File 通过io.Copy方法
}

type File struct {
	FieldName string // Post请求Form中的 文件对应字段
	FileName  string // 文件路径(存在Content时,可写文件名)
	Content   []byte // 文件写入内存中的内容
}

//解析yml文件
type BaseInfo struct {
	Port        int         `yaml:"port"`
	Ip          string      `yaml:"ip"`
	AdcVersion  string      `yaml:"adcversion"`
	LoginData   LoginData   `yaml:"login"`
	HttpRequest HttpRequest `yaml:"request"`
}

type LoginData struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Timeout  string `yaml:"timeout"`
}

func (req *RequestAdc) GetConf() {
	pwd, _ := os.Getwd()
	yamlFile, err := ioutil.ReadFile(path.Join(pwd, "conf.yml"))
	// fmt.Println(path.Join(pwd, "conf.yml"))
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &req.BaseInfo)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func (req *RequestAdc) GetAuthkey() (authkey string, err error) {
	// 尝试先logout
	// return
	if !req.Logout {
		req.LogoutF()
	}

	// 登录
	req, result, err := req.Login()
	if err != nil {
		return "", err
	}

	switch result.Errcode {
	case 1206462297:
		req.ErrLog(9, fmt.Errorf("用户名密码错误"))
		return "", fmt.Errorf("用户名密码错误")
	case 1492068187:
		req.ErrLog(14, fmt.Errorf("Session 上限"))
		return "", fmt.Errorf("Session 上限")
	case 1452616536:
		req.ErrLog(30, fmt.Errorf("用户锁定"))
		return "", fmt.Errorf("用户锁定")
	case 0:
		req.Authkey = result.Authkey
		// 更新
		return result.Authkey, nil
	default:
		req.ErrLog(30, fmt.Errorf("登录失败"))
		return "", fmt.Errorf("登录失败")
	}
}

// ADC登录
func (req *RequestAdc) Login() (*RequestAdc, *response.AuthkeyResult, error) {
	action := "login"
	result := response.AuthkeyResult{}
	req.URL = fmt.Sprintf("https://%s:%d/%s/?action=%s&username=%s&password=%s", req.BaseInfo.Ip, req.BaseInfo.Port, req.BaseInfo.AdcVersion, action, req.BaseInfo.LoginData.Username, req.BaseInfo.LoginData.Password)
	req.ResultPtr = &result
	// req = request.RequestAdc{
	// 	URL:       fmt.Sprintf("https://%s:%d/%s/?action=%s&username=%s&password=%s", req.BaseInfo.IP, req.BaseInfo.Port, action, action, req.LoginData.Username, req.LoginData.Password),
	// 	ResultPtr: &result,
	// }
	// fmt.Println("00000   ===>  ", req)

	req.Get()

	if req.Err != nil {
		req.ErrLog(7, fmt.Errorf("请确认IP是否可用"))
		return req, &result, fmt.Errorf("请确认IP是否可用")
	}
	return req, &result, nil
}

// 退出登录
func (req *RequestAdc) LogoutF() (err error) {
	result := response.RespResult{}
	req.URL = fmt.Sprintf("https://%s:%d/adcapi/v2.0/?authkey=%s&action=logout", req.BaseInfo.Ip, req.BaseInfo.Port, req.Authkey)
	req.ResultPtr = &result
	req.Logout = true
	// req = RequestAdc{
	// 	URL:       fmt.Sprintf("https://%s:%d/adcapi/v2.0/?authkey=%s&action=logout", IP, Port, authkey),
	// 	ResultPtr: &result,
	// 	Logout:    true,
	// }
	req.Get()

	if req.Err != nil {
		return req.Err
	}
	req.Authkey = ""
	return nil
}

// 打印错误日志 并赋值到Err Errcode
func (req *RequestAdc) ErrLog(Errcode uint, err error) {
	ErrMsg := fmt.Sprintf("Method: %s, Url: %s, Errcode: %d, Error: %s", req.Mehtod, req.URL, Errcode, err.Error())
	fmt.Println(ErrMsg)
	req.Err = err
	req.Errcode = Errcode
	return
}

func (req *RequestAdc) Log(message string) {
	InfoMsg := fmt.Sprintf("Method: %s, url: %s, %s", req.Mehtod, req.URL, message)
	fmt.Println(InfoMsg)
	return
}

// 默认传递URL即可
// GET请求
func (req *RequestAdc) Get() {
	var err error
	req.Mehtod = "GET"
	req.ParseURL()
	if req.Err != nil {
		return
	}
	req.Resp, err = global.HTTPClient.Get(req.URL)
	if err != nil {
		req.ErrLog(1, fmt.Errorf("Get请求第一次 %s", err.Error()))
		return
	}
	fmt.Println(" req.URL ===>  ", req.URL)

	req.ParseRespResult()
	// fmt.Println("--------------", req.ResultPtr)
	// 第二次重试
	if !req.Logout && req.Retry() {
		// 如果重试 将Err重置为空
		req.Err = nil
		req.Errcode = 0
		req.ParseURL()

		if req.Err != nil {
			return
		}
		req.Resp, err = global.HTTPClient.Get(req.URL)
		if err != nil {
			req.ErrLog(1, fmt.Errorf("Get请求第二次 %s", err.Error()))
			return
		}
		req.ParseRespResult()
		// fmt.Println("===============", req.ResultPtr)
	}
	if req.Err == nil {
		req.Log("请求成功")
	}
	return
}

func (req *RequestAdc) Post() {
	var err error
	req.Mehtod = "POST"
	req.ParseURL()
	if req.Err != nil {
		return
	}
	contentType := "application/json;charset=UTF8"
	//add post body
	body := &bytes.Buffer{}
	if req.Data != nil {
		var DataJson []byte
		DataJson, err = json.Marshal(req.Data)
		if err != nil {
			req.ErrLog(4, err)
			return
		}
		body = bytes.NewBuffer(DataJson)
	}

	// 填充文件
	if req.Files != nil {
		writer := multipart.NewWriter(body)
		for _, file := range req.Files {
			fileWriter, err := writer.CreateFormFile(file.FieldName, file.FileName)
			if err != nil {
				req.ErrLog(5, err)
				return
			}
			// 如果文件内容在内存中
			if file.Content != nil {
				body.Write(file.Content)
			} else {
				// 读取 文件
				f, err := os.Open(file.FileName)
				if err != nil {
					req.ErrLog(6, err)
				}
				io.Copy(fileWriter, f)
				f.Close()
			}
		}
		writer.Close()
		// 如果是文件类型需要添加 multipart/form-data; boundary=XXXXXXXXXXXXXXX
		contentType = writer.FormDataContentType()
	}

	req.Resp, err = global.HTTPClient.Post(req.URL, contentType, body)
	if err != nil {
		req.ErrLog(1, fmt.Errorf("Post请求第一次 %s", err.Error()))
		return
	}
	req.ParseRespResult()

	// 第二次重试
	if req.Retry() {
		// 如果重试 将Err重置为空
		req.Err = nil
		req.Errcode = 0
		req.ParseURL()
		if req.Err != nil {
			return
		}
		req.Resp, err = global.HTTPClient.Post(req.URL, contentType, body)
		if err != nil {
			req.ErrLog(1, fmt.Errorf("Post请求第二次 %s", err.Error()))
			return
		}
		req.ParseRespResult()
	}
	if req.Err == nil {
		req.Log("请求成功")
	}
	return
}

// 解析Url
func (req *RequestAdc) ParseURL() {
	u, err := url.Parse(req.URL)
	if err != nil {
		req.ErrLog(3, err)
		return
	}
	// fmt.Println("-----------", u.Host)
	if u.Host == "127.0.0.1" || u.Host == "" {
		// if req.BaseInfo == nil {
		// 	req.ErrLog(3, err)
		// 	return
		// }
		IP, Port := req.BaseInfo.Ip, req.BaseInfo.Port
		if req.IP != "" && req.Port != 0 {
			u.Host = fmt.Sprintf("%s:%d", req.IP, req.Port)
			u.Scheme = "https"
		} else if IP != "" && Port != 0 {
			u.Host = fmt.Sprintf("%s:%d", IP, Port)
			u.Scheme = "https"
		} else {
			req.ErrLog(3, err)
			return
		}
	} else {
		host_arr := strings.Split(u.Host, ":")
		req.IP = host_arr[0]
		if len(host_arr) == 2 {
			port, err := strconv.Atoi(host_arr[1])
			if err == nil {
				req.Port = port
				u.Scheme = "https"
			}
		} else {
			u.Scheme = "https"
			req.Port = 443
		}
	}
	// 如果Parms有值填充params
	params := u.Query()
	for key, value := range req.Params {
		params.Set(key, value)
	}
	// 需要authkey时, 用到设备
	vs := params["authkey"]
	if len(vs) != 0 {
		// authkey 优先级 Params > Authkey > Dev
		// 需要authkey 且需要填充
		// authkey 优先级高
		if req.Authkey != "" {
			params.Set("authkey", req.Authkey)
		} else {

			//  authkey为空  直接重试

			authkey, err := req.GetAuthkey()
			if err != nil {
				req.ErrLog(8, fmt.Errorf("authkey获取失败"))
				return
			}

			req.Authkey = authkey
			params.Set("authkey", authkey)
		}
	}

	u.RawQuery = params.Encode()
	// 转成url
	req.URL = u.String()
	return
}

func IsPtr(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Ptr
}

// 解析返回结果
func (req *RequestAdc) ParseRespResult() {
	if !IsPtr(req.ResultPtr) {
		req.ErrLog(10, fmt.Errorf("不是指针类型"))
		return
	}

	switch reflect.TypeOf(req.ResultPtr) {
	case nil:
		// 如果结构体为空, 默认赋值
		var result map[string]interface{}
		req.ResultPtr = &result
	}
	// 结束时关闭Body
	defer req.Resp.Body.Close()

	if req.Resp.StatusCode == 200 {
		if req.GetFile != nil {
			_, err := io.Copy(req.GetFile, req.Resp.Body)
			if err != nil {
				req.ErrLog(16, err)
			}
			return
		}
		// 解析resp body 转json
		body, err := ioutil.ReadAll(req.Resp.Body)
		if err != nil {
			req.ErrLog(11, err)
			return
		}
		// 打印返回结果
		if req.Print {
			fmt.Printf("\nURL: %v \nBody: %v \n", req.URL, string(body))
		}

		err = json.Unmarshal(body, req.ResultPtr)
		if err != nil {
			req.ErrLog(12, err)
			return
		}
	} else {
		req.ErrLog(13, fmt.Errorf("状态码为: %d", req.Resp.StatusCode))
		return
	}
	return
}

// Invalid authkey.进行重试
func (req *RequestAdc) Retry() bool {
	// 重试机制 -- 判断结构体中, 存在Errcode 字段, 并且不为Nil 值  如果是Invalid authkey.进行重试
	// reflect.ValueOf 需传递结构体
	// fmt.Println(reflect.ValueOf(req.ResultPtr).Elem())
	// fmt.Println(reflect.ValueOf(req.ResultPtr).Elem().Type())
	switch reflect.ValueOf(req.ResultPtr).Elem().Type().Kind() {
	case reflect.Struct:
		if Errcode := reflect.ValueOf(req.ResultPtr).Elem().FieldByName("Errcode"); Errcode.IsValid() {
			// Errcode == 1492067499 (Invalid authkey.)
			if Errcode.Int() == 1492067499 {
				authkey, err := req.GetAuthkey()
				if err != nil {
					return false
				}
				if req.Params != nil {
					req.Params["authkey"] = authkey
				}
				req.Authkey = authkey
				return true
			}
		}
	default:
		return false
	}

	return false
}

// Get请求 进一步封装 result传递指针
func GetAdcAPI(Url string, resultPtr interface{}) *RequestAdc {
	req := RequestAdc{
		URL:       Url,
		ResultPtr: resultPtr,
		// Print:     true,
	}
	req.Get()
	return &req
}
