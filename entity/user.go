package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string    `gorm:"column:username;type:varchar(100)"`
	FirstName   string    `gorm:"column:first_name;type:varchar(100)"`
	LastName    string    `gorm:"column:last_name;type:varchar(100)"`
	PhoneNumber string    `gorm:"column:phone_number;type:varchar(15)"`
	Email       string    `gorm:"column:email;type:varchar(100)"`
	Password    string    `gorm:"column:password;type:varchar(200)"`
	IsActive    bool      `gorm:"column:is_active;type:boolean"`
	UserRole    string    `gorm:"column:user_roles;type:varchar(100)"`
	Addresses   []Address `gorm:"ForeignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FileName    string    `gorm:"column:file_name;type:varchar(100)"`
	RefineryId  uint      `gorm:"column:refinery_id;type:int"`
}

func (User) TableName() string {
	return "tb_users"

}
