package binance

import (
	"crypto/hmac"
	"crypto/sha256"
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
		host:       "https://fapi.binance.com",
		httpClient: &http.Client{},
	}
}

func (ra *RestApi) createSign(secretKey string, params url.Values) string {
	// 币安签名方式：对查询字符串进行HMAC SHA256签名
	message := params.Encode()
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	return fmt.Sprintf("%x", h.Sum(nil)) // 币安使用16进制格式签名
}

func (ra *RestApi) sendRequest(path, method string, params map[string]interface{}, apiInfo *types.ApiInfo) (string, error) {
	var req *http.Request
	var err error

	// 构建完整URL
	fullURL := ra.host + path
	method = strings.ToUpper(method)
	var body string

	// 转换参数为url.Values
	queryParams := url.Values{}
	for k, v := range params {
		queryParams.Add(k, fmt.Sprintf("%v", v))
	}

	// 添加时间戳(币安API要求)
	if apiInfo != nil {
		queryParams.Add("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/1e6))
	}

	// 处理GET请求参数
	if method == "GET" {
		if apiInfo != nil {
			// 添加签名
			queryParams.Add("signature", ra.createSign(apiInfo.Secret, queryParams))
		}
		fullURL = fullURL + "?" + queryParams.Encode()
	} else if method == "POST" {
		// POST请求处理
		if apiInfo != nil {
			queryParams.Add("signature", ra.createSign(apiInfo.Secret, queryParams))
		}
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
		req.Header.Set("X-MBX-APIKEY", apiInfo.Key) // 币安使用X-MBX-APIKEY头
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

func (ra *RestApi) ExchangeInfo(params map[string]interface{}) (string, error) {
	url := "/fapi/v1/exchangeInfo"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

func (ra *RestApi) PlaceOrder(params map[string]interface{}, apiInfo *types.ApiInfo) (string, error) {
	url := "/fapi/v1/order"
	method := "POST"
	return ra.sendRequest(url, method, params, apiInfo)
}

func (ra *RestApi) PremiumIndex(params map[string]interface{}) (string, error) {
	url := "/fapi/v1/premiumIndex"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

func (ra *RestApi) FundingInfo(params map[string]interface{}) (string, error) {
	url := "/fapi/v1/fundingInfo"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}
