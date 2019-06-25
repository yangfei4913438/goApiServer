package dbs

import (
	"github.com/astaxie/beego"
	"goApiServer/tools"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type HTTP struct {
	BaseUrl    string
	BaseHeader map[string]string
}

var Request *HTTP

// 初始化 HTTP 类型
func initHttp() {
	// 初始化 HTTP 对象
	Request = &HTTP{
		BaseUrl: beego.AppConfig.String("base_url"),
	}
	beego.Info("Initialize the Http module successfully!")
}

type httpResult struct {
	StatusCode int    `json:"status_code"`
	Body       []byte `json:"body"`
}

func (api *HTTP) client(method string, url string, header map[string]string, body io.Reader) (*httpResult, error) {
	// 生成请求对象
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		beego.Error("生成请求对象出错:", err)
		return nil, err
	}

	// 添加默认头部信息
	if len(api.BaseHeader) > 0 {
		for k, v := range api.BaseHeader {
			request.Header.Add(k, v)
		}
	}

	// 添加自定义头部信息
	if len(header) > 0 {
		for k, v := range header {
			request.Header.Add(k, v)
		}
	}

	// 定义客户端
	client := http.Client{
		Timeout: tools.OneSecond * 5, // 设置超时
		Transport: &http.Transport{
			DisableKeepAlives: true, // http设置为短连接请求
		},
	}

	// 发起 http 请求
	response, err := client.Do(request)
	// 正确的关闭姿势，因为 response 可能是 nil
	if response != nil {
		defer response.Body.Close()
	}
	// 错误处理
	if err != nil {
		beego.Error("发起 http 请求出错:", err, response)
		return nil, err
	}

	// 生成返回数据
	send, err := ioutil.ReadAll(response.Body)
	if err != nil {
		beego.Error("读取返回数据出错", err)
		return nil, err
	}

	// 返回结果
	return &httpResult{response.StatusCode, send}, nil
}

func (api *HTTP) Get(url string, header map[string]string, params map[string]string) (*httpResult, error) {
	// 生成 URL
	url = api.BaseUrl + url
	if len(params) > 0 {
		index := 0
		for k, v := range params {
			if index == 0 {
				url += "?" + k + "=" + v
			} else {
				url += "&" + k + "=" + v
			}
			index++
		}
	}

	return api.client("GET", url, header, nil)
}

func (api *HTTP) Post(url string, header map[string]string, data []byte) (*httpResult, error) {
	// 生成 URL
	url = api.BaseUrl + url

	return api.client("POST", url, header, strings.NewReader(string(data)))
}

func (api *HTTP) Put(url string, header map[string]string, data []byte) (*httpResult, error) {
	// 生成 URL
	url = api.BaseUrl + url

	return api.client("PUT", url, header, strings.NewReader(string(data)))
}

// 如果是从 data 参数中定义删除参数，使用这个方法
func (api *HTTP) DeleteFromData(url string, header map[string]string, data []byte) (*httpResult, error) {
	// 生成 URL
	url = api.BaseUrl + url

	return api.client("DELETE", url, header, strings.NewReader(string(data)))
}

// 如果是从 url 参数中定义删除参数，使用这个方法
func (api *HTTP) DeleteFromParams(url string, header map[string]string, params map[string]string) (*httpResult, error) {
	// 生成 URL
	url = api.BaseUrl + url
	if len(params) > 0 {
		index := 0
		for k, v := range params {
			if index == 0 {
				url += "?" + k + "=" + v
			} else {
				url += "&" + k + "=" + v
			}
			index++
		}
	}

	return api.client("DELETE", url, header, nil)
}
