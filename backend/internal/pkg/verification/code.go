package verification

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateCode 生成指定长度的数字验证码
func GenerateCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("验证码长度必须大于0")
	}

	code := ""
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("生成验证码失败: %w", err)
		}
		code += num.String()
	}

	return code, nil
}

// GenerateVerificationCode 生成6位验证码
func GenerateVerificationCode() (string, error) {
	return GenerateCode(6)
}
