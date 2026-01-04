package request

// ========== 授权相关 ==========

// WechatLoginReq 微信登录请求
type WechatLoginReq struct {
	Code string `json:"code" binding:"required"` // 微信授权code
}

// WxJsConfigReq 获取微信JS-SDK配置请求
type WxJsConfigReq struct {
	Url string `form:"url" binding:"required"` // 当前页面URL
}

// ========== 用户相关 ==========

// UserRegisterReq 用户报名请求
type UserRegisterReq struct {
	RealName   string `json:"realName" binding:"required"` // 真实姓名
	Phone      string `json:"phone" binding:"required"`    // 手机号
	Department string `json:"department"`                  // 部门
	EmployeeNo string `json:"employeeNo"`                  // 工号
}

// UpdateUserInfoReq 更新用户信息请求
type UpdateUserInfoReq struct {
	RealName   string `json:"realName"`   // 真实姓名
	Phone      string `json:"phone"`      // 手机号
	Department string `json:"department"` // 部门
	EmployeeNo string `json:"employeeNo"` // 工号
}

// ========== 签到相关 ==========

// CheckInReq 签到请求
type CheckInReq struct {
	ActivityId uint `json:"activityId" binding:"required"` // 活动ID
}

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

// ========== 通用查询参数 ==========

// H5PageReq 分页请求
type H5PageReq struct {
	Page     int `form:"page" json:"page"`         // 页码
	PageSize int `form:"pageSize" json:"pageSize"` // 每页数量
}

// H5LimitReq 限制数量请求
type H5LimitReq struct {
	Limit int `form:"limit" json:"limit"` // 数量限制
}
