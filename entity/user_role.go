package entity

import "github.com/google/uuid"

type UserRole struct {
	Id     uuid.UUID `gorm:"primaryKey;column:user_role_id;type:varchar(36)"`
	Role   string    `gorm:"column:rAole;type:varchar(10)"`
	UserID uint      `gorm:"column:user_id;type:int"`
}

func (UserRole) TableName() string {
	return "tb_user_role"
}
