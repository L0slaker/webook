package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/l0slakers/webook/internal/domain"
	"github.com/l0slakers/webook/internal/repository"
)

var (
	ErrDuplicateEmail = repository.ErrDuplicateEmail
	ErrUnknownEmail   = repository.ErrUnknownEmail
	ErrWrongInfo      = errors.New("用户名或密码不正确！")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 密码加密
	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(pwd)

	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrWrongInfo
	}

	return user, nil
}

func (svc *UserService) Edit(ctx context.Context, user domain.User) error {
	return svc.repo.Update(ctx, user)
}

func (svc *UserService) Info(ctx context.Context, uid int64) (domain.User, error) {
	return svc.repo.FindByID(ctx, uid)
}
