package request

// ========== 弹幕相关 ==========

// SendDanmakuReq 发送弹幕请求
type SendDanmakuReq struct {
	ActivityId uint   `json:"activityId" binding:"required"` // 活动ID
	Content    string `json:"content" binding:"required"`    // 弹幕内容
	Color      string `json:"color"`                         // 弹幕颜色
}

// ========== 摇一摇相关 ==========

// JoinGameReq 加入游戏请求
type JoinGameReq struct {
	RoundId uint `json:"roundId" binding:"required"` // 场次ID
}

// SubmitScoreReq 提交分数请求
type SubmitScoreReq struct {
	RoundId uint `json:"roundId" binding:"required"` // 场次ID
	Score   int  `json:"score" binding:"required"`   // 分数（摇动次数）
}
