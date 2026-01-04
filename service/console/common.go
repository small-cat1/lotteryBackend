package console

// safeInt 安全获取int指针值
func safeInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// boolToInt bool转int
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
