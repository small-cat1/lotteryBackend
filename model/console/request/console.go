package request

// ==================== 活动配置 ====================

// ActivityIdReq 活动ID请求
type ActivityIdReq struct {
	ActivityId uint `json:"activityId" form:"activityId" binding:"required"` // 活动ID
}

// ==================== 签到管理 ====================

// OpenCheckInReq 开启签到请求
type OpenCheckInReq struct {
	ActivityId uint `json:"activityId" binding:"required"` // 活动ID
}

// CloseCheckInReq 关闭签到请求
type CloseCheckInReq struct {
	ActivityId uint `json:"activityId" binding:"required"` // 活动ID
}

// CheckInListReq 签到列表请求
type CheckInListReq struct {
	ActivityId uint   `json:"activityId" form:"activityId" binding:"required"` // 活动ID
	Page       int    `json:"page" form:"page"`                                // 页码
	PageSize   int    `json:"pageSize" form:"pageSize"`                        // 每页数量
	Keyword    string `json:"keyword" form:"keyword"`                          // 搜索关键词
}

// ==================== 弹幕管理 ====================

// DanmakuListReq 弹幕列表请求
type DanmakuListReq struct {
	ActivityId uint `json:"activityId" form:"activityId" binding:"required"` // 活动ID
	Limit      int  `json:"limit" form:"limit"`                              // 数量限制
	Status     int  `json:"status" form:"status"`                            // 状态筛选
}

// AuditDanmakuReq 审核弹幕请求
type AuditDanmakuReq struct {
	DanmakuId uint `json:"danmakuId" binding:"required"` // 弹幕ID
	Status    int  `json:"status" binding:"required"`    // 状态：1通过 2拒绝
}

// DanmakuDrawReq 弹幕抽奖请求
type DanmakuDrawReq struct {
	ActivityId uint `json:"activityId" binding:"required"` // 活动ID
	Count      int  `json:"count" binding:"required"`      // 抽取人数
	PrizeId    uint `json:"prizeId"`                       // 奖品ID（可选）
}

// ==================== 游戏控制 ====================

// PrepareGameReq 准备游戏请求
type PrepareGameReq struct {
	RoundId uint `json:"roundId" binding:"required"` // 场次ID
}

// StartGameReq 开始游戏请求
type StartGameReq struct {
	RoundId  uint   `json:"roundId" binding:"required"`  // 场次ID
	Password string `json:"password" binding:"required"` // 场次密码
}

// StopGameReq 停止游戏请求
type StopGameReq struct {
	RoundId uint `json:"roundId" binding:"required"` // 场次ID
}

// CancelGameReq 取消游戏请求
type CancelGameReq struct {
	RoundId uint `json:"roundId" binding:"required"` // 场次ID
}

// GameStatusReq 获取游戏状态请求
type GameStatusReq struct {
	RoundId uint `json:"roundId" form:"roundId" binding:"required"` // 场次ID
}

// RankingReq 获取排行榜请求
type RankingReq struct {
	RoundId uint `json:"roundId" form:"roundId" binding:"required"` // 场次ID
	Limit   int  `json:"limit" form:"limit"`                        // 数量限制
}

// WinnersReq 获取中奖名单请求
type WinnersReq struct {
	RoundId uint `json:"roundId" form:"roundId" binding:"required"` // 场次ID
}

// ==================== 随机抽奖 ====================

// ConsoleRandomDrawReq 控制台随机抽奖请求
type ConsoleRandomDrawReq struct {
	ActivityId uint `json:"activityId" binding:"required"` // 活动ID
	Count      int  `json:"count" binding:"required"`      // 抽取人数
	PrizeId    uint `json:"prizeId"`                       // 奖品ID（可选）
}
