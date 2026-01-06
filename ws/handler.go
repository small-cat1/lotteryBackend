package ws

import (
	"lotteryBackend/global"
	"lotteryBackend/pkg/h5jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该限制
	},
}

// Handler WebSocket处理器
type Handler struct {
	hub *Hub
}

// NewHandler 创建处理器
func NewHandler() *Handler {
	return &Handler{
		hub: GetHub(),
	}
}

// HandleConnection 处理WebSocket连接
// @Tags WebSocket
// @Summary WebSocket连接
// @Description H5前端WebSocket连接入口
// @Param token query string false "用户Token"
// @Param room query string false "房间ID"
// @Router /ws [get]
func (h *Handler) HandleConnection(c *gin.Context) {
	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("WebSocket升级失败", zap.Error(err))
		return
	}

	// 获取用户ID（可选，未登录用户也可连接）
	var userID uint = 0
	token := c.Query("token")
	if token != "" {
		// 验证Token获取用户ID
		userID = h.validateToken(token)
	}

	// 生成客户端ID
	clientID := uuid.New().String()

	// 创建客户端
	client := NewClient(clientID, userID, conn, h.hub)

	// 注册到Hub
	h.hub.Register <- client

	// 自动加入指定房间
	room := c.Query("room")
	if room != "" {
		client.JoinRoom(room)
	}

	// 发送连接成功消息
	client.SendMessage(TypeConnected, ConnectedPayload{
		ClientId: clientID,
		Message:  "连接成功",
	})

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump()
}

// HandleScreenConnection 大屏WebSocket连接
// @Tags WebSocket
// @Summary 大屏WebSocket连接
// @Description 大屏展示专用连接，自动加入screen房间
// @Param type query string true "大屏类型：checkin/danmaku/shake/draw"
// @Param activityId query string true "活动ID"
// @Router /ws/screen [get]
func (h *Handler) HandleScreenConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("WebSocket升级失败", zap.Error(err))
		return
	}

	clientID := uuid.New().String()
	client := NewClient(clientID, 0, conn, h.hub)

	h.hub.Register <- client

	// 根据类型加入对应房间
	screenType := c.Query("type")
	activityId := c.Query("activityId")
	// ⭐ 调试日志
	global.GVA_LOG.Info("主持人端连接",
		zap.String("clientId", clientID),
		zap.String("screenType", screenType),
		zap.String("activityId", activityId))

	// 加入通用大屏房间
	roomID := RoomTypeScreen + ":" + activityId
	client.JoinRoom(roomID)
	global.GVA_LOG.Info("主持人加入房间", zap.String("roomId", roomID))

	// 加入具体类型房间
	if screenType != "" {
		typeRoomID := screenType + ":" + activityId
		client.JoinRoom(typeRoomID)
		global.GVA_LOG.Info("主持人加入类型房间", zap.String("roomId", typeRoomID))
	}

	client.SendMessage(TypeConnected, ConnectedPayload{
		ClientId: clientID,
		Message:  "大屏连接成功",
	})

	go client.WritePump()
	go client.ReadPump()
}

// HandleH5Connection H5前端WebSocket连接
// @Tags WebSocket
// @Summary H5前端WebSocket连接
// @Description H5用户端连接
// @Param token query string true "用户Token"
// @Param activityId query string true "活动ID"
// @Router /ws/h5 [get]
func (h *Handler) HandleH5Connection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("WebSocket升级失败", zap.Error(err))
		return
	}

	// 验证Token
	token := c.Query("token")
	userID := h.validateToken(token)
	if userID == 0 {
		conn.WriteJSON(Message{
			Type: TypeError,
			Payload: ErrorPayload{
				Code:    401,
				Message: "未登录或Token无效",
			},
		})
		conn.Close()
		return
	}

	clientID := uuid.New().String()
	client := NewClient(clientID, userID, conn, h.hub)

	h.hub.Register <- client

	// 加入活动房间
	activityId := c.Query("activityId")
	if activityId != "" {
		// 加入弹幕房间（接收弹幕消息）
		client.JoinRoom(RoomTypeDanmaku + ":" + activityId)
		// 加入摇一摇房间（接收游戏消息）
		client.JoinRoom(RoomTypeShake + ":" + activityId)
	}

	client.SendMessage(TypeConnected, ConnectedPayload{
		ClientId: clientID,
		Message:  "连接成功",
	})

	go client.WritePump()
	go client.ReadPump()
}

// validateToken 验证Token并返回用户ID
func (h *Handler) validateToken(token string) uint {
	if token == "" {
		return 0
	}
	// 使用共享的 h5jwt 包验证 Token
	claims, err := h5jwt.ValidateToken(token)
	if err != nil {
		global.GVA_LOG.Error("Token验证失败", zap.Error(err))
		return 0
	}
	return claims.UserId

}
