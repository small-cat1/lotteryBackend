package app

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/app/response"
	"lotteryBackend/utils"
	"net/http"
	"sort"
	"strings"
	"time"
)

type H5AuthService struct{}

// H5Claims H5用户JWT Claims
type H5Claims struct {
	UserId   uint   `json:"userId"`
	OpenId   string `json:"openId"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

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

// WechatLogin 微信登录
func (s *H5AuthService) WechatLogin(code string) (*response.WechatLoginResp, error) {
	// 1. 获取微信配置
	appId := s.getWechatAppId()
	appSecret := s.getWechatAppSecret()

	if appId == "" || appSecret == "" {
		return nil, errors.New("微信配置未设置")
	}

	// 2. 通过code换取access_token
	tokenUrl := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		appId, appSecret, code,
	)

	tokenResp, err := http.Get(tokenUrl)
	if err != nil {
		return nil, fmt.Errorf("获取access_token失败: %v", err)
	}
	defer tokenResp.Body.Close()

	body, _ := io.ReadAll(tokenResp.Body)
	var accessToken WechatAccessToken
	if err := json.Unmarshal(body, &accessToken); err != nil {
		return nil, fmt.Errorf("解析access_token失败: %v", err)
	}

	if accessToken.ErrCode != 0 {
		return nil, fmt.Errorf("微信授权失败: %s", accessToken.ErrMsg)
	}

	// 3. 获取用户信息
	userInfoUrl := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken.AccessToken, accessToken.OpenId,
	)

	userResp, err := http.Get(userInfoUrl)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	defer userResp.Body.Close()

	body, _ = io.ReadAll(userResp.Body)
	var wxUser WechatUserInfo
	if err := json.Unmarshal(body, &wxUser); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v", err)
	}

	// 4. 查找或创建用户
	var user annual.AnnualUser
	result := global.GVA_DB.Where("open_id = ?", wxUser.OpenId).First(&user)

	if result.RowsAffected == 0 {
		// 创建新用户
		user = annual.AnnualUser{
			OpenId:   wxUser.OpenId,
			UnionId:  wxUser.UnionId,
			Nickname: wxUser.Nickname,
			Avatar:   wxUser.HeadImgUrl,
		}
		if err := global.GVA_DB.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("创建用户失败: %v", err)
		}
	} else {
		// 更新用户信息
		global.GVA_DB.Model(&user).Updates(map[string]interface{}{
			"nickname": wxUser.Nickname,
			"avatar":   wxUser.HeadImgUrl,
		})
	}

	// 5. 生成JWT Token
	token, err := s.GenerateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %v", err)
	}

	// 6. 构造响应
	isRegistered := 0
	if user.RealName != "" {
		isRegistered = 1
	}

	return &response.WechatLoginResp{
		Token: token,
		User: response.H5UserResp{
			ID:           user.ID,
			OpenId:       user.OpenId,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			RealName:     user.RealName,
			Phone:        user.Phone,
			Department:   user.Department,
			EmployeeNo:   user.EmployeeNo,
			IsRegistered: isRegistered,
			Status:       *user.Status,
		},
	}, nil
}

// GenerateToken 生成JWT Token
func (s *H5AuthService) GenerateToken(user *annual.AnnualUser) (string, error) {
	claims := H5Claims{
		UserId:   user.ID,
		OpenId:   user.OpenId,
		Nickname: user.Nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7天过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "annual-h5",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.getJwtSecret()))
}

// ValidateToken 验证Token
func (s *H5AuthService) ValidateToken(tokenString string) (*H5Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &H5Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.getJwtSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*H5Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新Token
func (s *H5AuthService) RefreshToken(userId uint) (string, error) {
	var user annual.AnnualUser
	if err := global.GVA_DB.First(&user, userId).Error; err != nil {
		return "", err
	}
	return s.GenerateToken(&user)
}

// GetWxJsConfig 获取微信JS-SDK配置
func (s *H5AuthService) GetWxJsConfig(url string) (*response.WxJsConfigResp, error) {
	// 获取jsapi_ticket
	ticket, err := s.getJsApiTicket()
	if err != nil {
		return nil, err
	}

	// 生成签名
	nonceStr := utils.RandomString(16)
	timestamp := time.Now().Unix()

	// 拼接签名字符串
	signStr := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket, nonceStr, timestamp, url)

	// SHA1签名
	h := sha1.New()
	h.Write([]byte(signStr))
	signature := hex.EncodeToString(h.Sum(nil))

	return &response.WxJsConfigResp{
		AppId:     s.getWechatAppId(),
		Timestamp: timestamp,
		NonceStr:  nonceStr,
		Signature: signature,
	}, nil
}

// GetWechatConfig 获取微信配置（AppID）
func (s *H5AuthService) GetWechatConfig() (*response.WechatConfigResp, error) {
	appId := s.getWechatAppId()
	if appId == "" {
		return nil, errors.New("微信配置未设置")
	}
	return &response.WechatConfigResp{
		AppId: appId,
	}, nil
}

// getJsApiTicket 获取jsapi_ticket（应该缓存）
func (s *H5AuthService) getJsApiTicket() (string, error) {
	// TODO: 从缓存获取，如果没有则请求微信接口
	// 这里简化处理，实际应该使用Redis缓存

	accessToken, err := s.getAccessToken()
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Ticket  string `json:"ticket"`
	}
	json.Unmarshal(body, &result)

	if result.ErrCode != 0 {
		return "", fmt.Errorf("获取jsapi_ticket失败: %s", result.ErrMsg)
	}

	return result.Ticket, nil
}

// getAccessToken 获取access_token（应该缓存）
func (s *H5AuthService) getAccessToken() (string, error) {
	// TODO: 从缓存获取，如果没有则请求微信接口
	appId := s.getWechatAppId()
	appSecret := s.getWechatAppSecret()

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		appId, appSecret)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	json.Unmarshal(body, &result)

	if result.ErrCode != 0 {
		return "", fmt.Errorf("获取access_token失败: %s", result.ErrMsg)
	}

	return result.AccessToken, nil
}

// 从配置表获取微信AppID
func (s *H5AuthService) getWechatAppId() string {
	var config annual.AnnualConfig
	global.GVA_DB.Where("config_key = ?", "wechat_app_id").First(&config)
	return config.ConfigValue
}

// 从配置表获取微信AppSecret
func (s *H5AuthService) getWechatAppSecret() string {
	var config annual.AnnualConfig
	global.GVA_DB.Where("config_key = ?", "wechat_app_secret").First(&config)
	return config.ConfigValue
}

// 从配置表获取JWT密钥
func (s *H5AuthService) getJwtSecret() string {
	var config annual.AnnualConfig
	global.GVA_DB.Where("config_key = ?", "h5_jwt_secret").First(&config)
	if config.ConfigValue == "" {
		return "annual-h5-jwt-secret-key" // 默认密钥
	}
	return config.ConfigValue
}

// GenerateSignature 生成签名（用于其他地方）
func GenerateSignature(params map[string]string, secret string) string {
	// 按key排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接字符串
	var builder strings.Builder
	for _, k := range keys {
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(params[k])
		builder.WriteString("&")
	}
	builder.WriteString("key=")
	builder.WriteString(secret)

	// SHA1签名
	h := sha1.New()
	h.Write([]byte(builder.String()))
	return hex.EncodeToString(h.Sum(nil))
}
