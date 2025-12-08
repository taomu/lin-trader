package bybit

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
		host:       "https://api.bybit.com",
		httpClient: &http.Client{},
	}
}

func (ra *RestApi) SetHost(host string) {
	ra.host = host
}

func (ra *RestApi) sendRequest(path, method string, params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	method = strings.ToUpper(method)

	fullURL := ra.host + path

	values := url.Values{}
	for k, v := range params {
		values.Set(k, fmt.Sprintf("%v", v))
	}

	bodyStr := ""
	if method == "GET" || method == "DELETE" {
		qs := values.Encode()
		if qs != "" {
			fullURL = fullURL + "?" + qs
		}
	} else {
		if len(params) > 0 {
			b, err := json.Marshal(params)
			if err != nil {
				return "", fmt.Errorf(fullURL+"JSON编码失败: err:%v params:%v", err, params)
			}
			bodyStr = string(b)
		}
	}

	var req *http.Request
	var err error
	if method == "GET" || method == "DELETE" {
		req, err = http.NewRequest(method, fullURL, nil)
	} else {
		req, err = http.NewRequest(method, fullURL, strings.NewReader(bodyStr))
		req.Header.Set("Content-Type", "application/json")
	}
	if err != nil {
		return "", fmt.Errorf(fullURL+"创建请求失败: %v", err)
	}

	if apiInfo != nil {
		ts := fmt.Sprintf("%d", time.Now().UnixMilli())
		recvWindow := "5000"
		var payload string
		if method == "GET" || method == "DELETE" {
			payload = ts + apiInfo.Key + recvWindow + values.Encode()
		} else {
			payload = ts + apiInfo.Key + recvWindow + bodyStr
		}
		mac := hmac.New(sha256.New, []byte(apiInfo.Secret))
		mac.Write([]byte(payload))
		sign := hex.EncodeToString(mac.Sum(nil))

		req.Header.Set("X-BAPI-API-KEY", apiInfo.Key)
		req.Header.Set("X-BAPI-SIGN", sign)
		req.Header.Set("X-BAPI-TIMESTAMP", ts)
		req.Header.Set("X-BAPI-RECV-WINDOW", recvWindow)
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := ra.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf(fullURL+"请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf(fullURL+"读取响应失败: %v", err)
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf(fullURL+"请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}
	return string(respBody), nil
}

func (ra *RestApi) Instruments(params map[string]interface{}) (string, error) {
	url := "/v5/market/instruments-info"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

func (ra *RestApi) Tickers24h(params map[string]interface{}) (string, error) {
	url := "/v5/market/tickers"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

func (ra *RestApi) GetPositions(params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	url := "/v5/position/list"
	method := "GET"
	return ra.sendRequest(url, method, params, apiInfo)
}

func (ra *RestApi) GetFundingRate(params map[string]interface{}) (string, error) {
	url := "/v5/market/funding-rate"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

func (ra *RestApi) PlaceOrder(params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	url := "/v5/order/create"
	method := "POST"
	return ra.sendRequest(url, method, params, apiInfo)
}

func (ra *RestApi) CancelOrder(params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	url := "/v5/order/cancel"
	method := "POST"
	return ra.sendRequest(url, method, params, apiInfo)
}
