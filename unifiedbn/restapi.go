package unifiedbn

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
		host:       "https://papi.binance.com", // Portfolio Margin API endpoint
		httpClient: &http.Client{},
	}
}

// HMAC-SHA256, returns hex lowercase
func createHMACSignatureHex(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (ra *RestApi) sendRequest(path, method string, params map[string]interface{}, apiInfo *types.ApiInfo) (string, error) {
	method = strings.ToUpper(method)

	values := url.Values{}
	for k, v := range params {
		values.Set(k, fmt.Sprintf("%v", v))
	}

	if apiInfo != nil {
		if values.Get("timestamp") == "" {
			values.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/1e6))
		}
		if values.Get("recvWindow") == "" {
			values.Set("recvWindow", "5000")
		}

		payload := values.Encode()
		hexSig := createHMACSignatureHex(apiInfo.Secret, payload)
		values.Set("signature", hexSig)
	}

	var req *http.Request
	var err error
	fullURL := ra.host + path

	if method == "GET" || method == "DELETE" {
		qs := values.Encode()
		if qs != "" {
			fullURL = fullURL + "?" + qs
		}
		req, err = http.NewRequest(method, fullURL, nil)
	} else if method == "POST" || method == "PUT" {
		body := values.Encode()
		req, err = http.NewRequest(method, fullURL, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		return "", fmt.Errorf("unsupported http method: %s", method)
	}
	if err != nil {
		return "", fmt.Errorf("creating request failed: %v", err)
	}

	if apiInfo != nil && apiInfo.Key != "" {
		req.Header.Set("X-MBX-APIKEY", apiInfo.Key)
	}

	resp, err := ra.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response failed: %v", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("request failed with status code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}

func (ra *RestApi) Account(params map[string]interface{}, apiInfo *types.ApiInfo) (string, error) {
	path := "/papi/v1/account"
	method := "GET"
	return ra.sendRequest(path, method, params, apiInfo)
}