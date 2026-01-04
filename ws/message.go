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

// CheckInPayload 签到消息
type CheckInPayload struct {
	ID          uint      `json:"id"`
	User        UserBrief `json:"user"`
	CheckInTime time.Time `json:"checkInTime"`
}

// CheckInStatsPayload 签到统计消息
type CheckInStatsPayload struct {
	CheckedCount int     `json:"checkedCount"`
	TotalCount   int     `json:"totalCount"`
	CheckRate    float64 `json:"checkRate"`
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

// RoundStartPayload 场次开始消息
type RoundStartPayload struct {
	Round RoundInfo `json:"round"`
}

// RoundEndPayload 场次结束消息
type RoundEndPayload struct {
	RoundId uint            `json:"roundId"`
	Ranking []RankingItem   `json:"ranking"`
	Winners []WinnerInfo    `json:"winners"`
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
	RoundId uint          `json:"roundId"`
	Ranking []RankingItem `json:"ranking"`
}

// CountdownPayload 倒计时消息
type CountdownPayload struct {
	RoundId    uint `json:"roundId"`
	RemainTime int  `json:"remainTime"`
}

// GameReadyPayload 游戏准备消息
type GameReadyPayload struct {
	Round     RoundInfo `json:"round"`
	Countdown int       `json:"countdown"` // 准备倒计时秒数
}

// ==================== 抽奖消息 ====================

// DrawStartPayload 抽奖开始消息
type DrawStartPayload struct {
	Prize      PrizeBrief  `json:"prize"`
	DrawCount  int         `json:"drawCount"`  // 本次抽取人数
	Candidates []UserBrief `json:"candidates"` // 候选人列表（用于滚动动画）
}

// WinnerInfo 中奖者信息
type WinnerInfo struct {
	ID        uint       `json:"id"`
	User      UserBrief  `json:"user"`
	Prize     PrizeBrief `json:"prize"`
	WinType   int        `json:"winType"`
	CreatedAt time.Time  `json:"createdAt"`
}

// DrawResultPayload 抽奖结果消息
type DrawResultPayload struct {
	Prize   PrizeBrief   `json:"prize"`
	Winners []WinnerInfo `json:"winners"`
}

// RollingUpdatePayload 滚动动画更新
type RollingUpdatePayload struct {
	Users []UserBrief `json:"users"`
}

// DrawResetPayload 重置抽奖
type DrawResetPayload struct {
	Message string `json:"message"`
}

// ==================== 消息类型常量 ====================

const (
	// 系统消息
	TypeConnected = "connected"
	TypeError     = "error"
	TypeHeartbeat = "heartbeat"
	TypePong      = "pong"

	// 签到消息
	TypeNewCheckIn   = "new_checkin"
	TypeCheckInStats = "checkin_stats"

	// 弹幕消息
	TypeNewDanmaku = "new_danmaku"
	TypeTopDanmaku = "top_danmaku"

	// 摇一摇消息
	TypeGameReady     = "game_ready"
	TypeRoundStart    = "round_start"
	TypeRoundEnd      = "round_end"
	TypeRankingUpdate = "ranking_update"
	TypeCountdown     = "countdown"

	// 抽奖消息
	TypeDrawStart   = "draw_start"
	TypeDrawResult  = "draw_result"
	TypeRollingUpdate = "rolling_update"
	TypeDrawReset   = "draw_reset"
)

// ==================== 房间类型常量 ====================

const (
	RoomTypeCheckIn = "checkin" // 签到房间
	RoomTypeDanmaku = "danmaku" // 弹幕房间
	RoomTypeShake   = "shake"   // 摇一摇房间
	RoomTypeDraw    = "draw"    // 抽奖房间
	RoomTypeScreen  = "screen"  // 大屏房间（接收所有消息）
)
