package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// Encrypt 加密明文密码，返回哈希字符串
func EncryptPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Matches 比较明文密码和哈希是否一致
func PasswordMatches(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
