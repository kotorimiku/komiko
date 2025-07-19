package utils

import (
	"regexp"
)

// IsValidPassword 检查密码是否为长度 6~20，且仅包含字母、数字、.-+_
func IsValidPassword(password string) bool {
	if len(password) <= 6 || len(password) >= 20 {
		return false
	}

	// 正则：仅允许 a-zA-Z0-9 ._-+
	regex := `^[a-zA-Z0-9._\-+]*$`
	match, _ := regexp.MatchString(regex, password)
	return match
}
