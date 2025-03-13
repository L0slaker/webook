package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"

	"github.com/l0slakers/webook/internal/domain"
	"github.com/l0slakers/webook/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 密码加密
	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(pwd)

	return svc.repo.CreateUser(ctx, u)
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}
