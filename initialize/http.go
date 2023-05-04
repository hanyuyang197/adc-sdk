package initialize

import (
	"adc-sdk-go/global"
	"adc-sdk-go/request"
	"crypto/tls"
	"net/http"
	"time"
)

func InitRequestsPool(req *request.RequestAdc) *http.Client {
	HTTPClient := &http.Client{
		// return &http.Client{
		Timeout: time.Duration(req.BaseInfo.HttpRequest.Timeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   req.BaseInfo.HttpRequest.MaxIdleConnsPerHost,
			MaxConnsPerHost:       req.BaseInfo.HttpRequest.MaxConnsPerHost,
			IdleConnTimeout:       time.Duration(req.BaseInfo.HttpRequest.IdleConnTimeout) * time.Second,
			TLSHandshakeTimeout:   time.Duration(req.BaseInfo.HttpRequest.TLSHandshakeTimeout) * time.Second,
			ExpectContinueTimeout: time.Duration(req.BaseInfo.HttpRequest.ExpectContinueTimeout) * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: req.BaseInfo.HttpRequest.InsecureSkipVerify},
		},
	}
	global.HTTPClient = HTTPClient
	return HTTPClient

}
