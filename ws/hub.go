package ws

import (
	"encoding/json"
	"lotteryBackend/global"
	"sync"

	"go.uber.org/zap"
)

// Hub WebSocket连接管理中心
type Hub struct {
	// 所有客户端
	Clients map[string]*Client
	// 客户端锁
	ClientsMutex sync.RWMutex

	// 房间 -> 客户端列表
	Rooms map[string]map[string]*Client
	// 房间锁
	RoomsMutex sync.RWMutex

	// 用户ID -> 客户端（用于点对点消息）
	UserClients map[uint]map[string]*Client
	// 用户客户端锁
	UserClientsMutex sync.RWMutex

	// 注册客户端通道
	Register chan *Client
	// 注销客户端通道
	Unregister chan *Client

	// 广播消息通道
	Broadcast chan *BroadcastMessage
	// 房间消息通道
	RoomBroadcast chan *RoomMessage
	// 用户消息通道
	UserMessage chan *UserMessage
}

// BroadcastMessage 广播消息
type BroadcastMessage struct {
	Type    string
	Payload interface{}
}

// RoomMessage 房间消息
type RoomMessage struct {
	RoomID  string
	Type    string
	Payload interface{}
}

// UserMessage 用户消息
type UserMessage struct {
	UserID  uint
	Type    string
	Payload interface{}
}

// 全局Hub实例
var (
	globalHub *Hub
	hubOnce   sync.Once
)

// GetHub 获取全局Hub实例
func GetHub() *Hub {
	hubOnce.Do(func() {
		globalHub = NewHub()
		go globalHub.Run()
	})
	return globalHub
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	return &Hub{
		Clients:       make(map[string]*Client),
		Rooms:         make(map[string]map[string]*Client),
		UserClients:   make(map[uint]map[string]*Client),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Broadcast:     make(chan *BroadcastMessage, 256),
		RoomBroadcast: make(chan *RoomMessage, 256),
		UserMessage:   make(chan *UserMessage, 256),
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcast(message)

		case message := <-h.RoomBroadcast:
			h.roomBroadcast(message)

		case message := <-h.UserMessage:
			h.userMessage(message)
		}
	}
}

// registerClient 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.ClientsMutex.Lock()
	h.Clients[client.ID] = client
	h.ClientsMutex.Unlock()

	// 注册用户客户端映射
	if client.UserID > 0 {
		h.UserClientsMutex.Lock()
		if h.UserClients[client.UserID] == nil {
			h.UserClients[client.UserID] = make(map[string]*Client)
		}
		h.UserClients[client.UserID][client.ID] = client
		h.UserClientsMutex.Unlock()
	}

	global.GVA_LOG.Info("WebSocket客户端已连接",
		zap.String("clientId", client.ID),
		zap.Uint("userId", client.UserID))
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.ClientsMutex.Lock()
	if _, ok := h.Clients[client.ID]; ok {
		delete(h.Clients, client.ID)
		close(client.Send)
	}
	h.ClientsMutex.Unlock()

	// 从所有房间移除
	h.RoomsMutex.Lock()
	for roomID := range client.Rooms {
		if room, ok := h.Rooms[roomID]; ok {
			delete(room, client.ID)
			if len(room) == 0 {
				delete(h.Rooms, roomID)
			}
		}
	}
	h.RoomsMutex.Unlock()

	// 从用户映射移除
	if client.UserID > 0 {
		h.UserClientsMutex.Lock()
		if clients, ok := h.UserClients[client.UserID]; ok {
			delete(clients, client.ID)
			if len(clients) == 0 {
				delete(h.UserClients, client.UserID)
			}
		}
		h.UserClientsMutex.Unlock()
	}

	global.GVA_LOG.Info("WebSocket客户端已断开",
		zap.String("clientId", client.ID),
		zap.Uint("userId", client.UserID))
}

// broadcast 广播消息给所有客户端
func (h *Hub) broadcast(message *BroadcastMessage) {
	data, err := json.Marshal(Message{
		Type:    message.Type,
		Payload: message.Payload,
	})
	if err != nil {
		return
	}

	h.ClientsMutex.RLock()
	defer h.ClientsMutex.RUnlock()

	for _, client := range h.Clients {
		select {
		case client.Send <- data:
		default:
			go h.unregisterClient(client)
		}
	}
}

