package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // 连接的客户端
var broadcast = make(chan Message)           // 广播通道

// 配置 WebSocket
var upgrader = websocket.Upgrader{}

// 消息结构
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	// 创建静态文件服务
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// 配置 WebSocket 路由
	http.HandleFunc("/ws", handleConnections)

	// 启动监听广播通道的 Goroutine
	go handleMessages()

	// 启动服务器
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 升级初始 GET 请求到 WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有源
		},
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 确保关闭连接
	defer ws.Close()

	// 注册新的客户端
	clients[ws] = true

	for {
		var msg Message
		// 读取新消息
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// 将新消息发送到广播通道
		broadcast <- msg
		fmt.Println(msg)
	}
}

func handleMessages() {
	for {
		// 从广播通道中获取下一个消息
		msg := <-broadcast
		// 将消息发送给所有客户端
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
