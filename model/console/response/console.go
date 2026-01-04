package response

import "time"

// ==================== 活动配置 ====================

// ActivityDetailResp 活动详情响应
type ActivityDetailResp struct {
	ID             uint       `json:"id"`
	Title          string     `json:"title"`
	Logo           string     `json:"logo"`
	Cover          string     `json:"cover"`
	Description    string     `json:"description"`
	StartTime      *time.Time `json:"startTime"`
	EndTime        *time.Time `json:"endTime"`
	CheckInEnabled int        `json:"checkInEnabled"`
	DanmakuEnabled int        `json:"danmakuEnabled"`
	DanmakuAudit   int        `json:"danmakuAudit"`
	WinnerExclude  int        `json:"winnerExclude"`
	Status         int        `json:"status"`
}

// PrizeListResp 奖品列表响应
type PrizeListResp struct {
	List []PrizeItem `json:"list"`
}

// PrizeItem 奖品信息
type PrizeItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Level       int    `json:"level"`
	TotalCount  int    `json:"totalCount"`
	RemainCount int    `json:"remainCount"`
}

// ==================== 签到管理 ====================

// CheckInStatsResp 签到统计响应
type CheckInStatsResp struct {
	IsOpen    bool    `json:"isOpen"`    // 签到是否开启
	CheckedIn int     `json:"checkedIn"` // 已签到人数
	Total     int     `json:"total"`     // 总人数
	Rate      float64 `json:"rate"`      // 签到率
}

// CheckInListResp 签到列表响应
type CheckInListResp struct {
	List     []CheckInItem `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

// CheckInItem 签到记录
type CheckInItem struct {
	ID          uint      `json:"id"`
	UserId      uint      `json:"userId"`
	CheckInTime time.Time `json:"checkInTime"`
	User        *UserInfo `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID         uint   `json:"id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	RealName   string `json:"realName"`
	Department string `json:"department"`
	Phone      string `json:"phone"`
}

// ==================== 弹幕管理 ====================

// DanmakuListResp 弹幕列表响应
type DanmakuListResp struct {
	List      []DanmakuItem `json:"list"`
	UserCount int           `json:"userCount"`
}

// DanmakuItem 弹幕信息
type DanmakuItem struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Color     string    `json:"color"`
	Status    int       `json:"status"`
	IsTop     int       `json:"isTop"`
	CreatedAt time.Time `json:"createdAt"`
	User      *UserInfo `json:"user"`
}

// ==================== 场次管理 ====================

// RoundListResp 场次列表响应
type RoundListResp struct {
	List []RoundItem `json:"list"`
}

// RoundItem 场次信息
type RoundItem struct {
	ID            uint       `json:"id"`
	RoundName     string     `json:"roundName"`
	Duration      int        `json:"duration"`
	WinnerCount   int        `json:"winnerCount"`
	Status        int        `json:"status"`
	ActualWinners int        `json:"actualWinners"`
	Prize         *PrizeItem `json:"prize"`
}

// RoundDetailResp 场次详情响应
type RoundDetailResp struct {
	ID            uint       `json:"id"`
	ActivityId    uint       `json:"activityId"`
	RoundName     string     `json:"roundName"`
	Duration      int        `json:"duration"`
	WinnerCount   int        `json:"winnerCount"`
	Status        int        `json:"status"`
	StartTime     *time.Time `json:"startTime"`
	EndTime       *time.Time `json:"endTime"`
	ActualWinners int        `json:"actualWinners"`
	PlayerCount   int        `json:"playerCount"`
	Prize         *PrizeItem `json:"prize"`
}

// ==================== 游戏控制 ====================

// GameStatusResp 游戏状态响应
type GameStatusResp struct {
	RoundId     uint       `json:"roundId"`
	Status      int        `json:"status"`      // -1未选择 0准备中 1进行中 2已结束
	PlayerCount int        `json:"playerCount"` // 参与人数
	StartTime   *time.Time `json:"startTime"`   // 开始时间
	Duration    int        `json:"duration"`    // 持续时间（秒）
	Remaining   int        `json:"remaining"`   // 剩余时间（秒）
}

// RankingListResp 排行榜响应
type RankingListResp struct {
	List        []RankingItem `json:"list"`
	PlayerCount int           `json:"playerCount"`
}

// RankingItem 排行信息
type RankingItem struct {
	Rank     int       `json:"rank"`
	UserId   uint      `json:"userId"`
	Score    int       `json:"score"`
	IsWinner bool      `json:"isWinner"`
	User     *UserInfo `json:"user"`
}

// WinnerListResp 中奖名单响应
type WinnerListResp struct {
	List []WinnerItem `json:"list"`
}

// WinnerItem 中奖信息
type WinnerItem struct {
	ID        uint       `json:"id"`
	Rank      int        `json:"rank"`
	UserId    uint       `json:"userId"`
	Score     int        `json:"score"`
	WinType   int        `json:"winType"` // 1摇一摇 2随机抽奖 3弹幕抽奖
	CreatedAt time.Time  `json:"createdAt"`
	User      *UserInfo  `json:"user"`
	Prize     *PrizeItem `json:"prize"`
}

type CurrentGameResp struct {
	Round     *RoundInfo    `json:"round"`     // 当前场次信息，nil表示无进行中游戏
	Status    int           `json:"status"`    // 状态 0待开始 1进行中 2已结束
	Remaining int           `json:"remaining"` // 剩余时间（秒）
	Ranking   []RankingItem `json:"ranking"`   // 当前排行榜
}

// RoundInfo 场次信息
type RoundInfo struct {
	ID          uint       `json:"id"`
	RoundName   string     `json:"roundName"`
	Duration    int        `json:"duration"`
	WinnerCount int        `json:"winnerCount"`
	Prize       *PrizeInfo `json:"prize"`
}

// PrizeInfo 奖品信息
type PrizeInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Level int    `json:"level"`
}

// ==================== 抽奖结果 ====================

// DrawResultResp 抽奖结果响应
type DrawResultResp struct {
	Winners []WinnerItem `json:"winners"`
	Prize   *PrizeItem   `json:"prize"`
}
