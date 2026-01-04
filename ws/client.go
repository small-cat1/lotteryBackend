package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 写入等待时间
	writeWait = 10 * time.Second
	// 读取pong的等待时间
	pongWait = 60 * time.Second
	// ping周期，必须小于pongWait
	pingPeriod = (pongWait * 9) / 10
	// 最大消息大小
	maxMessageSize = 4096
)

// Client WebSocket客户端
type Client struct {
	ID         string          // 客户端唯一ID
	UserID     uint            // 用户ID（已登录用户）
	Conn       *websocket.Conn // WebSocket连接
	Hub        *Hub            // 所属Hub
	Send       chan []byte     // 发送消息通道
	Rooms      map[string]bool // 已加入的房间
	RoomsMutex sync.RWMutex    // 房间锁
	CloseChan  chan struct{}   // 关闭信号
	Closed     bool            // 是否已关闭
	CloseMutex sync.Mutex      // 关闭锁
}

// NewClient 创建新客户端
func NewClient(id string, userID uint, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:        id,
		UserID:    userID,
		Conn:      conn,
		Hub:       hub,
		Send:      make(chan []byte, 256),
		Rooms:     make(map[string]bool),
		CloseChan: make(chan struct{}),
		Closed:    false,
	}
}

// ReadPump 读取消息循环
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 记录异常关闭
			}
			break
		}

		// 处理客户端消息
		c.handleMessage(message)
	}
}

// WritePump 写入消息循环
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 批量发送队列中的消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.CloseChan:
			return
		}
	}
}

// handleMessage 处理客户端消息
func (c *Client) handleMessage(data []byte) {
	var msg ClientMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.SendMessage(TypeError, ErrorPayload{
			Code:    400,
			Message: "消息格式错误",
		})
		return
	}

	switch msg.Type {
	case "ping", TypeHeartbeat:
		// 心跳响应
		c.SendMessage(TypePong, HeartbeatPayload{
			Timestamp: time.Now().UnixMilli(),
		})

	case "join":
		// 加入房间
		c.JoinRoom(msg.Payload)

	case "leave":
		// 离开房间
		c.LeaveRoom(msg.Payload)

	default:
		// 其他消息转发到Hub处理
		c.Hub.HandleClientMessage(c, &msg)
	}
}

// SendMessage 发送消息给客户端
func (c *Client) SendMessage(msgType string, payload interface{}) error {
	msg := Message{
		Type:    msgType,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.CloseMutex.Lock()
	if c.Closed {
		c.CloseMutex.Unlock()
		return nil
	}
	c.CloseMutex.Unlock()

	select {
	case c.Send <- data:
	default:
		// 通道满了，关闭连接
		c.Close()
	}

	return nil
}

// JoinRoom 加入房间
func (c *Client) JoinRoom(roomID string) {
	c.RoomsMutex.Lock()
	c.Rooms[roomID] = true
	c.RoomsMutex.Unlock()

	c.Hub.JoinRoom(c, roomID)
}

// LeaveRoom 离开房间
func (c *Client) LeaveRoom(roomID string) {
	c.RoomsMutex.Lock()
	delete(c.Rooms, roomID)
	c.RoomsMutex.Unlock()

	c.Hub.LeaveRoom(c, roomID)
}

// InRoom 检查是否在房间中
func (c *Client) InRoom(roomID string) bool {
	c.RoomsMutex.RLock()
	defer c.RoomsMutex.RUnlock()
	return c.Rooms[roomID]
}

// Close 关闭连接
func (c *Client) Close() {
	c.CloseMutex.Lock()
	defer c.CloseMutex.Unlock()

	if c.Closed {
		return
	}

	c.Closed = true
	close(c.CloseChan)
	c.Conn.Close()
}
