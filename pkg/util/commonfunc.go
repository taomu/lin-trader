package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GetUnitTime(interval string, timePrecision string) int64 {
	if !(interval == "1m" || interval == "5m" || interval == "15m" || interval == "1h" || interval == "4h" || interval == "1d") {
		panic("GetUnitTime param interval error, interval is " + interval)
	}
	if !(timePrecision == "s" || timePrecision == "ms") {
		panic("GetUnitTime param timePrecision error, timePrecision is" + timePrecision)
	}
	var miniTime int64
	if timePrecision == "ms" {
		miniTime = 1000
	} else {
		miniTime = 1
	}
	if interval == "1m" {
		return miniTime * 60
	}
	if interval == "5m" {
		return miniTime * 60 * 5
	} else if interval == "15m" {
		return miniTime * 60 * 15
	} else if interval == "1h" {
		return miniTime * 60 * 60
	} else if interval == "4h" {
		return miniTime * 60 * 60 * 4
	} else {
		return miniTime * 60 * 60 * 24
	}
}
func StrToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}
func StrToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}
func GetFile(fullPath string) *os.File {
	logDir := filepath.Dir(fullPath)
	err := os.MkdirAll(logDir, 0755) // 0755 表示目录的权限（可读写执行）
	if err != nil {
		fmt.Println("Failed to create log directory:", err)
		panic(err)
	}
	file, err2 := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err2 != nil {
		fmt.Println("Failed to open log file:", err2)
		panic(err2)
	}
	return file
}

// FeishuMessage 请求体结构
type FeishuMessage struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

// SendFeishuMessage 发送消息到飞书群聊
func SendFeishuMessage(webhookURL, text string) error {
	// 创建消息内容
	message := FeishuMessage{
		MsgType: "text",
	}
	message.Content.Text = text

	// 序列化消息
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// 发起 HTTP POST 请求
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(messageJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	fmt.Println("Message sent successfully!")
	return nil
}

// CountDecimalPlaces 获取小数点后面位数
func CountDecimalPlaces(input string) int {
	// 找到小数点的位置
	dotIndex := strings.Index(input, ".")
	if dotIndex == -1 {
		// 如果没有小数点，返回0
		return 0
	}
	// 获取小数点后面的部分并去掉末尾的0
	decimalPart := strings.TrimRight(input[dotIndex+1:], "0")

	// 返回去0后的小数位数
	return len(decimalPart)
}

// Round 保留n位小数
func Round(f float64, n int) float64 {
	// 通过 10 的 n 次方放大后再缩小
	pow := math.Pow(10, float64(n))
	return math.Round(f*pow) / pow
}

// SliceReverse 泛型方法，支持任意类型切片倒序
func SliceReverse[T any](slice *[]T) {
	if slice == nil {
		return
	}
	s := *slice
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// ToTimeStamp 将时间字符串转换为时间戳
func ToTimeStamp(timeStr string, timeStampType string) int64 {
	// 时间格式
	const layout = "2006-01-02 15:04:05"
	// 加载东八区时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0
	}
	// 解析时间字符串
	t, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return 0
	}
	// 转换为毫秒级时间戳
	if timeStampType == "ms" {
		return t.UnixMilli()
	}
	if timeStampType == "s" {
		return t.Unix()
	}
	return 0
}

// TsToTimeUtc8 时间戳转UTC+8时间，返回格式化的字符串
func TsToTimeUtc8(ts int64) string {
	// 时间格式
	const layout = "2006-01-02 15:04:05"

	// 判断时间戳是秒级还是毫秒级
	var timestamp int64
	if ts > 1000000000000 {
		// 毫秒级时间戳，转为秒级
		timestamp = ts / 1000
	} else {
		// 秒级时间戳
		timestamp = ts
	}

	// 将时间戳转换为 UTC 时间
	utcTime := time.Unix(timestamp, 0).UTC()

	// 将 UTC 时间转换为 UTC+8（北京时间）
	// 假设时区为 UTC+8
	utc8 := time.FixedZone("UTC+8", 8*60*60)
	localTime := utcTime.In(utc8)

	// 返回格式化的时间字符串
	return localTime.Format(layout)
}
func TimeUtc8ToTs(timeStr string, tsType string) int64 {
	// 支持多种时间格式
	layouts := []string{
		"2006年1月2日下午3:04",     // 中文格式下午时间
		"2006年1月2日上午3:04",     // 中文格式上午时间
		"2006年1月2日15:04",      // 中文格式24小时制
		"2006-1-2 15:04:05",   // 数字格式
		"2006-01-02 15:04:05", // 标准数字格式
	}

	// 加载东八区时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0
	}

	// 预处理时间字符串
	timeStr = strings.ReplaceAll(timeStr, " ", "") // 移除所有空格

	// 处理下午时间(转换为24小时制)
	if strings.Contains(timeStr, "下午") {
		parts := strings.Split(timeStr, "下午")
		timePart := strings.Split(parts[1], ":")
		hour, _ := strconv.Atoi(timePart[0])
		if hour < 12 {
			hour += 12
		}
		timeStr = parts[0] + strconv.Itoa(hour) + ":" + timePart[1]
	} else if strings.Contains(timeStr, "上午") {
		timeStr = strings.Replace(timeStr, "上午", "", 1)
	}

	// 尝试用各种格式解析时间
	var t time.Time
	for _, layout := range layouts {
		t, err = time.ParseInLocation(layout, timeStr, loc)
		if err == nil {
			break
		}
	}
	if err != nil {
		return 0
	}

	// 转换为时间戳
	if tsType == "ms" {
		return t.UnixMilli()
	}
	if tsType == "s" {
		return t.Unix()
	}
	return 0
}
func UtcToUtc8(utcTimeStr string) string {
	//示范参数utcTimeStr := "2025-02-11 15:45 (UTC)"
	// 解析 UTC 时间
	locUTC, _ := time.LoadLocation("UTC")
	parsedTime, err := time.ParseInLocation("2006-01-02 15:04", strings.TrimSuffix(utcTimeStr, " (UTC)"), locUTC)
	if err != nil {
		fmt.Println("时间解析失败:", err)
		return ""
	}
	// 转换为 UTC+8
	locUTC8, _ := time.LoadLocation("Asia/Shanghai")
	utc8Time := parsedTime.In(locUTC8)
	return utc8Time.Format("2006-01-02 15:04:05")
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 创建带种子的随机数生成器
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = letters[r.Intn(len(letters))]
	}
	return string(result)
}
func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
func StructToJson(v interface{}) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
