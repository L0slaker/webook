package dao

import "gorm.io/gorm"

// TODO 自动建表
func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
