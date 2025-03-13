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
	return r.userDao.Insert(ctx, domainToDAO(u))
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.userDao.FirstByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return daoToDomain(u), nil
}

func domainToDAO(u domain.User) dao.User {
	return dao.User{
		Email:    u.Email,
		Password: u.Password,
	}
}

func daoToDomain(u dao.User) domain.User {
	return domain.User{
		Email:    u.Email,
		Password: u.Password,
	}
}
