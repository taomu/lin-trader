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

	"github.com/taomu/lin-trader/pkg/lintypes"
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

func (ra *RestApi) SetHost(host string) {
	ra.host = host
}

func (ra *RestApi) sendRequest(path, method string, params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	method = strings.ToUpper(method)

	// --------- 1) 构建 requestPath ---------
	requestPath := path
	fullURL := ra.host + path
	if method == "GET" || method == "DELETE" {
		if len(params) > 0 {
			q := url.Values{}
			for k, v := range params {
				q.Add(k, fmt.Sprintf("%v", v))
			}
			query := q.Encode()
			requestPath = path + "?" + query
			fullURL = fullURL + "?" + query
		}
	}

	// --------- 2) 构建 body ---------
	bodyStr := ""
	if method == "POST" || method == "PUT" {
		if len(params) > 0 {
			b, err := json.Marshal(params)
			if err != nil {
				return "", fmt.Errorf("JSON编码失败: %v", err)
			}
			bodyStr = string(b)
		}
	}

	// --------- 3) 创建请求 ---------
	var req *http.Request
	var err error
	if method == "GET" || method == "DELETE" {
		req, err = http.NewRequest(method, fullURL, nil)
	} else {
		req, err = http.NewRequest(method, fullURL, strings.NewReader(bodyStr))
		req.Header.Set("Content-Type", "application/json")
	}
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// --------- 4) 设置认证头 ---------
	if apiInfo != nil {
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		prehash := timestamp + method + requestPath + bodyStr

		mac := hmac.New(sha256.New, []byte(apiInfo.Secret))
		mac.Write([]byte(prehash))
		sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		req.Header.Set("OK-ACCESS-KEY", apiInfo.Key)
		req.Header.Set("OK-ACCESS-SIGN", sign)
		req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Set("OK-ACCESS-PASSPHRASE", apiInfo.Passphrase)
		req.Header.Set("Content-Type", "application/json")
	}

	// --------- 5) 发送请求 ---------
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

// 获取24小时tick数据
func (ra *RestApi) Tickers24h(params map[string]interface{}) (string, error) {
	url := "/api/v5/market/tickers"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

// 获取所有合约交易对信息
func (ra *RestApi) Instruments(params map[string]interface{}) (string, error) {
	url := "/api/v5/public/instruments"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

// 获取持仓信息
func (ra *RestApi) GetPositions(params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	url := "/api/v5/account/positions"
	method := "GET"
	return ra.sendRequest(url, method, params, apiInfo)
}

// 获取资金费率信息
func (ra *RestApi) GetPremium(params map[string]interface{}) (string, error) {
	url := "/api/v5/public/funding-rate"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

// 提交订单
func (ra *RestApi) PlaceOrder(params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	url := "/api/v5/trade/order"
	method := "POST"
	return ra.sendRequest(url, method, params, apiInfo)
}
