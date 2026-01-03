package dto

// WechatConfig 微信配置结构
type WechatConfig struct {
	AppID  string `json:"appId"`
	Secret string `json:"secret"`
	Token  string `json:"token"`
}

// SiteConfig 站点配置结构
type SiteConfig struct {
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	Copyright string `json:"copyright"`
}

// AllConfig 所有配置
type AllConfig struct {
	Wechat WechatConfig `json:"wechat"`
	Site   SiteConfig   `json:"site"`
}
