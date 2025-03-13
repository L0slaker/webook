package repository

import (
	"context"
	"github.com/l0slakers/webook/internal/domain"
	"github.com/l0slakers/webook/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUnknownEmail   = dao.ErrUnknownEmail
)

type UserRepository struct {
	userDao *dao.UserDAO
}

func NewUserService(userDao *dao.UserDAO) *UserRepository {
	return &UserRepository{userDao: userDao}
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

func (r *UserRepository) FindByID(ctx context.Context, uid int64) (domain.User, error) {
	u, err := r.userDao.FirstByID(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	return daoToDomain(u), nil
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
