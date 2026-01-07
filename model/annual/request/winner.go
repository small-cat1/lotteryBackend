package request

// WinnerReceive 确认领奖请求
type WinnerReceive struct {
	Id         uint   `json:"id" binding:"required"`
	VerifyCode string `json:"verifyCode" binding:"required"` // 核销码
	Remark     string `json:"remark"`                        // 备注（可选）
}
