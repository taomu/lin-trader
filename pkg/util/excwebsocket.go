package util

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ExcWebsocket struct {
	URL          string
	conn         *websocket.Conn
	OnConnect    func()
	OnMessage    func(msg string)
	OnClose      func()
	PingInterval time.Duration
	manualClose  bool // 标志是否为手动关闭
	mu           sync.Mutex
}

func NewExcWebsocket(url string) *ExcWebsocket {
	return &ExcWebsocket{
		URL:          url,
		PingInterval: 10 * time.Second, // 默认心跳间隔为 10 秒
		manualClose:  false,
	}
}

func (sw *ExcWebsocket) Connect() error {
	// sw.mu.Lock()
	// defer sw.mu.Unlock()
	sw.manualClose = false
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}
	if sw.URL == "" {
		return fmt.Errorf("url is empty")
	}
	conn, _, err := dialer.Dial(sw.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	sw.conn = conn

	// 设置自定义 Ping 处理器
	sw.conn.SetPingHandler(func(appData string) error {
		// 手动回复 Pong（或按需处理）
		err := sw.conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
		return err
	})

	// 触发 OnConnect 回调
	if sw.OnConnect != nil {
		sw.OnConnect()
	}

	go sw.readLoop()
	go sw.pingLoop()

	return nil
}

func (sw *ExcWebsocket) readLoop() {
	defer sw.handleClose()
	for {
		_, message, err := sw.conn.ReadMessage()
		if err != nil {
			break
		}
		// 判断是否为字符串 "ping"
		if string(message) == "ping" {
			err3 := sw.conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			if err3 != nil {
				return
			}
			continue
		}
		// 触发 OnMessage 回调
		if sw.OnMessage != nil {
			sw.OnMessage(string(message))
		}
	}
}

// 修改 pingLoop 逻辑
func (sw *ExcWebsocket) pingLoop() {
	ticker := time.NewTicker(sw.PingInterval)
	defer ticker.Stop()

	//使用range替换
	for range ticker.C {
		sw.mu.Lock()
		err := sw.conn.WriteMessage(websocket.PingMessage, nil)
		sw.mu.Unlock()

		if err != nil {
			sw.Close() // ✅ 主动触发关闭流程
			return
		}
	}
}

func (sw *ExcWebsocket) handleClose() {
	if sw.OnClose != nil {
		sw.OnClose()
	}
	// 自动重连逻辑
	if !sw.manualClose { // 只有非手动关闭时才进行自动重连
		for {
			time.Sleep(3 * time.Second) // 重连间隔
			err := sw.Connect()
			if err == nil {
				return
			}
		}
	}
}

func (sw *ExcWebsocket) Push(data string) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.conn.WriteMessage(websocket.TextMessage, []byte(data))
}

func (sw *ExcWebsocket) Close() error {
	sw.manualClose = true // 设置为手动关闭
	return sw.conn.Close()
}

// // example
// func main() {
// 	client := NewSimpleWebsocket("wss://example.com/socket")

// 	client.OnConnect = func() {
// 		fmt.Println("Connected to server")
// 	}

// 	client.OnMessage = func(msg string) {
// 		fmt.Printf("Received message: %s\n", msg)
// 	}

// 	client.OnClose = func() {
// 		fmt.Println("Connection closed")
// 	}

// 	err := client.Connect()
// 	if err != nil {
// 		fmt.Println("Failed to connect:", err)
// 		return
// 	}

// 	// 模拟发送消息
// 	err = client.Push("Hello, WebSocket!")
// 	if err != nil {
// 		fmt.Println("Failed to send message:", err)
// 		return
// 	}

// 	// 等待关闭
// 	time.Sleep(30 * time.Second)
// 	client.Close()
// }