// roomBroadcast 广播消息给房间内的客户端
func (h *Hub) roomBroadcast(message *RoomMessage) {
	data, err := json.Marshal(Message{
		Type:    message.Type,
		Payload: message.Payload,
	})
	if err != nil {
		return
	}

	h.RoomsMutex.RLock()
	defer h.RoomsMutex.RUnlock()

	if room, ok := h.Rooms[message.RoomID]; ok {
		for _, client := range room {
			select {
			case client.Send <- data:
			default:
				go h.unregisterClient(client)
			}
		}
	}
}

// userMessage 发送消息给指定用户
func (h *Hub) userMessage(message *UserMessage) {
	data, err := json.Marshal(Message{
		Type:    message.Type,
		Payload: message.Payload,
	})
	if err != nil {
		return
	}

	h.UserClientsMutex.RLock()
	defer h.UserClientsMutex.RUnlock()

	if clients, ok := h.UserClients[message.UserID]; ok {
		for _, client := range clients {
			select {
			case client.Send <- data:
			default:
				go h.unregisterClient(client)
			}
		}
	}
}

// JoinRoom 客户端加入房间
func (h *Hub) JoinRoom(client *Client, roomID string) {
	h.RoomsMutex.Lock()
	defer h.RoomsMutex.Unlock()

	if h.Rooms[roomID] == nil {
		h.Rooms[roomID] = make(map[string]*Client)
	}
	h.Rooms[roomID][client.ID] = client

	global.GVA_LOG.Debug("客户端加入房间",
		zap.String("clientId", client.ID),
		zap.String("roomId", roomID))
}

// LeaveRoom 客户端离开房间
func (h *Hub) LeaveRoom(client *Client, roomID string) {
	h.RoomsMutex.Lock()
	defer h.RoomsMutex.Unlock()

	if room, ok := h.Rooms[roomID]; ok {
		delete(room, client.ID)
		if len(room) == 0 {
			delete(h.Rooms, roomID)
		}
	}

	global.GVA_LOG.Debug("客户端离开房间",
		zap.String("clientId", client.ID),
		zap.String("roomId", roomID))
}

// HandleClientMessage 处理客户端消息
func (h *Hub) HandleClientMessage(client *Client, msg *ClientMessage) {
	switch msg.Type {
	case TypeShakeScore:
		GetShakeHandler().HandleShakeScore(client, msg.Payload)
	}
}

// ==================== 对外广播方法 ====================

// BroadcastToAll 广播给所有客户端
func (h *Hub) BroadcastToAll(msgType string, payload interface{}) {
	h.Broadcast <- &BroadcastMessage{
		Type:    msgType,
		Payload: payload,
	}
}

// BroadcastToRoom 广播给指定房间
func (h *Hub) BroadcastToRoom(roomID string, msgType string, payload interface{}) {
	h.RoomBroadcast <- &RoomMessage{
		RoomID:  roomID,
		Type:    msgType,
		Payload: payload,
	}
}

// SendToUser 发送给指定用户
func (h *Hub) SendToUser(userID uint, msgType string, payload interface{}) {
	h.UserMessage <- &UserMessage{
		UserID:  userID,
		Type:    msgType,
		Payload: payload,
	}
}

// GetRoomClientCount 获取房间客户端数量
func (h *Hub) GetRoomClientCount(roomID string) int {
	h.RoomsMutex.RLock()
	defer h.RoomsMutex.RUnlock()

	if room, ok := h.Rooms[roomID]; ok {
		return len(room)
	}
	return 0
}

// GetOnlineCount 获取在线客户端数量
func (h *Hub) GetOnlineCount() int {
	h.ClientsMutex.RLock()
	defer h.ClientsMutex.RUnlock()
	return len(h.Clients)
}

// GetOnlineUserCount 获取在线用户数量
func (h *Hub) GetOnlineUserCount() int {
	h.UserClientsMutex.RLock()
	defer h.UserClientsMutex.RUnlock()
	return len(h.UserClients)
}

// IsUserOnline 检查用户是否在线
func (h *Hub) IsUserOnline(userID uint) bool {
	h.UserClientsMutex.RLock()
	defer h.UserClientsMutex.RUnlock()
	clients, ok := h.UserClients[userID]
	return ok && len(clients) > 0
}
