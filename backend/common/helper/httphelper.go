package helper

import (
	"bytes"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ApiListData struct {
	CurPage    int         `json:"curPage"`
	LastPage   int         `json:"lastPage"`
	List       interface{} `json:"list"`
	PageSize   int         `json:"pageSize"`
	TotalCount int         `json:"totalCount"`
}

func HttpPost(httpUrl string, params string) (int, string) {
	resp, err := http.Post(httpUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(params))
	if err != nil {
		logx.Error(err)
		return http.StatusInternalServerError, err.Error()
	}

	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	statusCode := resp.StatusCode
	if readErr != nil {
		logx.Error(readErr)
		return http.StatusInternalServerError, readErr.Error()
	}

	return statusCode, string(body)
}

func HttpPostForm(httpUrl string, params map[string][]string) (int, string) {
	resp, err := http.PostForm(httpUrl,
		params)
	if err != nil {
		logx.Error(err)
		return http.StatusInternalServerError, err.Error()
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	statusCode := resp.StatusCode
	if readErr != nil {
		logx.Error(readErr)
		return http.StatusInternalServerError, readErr.Error()
	}

	return statusCode, string(body)
}

func HttpPostJson(httpUrl string, params string) (int, string) {
	logx.Infof("HttpPostJson请求的url：" + httpUrl + "请求的参数：" + params)
	jsonStr := []byte(params)
	url := httpUrl
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logx.Error(err)
		return http.StatusInternalServerError, err.Error()
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, clientErr := client.Do(req)
	if clientErr != nil {
		logx.Error(clientErr)
		return http.StatusInternalServerError, clientErr.Error()
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	body, _ := ioutil.ReadAll(resp.Body)

	return statusCode, string(body)
}

func HttpGet(url string, headers map[string]string) (int, string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logx.Error(err)
		return http.StatusInternalServerError, err.Error()
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, clientErr := client.Do(req)
	if clientErr != nil {
		logx.Error(clientErr)
		return http.StatusInternalServerError, clientErr.Error()
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	statusCode := resp.StatusCode
	if readErr != nil {
		logx.Error(readErr)
		return http.StatusInternalServerError, readErr.Error()
	}

	return statusCode, string(body)
}

func HttpBuildQuery(params map[string]string) (paramStr string) {
	paramsArr := make([]string, 0, len(params))
	for k, v := range params {
		paramsArr = append(paramsArr, fmt.Sprintf("%s=%s", k, v))
	}
	//fmt.Println(params_arr)
	paramStr = strings.Join(paramsArr, "&")
	return paramStr
}

func ParseResponseBody(resp *http.Response) (int, string) {
	body, readErr := ioutil.ReadAll(resp.Body)
	statusCode := resp.StatusCode
	if readErr != nil {
		logx.Error(readErr)
		return http.StatusInternalServerError, readErr.Error()
	}

	return statusCode, string(body)
}

func ApiSuccess(data interface{}) *ApiResponse {
	return &ApiResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func ApiError(message string, data interface{}) *ApiResponse {
	return &ApiResponse{
		Code:    1,
		Message: message,
		Data:    data,
	}
}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
