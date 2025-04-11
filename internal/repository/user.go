package repository

import (
	"context"
	"github.com/l0slakers/webook/internal/domain"
	"github.com/l0slakers/webook/internal/repository/cache"
	"github.com/l0slakers/webook/internal/repository/dao"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUnknownEmail   = dao.ErrUnknownEmail
)

type UserRepository struct {
	userDao   dao.UserDAO
	userCache cache.UserCache
}

func NewUserRepository(userDao dao.UserDAO, userCache cache.UserCache) *UserRepository {
	return &UserRepository{
		userDao:   userDao,
		userCache: userCache,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.userDao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.userDao.FirstByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return daoToDomain(u), nil
}

func (r *UserRepository) Update(ctx context.Context, user domain.User) error {
	return r.userDao.Update(ctx, dao.User{
		ID:           user.ID,
		Nickname:     user.Nickname,
		Birthday:     user.Birthday,
		Introduction: user.Introduction,
	})
}

func (r *UserRepository) FindByID(ctx context.Context, uid int64) (u domain.User, err error) {
	u, err = r.userCache.Get(ctx, uid)
	// 此时出现 error 可能有两种情况：
	// 1.没查到Key，正常查库即可
	// 2.redis崩溃，应该保护系统，不执行查库
	switch err {
	case redis.Nil:
		user, err := r.userDao.FirstByID(ctx, uid)
		if err != nil {
			return domain.User{}, err
		}
		u = daoToDomain(user)
		// 回写缓存
		if err = r.userCache.Set(ctx, u); err != nil {
			// 网络崩了，也可能是 redis 崩了
			log.Println(err)
		}
		return u, nil
	case nil:
		return u, nil
	default:
		return domain.User{}, err
	}
}

func daoToDomain(u dao.User) domain.User {
	return domain.User{
		ID:           u.ID,
		Email:        u.Email,
		Password:     u.Password,
		Nickname:     u.Nickname,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
	}
}
