package request

// CreateShakeRoundReq 创建摇一摇场次请求
type CreateShakeRoundReq struct {
	ActivityId  uint   `json:"activityId" binding:"required,gt=0"`
	RoundName   string `json:"roundName" binding:"required,min=1,max=64"`
	Duration    int    `json:"duration" binding:"required,min=10,max=120"`
	WinnerCount int    `json:"winnerCount" binding:"required,min=1,max=100"`
	PrizeId     uint   `json:"prizeId" binding:"required,gt=0"`
	Sort        int    `json:"sort" binding:"min=0,max=999"`
}

// UpdateShakeRoundReq 更新摇一摇场次请求
type UpdateShakeRoundReq struct {
	ID          uint   `json:"ID" binding:"required,gt=0"`
	ActivityId  uint   `json:"activityId" binding:"required,gt=0"`
	RoundName   string `json:"roundName" binding:"required,min=1,max=64"`
	Duration    int    `json:"duration" binding:"required,min=10,max=120"`
	WinnerCount int    `json:"winnerCount" binding:"required,min=1,max=100"`
	PrizeId     uint   `json:"prizeId" binding:"required,gt=0"`
	Sort        int    `json:"sort" binding:"min=0,max=999"`
}
