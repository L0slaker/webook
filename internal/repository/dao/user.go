package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

const (
	uniqueIndexErr uint16 = 1062 // 唯一索引冲突
)

var (
	ErrDuplicateEmail = errors.New("该邮箱已被注册！")
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

// TODO 未设置索引
// TODO 未设置数据存储长度
type User struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Email    string `gorm:"column:email;unique"`
	Password string `gorm:"column:password;"`
	// 不使用time.Time，避免时区问题，统一使用 UTC 0毫秒数存储
	CreatedAt int64 `gorm:"column:created_at;"`
	UpdatedAt int64 `gorm:"column:updated_at;"`
}

func (User) TableName() string {
	return "user"
}

func (u *UserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	// 唯一索引冲突
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		if err.(*mysql.MySQLError).Number == uniqueIndexErr {
			return ErrDuplicateEmail
		}
	}
	return nil
}
