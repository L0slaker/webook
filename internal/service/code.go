package service

import (
	"context"
	"github.com/l0slakers/webook/internal/pkg/maths"
	"github.com/l0slakers/webook/internal/repository"
	"github.com/l0slakers/webook/internal/service/sms"
)

const codeTplId = "1877556"

type CodeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo repository.CodeRepository, sms sms.Service) *CodeService {
	return &CodeService{
		repo: repo,
		sms:  sms,
	}
}

// Send 发送验证码
// 1.如果Redis中不存在该key，则直接发送
// 2.如果Redis中存在该key，
// 2.1 但没有过期时间说明系统异常
// 2.2 如果key有过期时间，且过期时间还有剩余（具体业务为准），则认为发送太频繁，拒绝
// 2.3 否则，重新发送一个验证码
func (c *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := maths.GenerateCode()
	err := c.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	return c.sms.Send(ctx, codeTplId, []string{code}, phone)
}

// Verify 校验验证码
// 1.如果验证码不存在，说明还未发送
// 2.如果存在
// 2.1 验证次数<=3次，比较验证码是否相等
// 2.2 验证次数>3次，则认为验证码错误次数过多，直接返回false
func (c *CodeService) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.repo.Verify(ctx, biz, phone, code)
}
