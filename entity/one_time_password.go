package entity

import (
	"gorm.io/gorm"
	"time"
)

type OneTimePassword struct {
	gorm.Model
	UserId    uint      `gorm:"column:user_id;type:int"`
	Code      string    `gorm:"column:code;type:varchar(10)"`
	IsUsed    bool      `gorm:"column:is_used;type:boolean"`
	ExpiredAt time.Time `gorm:"column:expired_at;type:timestamp"`
}

func (OneTimePassword) TableName() string {
	return "tb_one_time_passwords"
}
