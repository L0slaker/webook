package sms

import (
	"context"
	"errors"
	"fmt"
)

// Service 短信服务
type Service interface {
	Send(ctx context.Context, tplId string, args []string, phone ...string) error
}

func SendSmsError(code, msg string) error {
	return errors.New(fmt.Sprintf("send maths error: %s, msg: %s", code, msg))
}
