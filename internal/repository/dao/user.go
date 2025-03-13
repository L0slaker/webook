package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
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
	return u.db.WithContext(ctx).Create(&user).Error
}
