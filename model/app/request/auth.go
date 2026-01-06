package request

// ========== 授权相关 ==========

type WechatLoginReq struct {
	Code       string `json:"code" binding:"required"`
	ActivityId uint   `json:"activityId" binding:"required"` // ✅ 可选，类型为 uint
}

type WxJsConfigReq struct {
	Url string `form:"url" binding:"required"`
}
