package h5jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
)

// H5Claims H5用户JWT Claims
type H5Claims struct {
	UserId   uint   `json:"userId"`
	OpenId   string `json:"openId"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(user *annual.AnnualUser) (string, error) {
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
	return token.SignedString([]byte(getJwtSecret()))
}

// ValidateToken 验证Token
func ValidateToken(tokenString string) (*H5Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &H5Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(getJwtSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*H5Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetUserIdFromToken 从Token获取用户ID（便捷方法）
func GetUserIdFromToken(tokenString string) uint {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return 0
	}
	return claims.UserId
}

// 从配置表获取JWT密钥
func getJwtSecret() string {
	var config annual.AnnualConfig
	global.GVA_DB.Where("config_key = ?", "h5_jwt_secret").First(&config)
	if config.ConfigValue == "" {
		return "annual-h5-jwt-secret-key" // 默认密钥
	}
	return config.ConfigValue
}
