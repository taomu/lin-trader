package binance

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
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
		host:       "https://fapi.binance.com",
		httpClient: &http.Client{},
	}
}

// SetUrl 设置 REST API 主机地址
func (ra *RestApi) SetHost(host string) {
	ra.host = host
}

// HMAC-SHA256, 返回 hex 小写
func createHMACSignatureHex(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// 尝试解析 PEM (PKCS#8 或 PKCS#1) 并用 RSA PKCS1v1.5+SHA256 签名，返回 base64 (未 url-encode)
func createRSASignBase64(privateKeyPEM []byte, payload string) (string, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return "", errors.New("invalid private key PEM")
	}

	var priv *rsa.PrivateKey
	// 先尝试 PKCS#8
	if parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if k, ok := parsed.(*rsa.PrivateKey); ok {
			priv = k
		} else {
			return "", errors.New("pkcs8 key is not RSA")
		}
	} else {
		// 再尝试 PKCS#1
		if k, err2 := x509.ParsePKCS1PrivateKey(block.Bytes); err2 == nil {
			priv = k
		} else {
			return "", fmt.Errorf("parse private key failed: %v / %v", err, err2)
		}
	}

	h := sha256.Sum256([]byte(payload))
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, h[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

// 统一发送函数（支持 GET/DELETE 将参数放 URL，POST/PUT 将参数放 body (form-urlencoded)）
// 若 apiInfo != nil 则会自动加 timestamp/recvWindow（如果未提供），并签名（默认 HMAC；若 apiInfo.Secret 看起来像 PEM，则尝试 RSA）
func (ra *RestApi) sendRequest(path, method string, params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	method = strings.ToUpper(method)

	// 转成 url.Values（保证 deterministic order）
	values := url.Values{}
	for k, v := range params {
		values.Set(k, fmt.Sprintf("%v", v))
	}

	// 需要签名的接口：加入 timestamp（毫秒）与默认 recvWindow（可覆盖）
	if apiInfo != nil {
		if values.Get("timestamp") == "" {
			values.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/1e6))
		}
		// 如果调用方没传 recvWindow，默认 5000 ms（示例和建议）
		if values.Get("recvWindow") == "" {
			values.Set("recvWindow", "5000")
		}

		// 生成签名（在将 signature 放入 values 之前，用 values.Encode() 作为 payload）
		payload := values.Encode()

		// 简单 heuristic 判断是否 RSA PEM（如果 secret 是 PEM 文本）
		secretTrim := strings.TrimSpace(apiInfo.Secret)
		if strings.HasPrefix(secretTrim, "-----BEGIN") {
			// RSA：返回 base64，需要 url-encode 放到 signature
			base64Sig, err := createRSASignBase64([]byte(apiInfo.Secret), payload)
			if err != nil {
				return "", fmt.Errorf("BN_API RSA 签名失败: %v", err)
			}
			// url.Values.Encode 会自动做 QueryEscape，所以直接 Set(base64) 再 Encode 即可
			values.Set("signature", base64Sig)
		} else {
			// HMAC：hex 小写
			hexSig := createHMACSignatureHex(apiInfo.Secret, payload)
			values.Set("signature", hexSig)
		}
	}

	var req *http.Request
	var err error
	fullURL := ra.host + path

	switch method {
	case "GET", "DELETE":
		qs := values.Encode()
		if qs != "" {
			fullURL = fullURL + "?" + qs
		}
		req, err = http.NewRequest(method, fullURL, nil)
	case "POST", "PUT":
		// POST/PUT: 使用 form-urlencoded body（不要用 JSON）
		body := values.Encode()
		req, err = http.NewRequest(method, fullURL, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	default:
		return "", fmt.Errorf("%s BN_API unsupported http method: %s", fullURL, method)
	}
	if err != nil {
		return "", fmt.Errorf("%s %s BN_API 创建请求失败: %v", fullURL, method, err)
	}

	// API key header（不管 HMAC 还是 RSA 都需要）
	if apiInfo != nil && apiInfo.Key != "" {
		req.Header.Set("X-MBX-APIKEY", apiInfo.Key)
	}

	// 调试日志（必要时打开）：打印最终 payload、签名、请求 URL/Body，便于与 openssl/curl 输出比对
	//fmt.Println("DEBUG payload:", values.Encode())
	//fmt.Println("DEBUG request url:", req.URL.String())
	//if method == "POST" { bodyBytes, _ := io.ReadAll(req.Body); fmt.Println("DEBUG body:", string(bodyBytes)); req.Body = io.NopCloser(strings.NewReader(string(bodyBytes))) }

	resp, err := ra.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s %s BN_API 请求发送失败: %v", fullURL, method, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s %s BN_API 读取响应失败: %v", fullURL, method, err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("%s %s BN_API 请求失败，状态码: %d, 响应: %s", fullURL, method, resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}

func (ra *RestApi) ExchangeInfo(params map[string]interface{}) (string, error) {
	url := "/fapi/v1/exchangeInfo"
	method := "GET"
	return ra.sendRequest(url, method, params, nil)
}

func (ra *RestApi) PlaceOrder(params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
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

func (ra *RestApi) Account(params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	url := "/fapi/v3/account"
	method := "GET"
	return ra.sendRequest(url, method, params, apiInfo)
}

// 获取用户数据流 listenKey（仅需 API Key，不签名）
func (ra *RestApi) StartUserDataStream(apiKey string) (string, error) {
	fullURL := ra.host + "/fapi/v1/listenKey"
	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("X-MBX-APIKEY", apiKey)

	resp, err := ra.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}
	return string(body), nil
}

// 保活用户数据流 listenKey（仅需 API Key，不签名）
func (ra *RestApi) KeepaliveUserDataStream(apiKey, listenKey string) error {
	fullURL := ra.host + "/fapi/v1/listenKey"
	values := url.Values{}
	values.Set("listenKey", listenKey)
	req, err := http.NewRequest("PUT", fullURL, strings.NewReader(values.Encode()))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("X-MBX-APIKEY", apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := ra.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (ra *RestApi) CancelOrder(params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	url := "/fapi/v1/order"
	method := "DELETE"
	return ra.sendRequest(url, method, params, apiInfo)
}

func (ra *RestApi) GetLeverageBrackets(params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	url := "/fapi/v1/leverageBracket"
	method := "GET"
	return ra.sendRequest(url, method, params, apiInfo)
}

// 获取账户设置
func (ra *RestApi) GetAccountConfig(params map[string]interface{}, apiInfo *lintypes.ApiInfo) (string, error) {
	url := "/fapi/v1/accountConfig"
	method := "GET"
	return ra.sendRequest(url, method, params, apiInfo)
}
