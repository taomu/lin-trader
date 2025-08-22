package okx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/taomu/lin-trader/pkg/types"
)

type RestApi struct {
	host       string
	httpClient *http.Client
}

func NewRestApi() *RestApi {
	return &RestApi{
		host:       "https://www.okx.com",
		httpClient: &http.Client{},
	}
}

func (ra *RestApi) createSign(timestamp, method, requestPath, secretKey, body string) string {
	message := timestamp + method + requestPath + body
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (ra *RestApi) sendRequest(path, method string, params map[string]interface{}, apiInfo *types.ApiInfo) (string, error) {
	var req *http.Request
	var err error

	// 构建完整URL
	fullURL := ra.host + path
	method = strings.ToUpper(method)
	var body string

	// 处理GET请求参数
	if method == "GET" && len(params) > 0 {
		query := url.Values{}
		for k, v := range params {
			query.Add(k, fmt.Sprintf("%v", v))
		}
		fullURL = fullURL + "?" + query.Encode()
	} else if method == "POST" && len(params) > 0 {
		jsonData, err2 := json.Marshal(params)
		if err2 != nil {
			return "", fmt.Errorf("JSON编码失败: %v", err2)
		}
		body = string(jsonData)
	}

	// 创建HTTP请求
	req, err = http.NewRequest(method, fullURL, strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 处理需要认证的请求
	if apiInfo != nil {
		timestamp := time.Now().UTC().Format(time.RFC3339)
		sign := ra.createSign(timestamp, method, path, apiInfo.Secret, body)

		req.Header.Set("OK-ACCESS-KEY", apiInfo.Key)
		req.Header.Set("OK-ACCESS-PASSPHRASE", apiInfo.Passphrase)
		req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Set("OK-ACCESS-SIGN", sign)
	}

	// 发送请求
	resp, err := ra.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}

func (ra *RestApi) Tickers24h(params map[string]interface{}) (string, error) {
	url := "/api/v5/market/tickers"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

func (ra *RestApi) Instruments(params map[string]interface{}) (string, error) {
	url := "/api/v5/public/instruments"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

func (ra *RestApi) GetPositions(params map[string]interface{}, apiInfo *types.ApiInfo) (string, error) {
	url := "/api/v5/account/positions"
	method := "GET"
	return ra.sendRequest(url, method, params, apiInfo)
}

func (ra *RestApi) GetPremium(params map[string]interface{}) (string, error) {
	url := "/api/v5/public/funding-rate"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

func (ra *RestApi) IndexTickers(params map[string]interface{}) (string, error) {
	url := "/api/v5/market/index-tickers"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}
