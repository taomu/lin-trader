package binance

import (
	"crypto/hmac"
	"crypto/sha256"
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

	// 转换参数为 url.Values
	queryParams := url.Values{}
	for k, v := range params {
		queryParams.Add(k, fmt.Sprintf("%v", v))
	}

	// 处理需要签名的接口
	if apiInfo != nil {
		queryParams.Add("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/1e6))
		signature := ra.createSign(apiInfo.Secret, queryParams)
		queryParams.Add("signature", signature)
	}

	if method == "GET" || method == "DELETE" {
		// GET/DELETE 参数放在 URL
		fullURL = fullURL + "?" + queryParams.Encode()
		req, err = http.NewRequest(method, fullURL, nil)
	} else if method == "POST" {
		// POST 参数放在 body (form-url-encoded)
		req, err = http.NewRequest(method, fullURL, strings.NewReader(queryParams.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		return "", fmt.Errorf("不支持的HTTP方法: %s", method)
	}

	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置 API KEY
	if apiInfo != nil {
		req.Header.Set("X-MBX-APIKEY", apiInfo.Key)
	}

	// 发送请求
	resp, err := ra.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

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

func (ra *RestApi) Depth(params map[string]interface{}) (string, error) {
	url := "/fapi/v1/depth"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

func (ra *RestApi) Account(params map[string]interface{}, apiInfo *types.ApiInfo) (string, error) {
	url := "/fapi/v3/account"
	method := "GET"
	return ra.sendRequest(url, method, params, apiInfo)
}
