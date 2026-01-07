package response

import "time"

// ========== 授权相关 ==========

// WechatLoginResp 微信登录响应
type WechatLoginResp struct {
	Token string     `json:"token"` // JWT Token
	User  H5UserResp `json:"user"`  // 用户信息
}

// WxJsConfigResp 微信JS-SDK配置响应
type WxJsConfigResp struct {
	AppId     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

// WechatConfigResp 微信配置响应
type WechatConfigResp struct {
	AppId string `json:"appId"` // 微信AppID
}

// ========== 用户相关 ==========

// ========== 活动相关 ==========

// H5ActivityResp 活动信息响应
type H5ActivityResp struct {
	ID             uint      `json:"ID"`
	Title          string    `json:"title"`
	Cover          string    `json:"cover"`
	Description    string    `json:"description"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	Status         int       `json:"status"` // 0未开始 1进行中 2已结束
	CheckInEnabled int       `json:"checkInEnabled"`
	DanmakuEnabled int       `json:"danmakuEnabled"`
	DanmakuAudit   int       `json:"danmakuAudit"`
}

// ========== 签到相关 ==========

// CheckInStatusResp 签到状态响应
type CheckInStatusResp struct {
	IsCheckedIn bool       `json:"isCheckedIn"` // 是否已签到
	CheckInTime *time.Time `json:"checkInTime"` // 签到时间
}

// CheckInStatsResp 签到统计响应
type CheckInStatsResp struct {
	CheckedCount int     `json:"checkedCount"` // 已签到人数
	TotalCount   int     `json:"totalCount"`   // 总人数
	CheckRate    float64 `json:"checkRate"`    // 签到率
}

// ========== 弹幕相关 ==========

// DanmakuResp 弹幕响应
type DanmakuResp struct {
	ID        uint        `json:"id"`
	User      H5UserBrief `json:"user"`
	Content   string      `json:"content"`
	Color     string      `json:"color"`
	IsTop     int         `json:"isTop"`
	CreatedAt time.Time   `json:"createdAt"`
}

// ========== 摇一摇相关 ==========

// ShakeRoundResp 摇一摇场次响应
type ShakeRoundResp struct {
	ID          uint          `json:"ID"`
	RoundName   string        `json:"roundName"`
	Duration    int           `json:"duration"`    // 游戏时长（秒）
	WinnerCount int           `json:"winnerCount"` // 获奖人数
	Status      int           `json:"status"`      // 0未开始 1进行中 2已结束
	Prize       *H5PrizeBrief `json:"prize"`       // 关联奖品
	StartTime   *time.Time    `json:"startTime"`
	EndTime     *time.Time    `json:"endTime"`
	EndTimeMs   int64         `json:"endTimeMs"` // 结束时间戳（毫秒）
	MyScore     int           `json:"myScore"`   // ⭐ 新增
}

// ShakeRankingResp 摇一摇排名响应
type ShakeRankingResp struct {
	UserId   uint        `json:"userId"`
	User     H5UserBrief `json:"user"`
	Score    int         `json:"score"`
	Rank     int         `json:"rank"`
	IsWinner bool        `json:"isWinner"`
}

// MyScoreResp 我的成绩响应
type MyScoreResp struct {
	Score    int  `json:"score"`
	Rank     int  `json:"rank"`
	IsWinner bool `json:"isWinner"`
}

// RoundResultResp 场次结果响应
type RoundResultResp struct {
	Ranking  []ShakeRankingResp `json:"ranking"`
	MyRank   int                `json:"myRank"`
	IsWinner bool               `json:"isWinner"`
	WinInfo  *WinningResp       `json:"winInfo"`
}

// ========== 奖品相关 ==========

// H5PrizeResp 奖品响应
type H5PrizeResp struct {
	ID          uint   `json:"ID"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Level       int    `json:"level"` // 1特等奖 2一等奖 3二等奖 4三等奖 5参与奖
	TotalCount  int    `json:"totalCount"`
	RemainCount int    `json:"remainCount"`
}

// H5PrizeBrief 奖品简要信息
type H5PrizeBrief struct {
	ID    uint   `json:"ID"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Level int    `json:"level"`
}

// WinningResp 中奖记录响应
type WinningResp struct {
	ID          uint          `json:"id"`
	User        *H5UserBrief  `json:"user"`
	Prize       *H5PrizeBrief `json:"prize"`
	PrizeId     uint          `json:"prizeId"`
	WinType     int           `json:"winType"` // 1摇一摇 2随机抽奖 3弹幕抽奖
	Status      int           `json:"status"`  // 0未领取 1已领取
	ReceiveTime *time.Time    `json:"receiveTime"`
	CreatedAt   time.Time     `json:"createdAt"`
	ReceiveCode string        `json:"receiveCode"` // ✅ 新增
}
