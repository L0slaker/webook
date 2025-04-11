package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string

	ErrSendTooMany       = errors.New("发送验证码过于频繁")
	ErrCodeVerifyTooMany = errors.New("验证次数过多")
)

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{
		client: client,
	}
}

func (c *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		// 调用 redis 时出现问题
		return err
	}
	switch res {
	case -2:
		return errors.New("验证码存在，但没有过期时间")
	case -1:
		return ErrSendTooMany
	default:
		// 发送成功
		return nil
	}
}

func (c *CodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		// 调用 redis 时出现问题
		return false, err
	}
	switch res {
	case -2:
		return false, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		// 发送成功
		return true, nil
	}
}

func (c *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone-maths:%s-%s", biz, phone)
}
