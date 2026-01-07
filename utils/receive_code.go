package utils

import (
	"math/rand"
	"time"
)

// GenerateReceiveCode 生成6位核销密码（数字+大写字母）
func GenerateReceiveCode() string {
	const charset = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ" // 去掉容易混淆的 I O
	const length = 6

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[r.Intn(len(charset))]
	}
	return string(code)
}
