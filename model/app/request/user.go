package request

// ========== 用户相关 ==========

// CheckInReq 签到请求
type CheckInReq struct {
	ActivityId uint   `json:"activityId" binding:"required"` // 活动ID
	RealName   string `json:"realName" binding:"required"`   // 真实姓名
	Phone      string `json:"phone" binding:"required"`      // 手机号
	Department string `json:"department"`                    // 部门
	EmployeeNo string `json:"employeeNo"`                    // 工号
}

// GetUserInfoReq 获取用户信息请求
type GetUserInfoReq struct {
	ActivityId uint `form:"activityId"` // 活动ID
}
