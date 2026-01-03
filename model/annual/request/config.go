package request

// WechatConfigReq 微信配置请求
type WechatConfigReq struct {
	AppID  string `json:"appId" binding:"required"`
	Secret string `json:"secret" binding:"required"`
	Token  string `json:"token"`
}

// SiteConfigReq 站点配置请求
type SiteConfigReq struct {
	Name      string `json:"name" binding:"required"`
	Logo      string `json:"logo"`
	Copyright string `json:"copyright"`
}
