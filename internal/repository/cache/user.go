package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/l0slakers/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) UserCache {
	return UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (u *UserCache) key(uid int64) string {
	return fmt.Sprintf("user:id:%d", uid)
}

func (u *UserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := u.key(uid)
	val, err := u.client.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}

	var user domain.User
	err = json.Unmarshal([]byte(val), &user)

	return user, err
}

func (u *UserCache) Set(ctx context.Context, user domain.User) error {
	key := u.key(user.ID)
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return u.client.Set(ctx, key, val, u.expiration).Err()
}
