package dto

// WechatUserInfo 微信用户信息
type WechatUserInfo struct {
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgUrl string `json:"headimgurl"`
}

// WechatAccessToken 微信AccessToken
type WechatAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}
