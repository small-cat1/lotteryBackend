package ws

import "time"

// ==================== 基础消息结构 ====================

// Message WebSocket基础消息结构
type Message struct {
	Type    string      `json:"type"`    // 消息类型
	Payload interface{} `json:"payload"` // 消息内容
}

// ClientMessage 客户端发送的消息
type ClientMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"` // JSON字符串，根据Type解析
}

// ==================== 系统消息 ====================

// ConnectedPayload 连接成功消息
type ConnectedPayload struct {
	ClientId string `json:"clientId"`
	Message  string `json:"message"`
}

// ErrorPayload 错误消息
type ErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HeartbeatPayload 心跳消息
type HeartbeatPayload struct {
	Timestamp int64 `json:"timestamp"`
}

// ==================== 用户信息 ====================

// UserBrief 用户简要信息
type UserBrief struct {
	ID         uint   `json:"ID"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	RealName   string `json:"realName"`
	Department string `json:"department"`
}

// ==================== 签到消息 ====================

// CheckInStatsPayload 签到统计消息（单独推送统计时使用）
type CheckInStatsPayload struct {
	Total    int `json:"total"`    // 总签到数
	Pending  int `json:"pending"`  // 待审核
	Approved int `json:"approved"` // 已通过
	Rejected int `json:"rejected"` // 已拒绝
}

// ==================== 弹幕消息 ====================

// DanmakuPayload 弹幕消息
type DanmakuPayload struct {
	ID        uint      `json:"id"`
	User      UserBrief `json:"user"`
	Content   string    `json:"content"`
	Color     string    `json:"color"`
	IsTop     int       `json:"isTop"`
	CreatedAt time.Time `json:"createdAt"`
}

// TopDanmakuPayload 置顶弹幕消息
type TopDanmakuPayload struct {
	ID      uint      `json:"id"`
	User    UserBrief `json:"user"`
	Content string    `json:"content"`
	Color   string    `json:"color"`
}

// ==================== 摇一摇消息 ====================

// PrizeBrief 奖品简要信息
type PrizeBrief struct {
	ID    uint   `json:"ID"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Level int    `json:"level"`
}

// RoundInfo 场次信息
type RoundInfo struct {
	ID          uint       `json:"ID"`
	RoundName   string     `json:"roundName"`
	Duration    int        `json:"duration"`
	WinnerCount int        `json:"winnerCount"`
	Status      int        `json:"status"`
	Prize       PrizeBrief `json:"prize"`
}

// RankingItem 排名项
type RankingItem struct {
	Rank     int       `json:"rank"`
	UserId   uint      `json:"userId"`
	User     UserBrief `json:"user"`
	Score    int       `json:"score"`
	IsWinner bool      `json:"isWinner"`
}

// RankingUpdatePayload 排名更新消息
type RankingUpdatePayload struct {
	RoundId     uint          `json:"roundId"`
	PlayerCount uint          `json:"playerCount"`
	Ranking     []RankingItem `json:"ranking"`
}

// ==================== 抽奖消息 ====================

// WinnerInfo 中奖者信息
type WinnerInfo struct {
	ID        uint       `json:"id"`
	User      UserBrief  `json:"user"`
	Prize     PrizeBrief `json:"prize"`
	WinType   int        `json:"winType"`
	CreatedAt time.Time  `json:"createdAt"`
}

// ShakeScorePayload 用户上报分数消息
type ShakeScorePayload struct {
	RoundId uint `json:"roundId"`
	Score   int  `json:"score"`
}

// GameStartPayload 游戏开始消息
type GameStartPayload struct {
	RoundId   uint      `json:"roundId"`
	Round     RoundInfo `json:"round"`
	Duration  int       `json:"duration"`  // 游戏时长（秒）
	EndTime   int64     `json:"endTime"`   // 游戏结束时间戳（毫秒）
	StartTime int64     `json:"startTime"` // 游戏开始时间戳（毫秒）
}

// GameStopPayload 游戏结束消息
type GameStopPayload struct {
	RoundId uint `json:"roundId"`
}

// ==================== 消息类型常量 ====================

const (
	// 系统消息
	TypeConnected = "connected"
	TypeError     = "error"
	TypeHeartbeat = "heartbeat"
	TypePong      = "pong"
	// 签到消息
	TypeCheckInStats = "checkin_stats"
	// 弹幕消息
	TypeNewDanmaku = "new_danmaku"
	// 摇一摇消息
	TypeShakeScore    = "shake_score"    // 用户上报分数
	TypeRankingUpdate = "ranking_update" // 排名更新广播

	// 游戏控制消息
	TypeGameStart = "game_start" // 游戏开始广播
	TypeGameStop  = "game_stop"  // 游戏结束广播
)

// ==================== 房间类型常量 ====================

const (
	RoomTypeDanmaku = "danmaku" // 弹幕房间
	RoomTypeShake   = "shake"   // 摇一摇房间
	RoomTypeScreen  = "screen"  // 大屏房间（接收所有消息）
)
