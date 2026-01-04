package request

// ==================== 游戏控制 ====================

// StartGameReq 开始游戏请求
type StartGameReq struct {
	RoundId  uint   `json:"roundId" binding:"required"`  // 场次ID
	Password string `json:"password" binding:"required"` // 场次密码
}

// StopGameReq 停止游戏请求
type StopGameReq struct {
	RoundId uint `json:"roundId" binding:"required"` // 场次ID
}

// SettleGameReq 结算游戏请求
type SettleGameReq struct {
	RoundId uint `json:"roundId" binding:"required"` // 场次ID
}

// ==================== 随机抽奖 ====================

// RandomDrawReq 随机抽奖请求
type RandomDrawReq struct {
	ActivityId uint `json:"activityId" binding:"required"` // 活动ID
	Count      int  `json:"count" binding:"required"`      // 抽取人数
}

// ==================== H5摇一摇 ====================

// JoinShakeReq 加入游戏请求
type JoinShakeReq struct {
	RoundId uint `json:"roundId" binding:"required"` // 场次ID
}

// SubmitShakeScoreReq 提交分数请求
type SubmitShakeScoreReq struct {
	RoundId uint `json:"roundId" binding:"required"` // 场次ID
	Score   int  `json:"score" binding:"required"`   // 分数
}
