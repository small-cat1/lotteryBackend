package request

import (
	"lotteryBackend/model/common/request"
)

// UserSearch 用户搜索
type UserSearch struct {
	request.PageInfo
	Nickname string `json:"nickname"`
}

// ===== 活动相关 =====

// ActivitySearch 活动搜索
type ActivitySearch struct {
	request.PageInfo
	Title  string `json:"title"`
	Status *int   `json:"status"`
}

// ===== 弹幕相关 =====

// DanmakuSend 发送弹幕
type DanmakuSend struct {
	ActivityId uint   `json:"activityId" binding:"required"`
	Content    string `json:"content" binding:"required,max=200"`
	Color      string `json:"color"`
}

// DanmakuSearch 弹幕搜索
type DanmakuSearch struct {
	request.PageInfo
	ActivityId uint   `json:"activityId"`
	UserId     uint   `json:"userId"`
	Content    string `json:"content"`
	Status     *int   `json:"status"`
}

// DanmakuAudit 弹幕审核
type DanmakuAudit struct {
	Ids    []uint `json:"ids" binding:"required"`
	Status int    `json:"status" binding:"required,oneof=1 2"` // 1通过 2拒绝
}

// ===== 摇一摇相关 =====

// ShakeRoundSearch 摇一摇场次搜索
type ShakeRoundSearch struct {
	request.PageInfo
	ActivityId uint   `json:"activityId"`
	RoundName  string `json:"roundName"`
	Status     *int   `json:"status"`
}

// ShakeScoreUpdate 摇一摇分数上报
type ShakeScoreUpdate struct {
	RoundId uint `json:"roundId" binding:"required"`
	Score   int  `json:"score" binding:"required,min=1"`
}

// ===== 奖品相关 =====

// PrizeSearch 奖品搜索
type PrizeSearch struct {
	request.PageInfo
	ActivityId uint   `json:"activityId"`
	Name       string `json:"name"`
	Level      *int   `json:"level"`
}

// ===== 中奖相关 =====

// WinnerSearch 中奖记录搜索
type WinnerSearch struct {
	request.PageInfo
	ActivityId uint `json:"activityId"`
	UserId     uint `json:"userId"`
	PrizeId    uint `json:"prizeId"`
	RoundId    uint `json:"roundId"`
	WinType    *int `json:"winType"`
	Status     *int `json:"status"`
}
