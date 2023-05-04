package response

// 请求ADC 错误返回值
type RespResult struct {
	Result  string `json:"result"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type AuthkeyResult struct {
	Authkey string `json:"authkey"`
	RespResult
}
