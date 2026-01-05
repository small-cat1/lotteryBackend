package request

// ===== 通用请求 =====

type GetById struct {
	ID uint `json:"id" binding:"required"`
}

type GetByIds struct {
	Ids []uint `json:"ids" binding:"required"`
}

type UpdateStatus struct {
	ID     uint `json:"id" binding:"required"`
	Status int  `json:"status"`
}

// ===== 用户相关 =====

// ===== 活动相关 =====

// ===== 弹幕相关 =====

type DanmakuTop struct {
	ID    uint `json:"id" binding:"required"`
	IsTop int  `json:"isTop"`
}

// ===== 摇一摇相关 =====

// ===== 奖品相关 =====

// ===== 中奖相关 =====

type RandomDraw struct {
	ActivityId uint `json:"activityId" binding:"required"`
	PrizeId    uint `json:"prizeId" binding:"required"`
	Count      int  `json:"count" binding:"required,min=1"`
}
