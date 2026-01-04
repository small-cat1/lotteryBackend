package utils

import (
	"math/rand"
	"time"
)

// GeneratePassword 生成随机密码
func GeneratePassword(length int) string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 排除易混淆字符 I/1/O/0

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
