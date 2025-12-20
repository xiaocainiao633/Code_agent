package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/config"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/utils"
)

// MessageType 消息类型
type MessageType string

const (
	MessageTypeTaskProgress MessageType = "task_progress"
	MessageTypeAgentThought MessageType = "agent_thought"
	MessageTypeSystem       MessageType = "system"
	MessageTypePing         MessageType = "ping"
	MessageTypePong         MessageType = "pong"
)

// Message WebSocket消息
type Message struct {
	Type      MessageType            `json:"type"`
	TaskID    string                 `json:"task_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Client WebSocket客户端
type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	hub      *Hub
	taskID   string
	clientID string
}

// Hub WebSocket连接管理中心
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	config     *config.WebSocketConfig
	mu         sync.RWMutex
}

// NewHub 创建WebSocket Hub
func NewHub(cfg *config.WebSocketConfig) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		config:     cfg,
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	utils.Info("WebSocket Hub started")
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			utils.Info("WebSocket client registered: %s", client.clientID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				utils.Info("WebSocket client unregistered: %s", client.clientID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// 客户端发送通道已满，关闭连接
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToTask 向特定任务广播消息
func (h *Hub) BroadcastToTask(taskID string, message Message) {
	message.TaskID = taskID
	data, err := json.Marshal(message)
	if err != nil {
		utils.Error("Failed to marshal message for task %s: %v", taskID, err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.taskID == taskID {
			select {
			case client.send <- data:
			default:
				// 客户端发送通道已满，准备关闭
				go h.unregisterClient(client)
			}
		}
	}
}

// BroadcastAgentThought 广播Agent思考流
func (h *Hub) BroadcastAgentThought(taskID string, thought string, step string) {
	message := Message{
		Type:      MessageTypeAgentThought,
		TaskID:    taskID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"thought": thought,
			"step":    step,
		},
	}
	h.BroadcastToTask(taskID, message)
}

// BroadcastTaskProgress 广播任务进度
func (h *Hub) BroadcastTaskProgress(taskID string, progress int, status string, message string) {
	msg := Message{
		Type:      MessageTypeTaskProgress,
		TaskID:    taskID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"progress": progress,
			"status":   status,
			"message":  message,
		},
	}
	h.BroadcastToTask(taskID, msg)
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	select {
	case h.unregister <- client:
	default:
		// Hub已关闭，直接关闭客户端
		client.conn.Close()
	}
}

// WebSocketHandler WebSocket处理器
type WebSocketHandler struct {
	hub *Hub
	upgrader websocket.Upgrader
}

// NewWebSocketHandler 创建WebSocket处理器
func NewWebSocketHandler(hub *Hub, corsConfig *config.CORSConfig) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// CORS检查 - 开发环境下允许所有来源
				origin := r.Header.Get("Origin")
				if origin == "" {
					return true
				}
				
				// 开发环境：允许所有来源
				for _, allowedOrigin := range corsConfig.AllowedOrigins {
					if allowedOrigin == "*" {
						return true
					}
				}
				
				// 生产环境：检查具体来源
				for _, allowedOrigin := range corsConfig.AllowedOrigins {
					if allowedOrigin == origin {
						return true
					}
				}
				
				// 默认允许本地开发
				if origin == "http://localhost:3000" || origin == "http://localhost:3001" {
					return true
				}
				
				return false
			},
		},
	}
}

// HandleTaskWebSocket 处理任务WebSocket连接
func (h *WebSocketHandler) HandleTaskWebSocket(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskId")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	utils.Info("WebSocket connection request for task: %s", taskID)

	// 升级HTTP连接到WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Error("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	clientID := fmt.Sprintf("ws_%s_%d", taskID, time.Now().UnixNano())
	client := &Client{
		conn:     conn,
		send:     make(chan []byte, 256),
		hub:      h.hub,
		taskID:   taskID,
		clientID: clientID,
	}

	// 注册客户端
	h.hub.register <- client

	// 启动客户端读写协程
	go client.writePump()
	go client.readPump()

	utils.Info("WebSocket client connected: %s for task: %s", clientID, taskID)
}

// readPump 读取客户端消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(c.hub.config.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(c.hub.config.PongTimeout))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.hub.config.PongTimeout))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.Error("WebSocket read error: %v", err)
			}
			break
		}

		// 处理ping消息
		var msg Message
		if err := json.Unmarshal(message, &msg); err == nil && msg.Type == MessageTypePing {
			// 回复pong
			pong := Message{
				Type:      MessageTypePong,
				Timestamp: time.Now(),
				Data:      map[string]interface{}{},
			}
			if data, err := json.Marshal(pong); err == nil {
				c.send <- data
			}
		}
	}
}

// writePump 向客户端发送消息
func (c *Client) writePump() {
	ticker := time.NewTicker(c.hub.config.PingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub关闭了通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			// 发送ping消息
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			ping := Message{
				Type:      MessageTypePing,
				Timestamp: time.Now(),
				Data:      map[string]interface{}{},
			}
			if data, err := json.Marshal(ping); err == nil {
				if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
					return
				}
			}
		}
	}
}