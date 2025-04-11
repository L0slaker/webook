package maths

import (
	"fmt"
	"math/rand"
)

// GenerateCode 生成验证码(6位随机数字，不足零的补零)
func GenerateCode() string {
	n := rand.Intn(1000000)
	return fmt.Sprintf("%06d", n)
}
