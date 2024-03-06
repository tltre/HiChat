package models

import (
	"gorm.io/gorm"

	"time"
)

// Model record model implemented by authority, is consistent with gorm.Model
type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// UserBasic Basic User model
type UserBasic struct {
	Model
	Name          string
	PassWord      string
	Avatar        string // profile photo
	Gender        string `gorm:"column:gender;default:male;type:varchar(6);comment:'value in {male, female}'"`
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string `valid:"ipv4"`
	ClientPort    string
	Salt          string
	LoginTime     *time.Time `gorm:"column:login_time"`
	HeartBeatTime *time.Time `gorm:"column:heart_beat_time"`
	LoginOutTime  *time.Time `gorm:"column:login_out_time"`
	IsLoginOut    bool
	DeviceInfo    string // the device of login in
}

// UserTableName 返回用户表的名字
func (b *UserBasic) UserTableName() string {
	return "user_basic"
}
