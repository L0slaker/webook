package repository

import (
	"context"
	"github.com/l0slakers/webook/internal/domain"
	"github.com/l0slakers/webook/internal/repository/dao"
)

type UserRepository struct {
	userDao *dao.UserDAO
}

func NewUserService(userDao *dao.UserDAO) *UserRepository {
	return &UserRepository{userDao: userDao}
}

func (r *UserRepository) CreateUser(ctx context.Context, u domain.User) error {
	return r.userDao.Insert(ctx, domainToDAOUser(u))
}

func domainToDAOUser(u domain.User) dao.User {
	//now := time.Now()
	return dao.User{
		//ID:        0,
		Email:    u.Email,
		Password: u.Password,
		//CreatedAt: now,
		//UpdatedAt: now,
	}
}
