package request

// ========== 授权相关 ==========

type WechatLoginReq struct {
	Code string `json:"code" binding:"required"`
}

type WxJsConfigReq struct {
	Url string `form:"url" binding:"required"`
}
