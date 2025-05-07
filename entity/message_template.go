package entity

import "gorm.io/gorm"

type MessageTemplate struct {
	gorm.Model
	Name           string `gorm:"column:template_name;type:varchar(100)"`
	Subject        string `gorm:"column:subject;type:varchar(100)"`
	Message        string `gorm:"column:message;type:text"`
	IsEmailMessage bool   `gorm:"column:is_email_message;type:bool"`
	IsSMSMessage   bool   `gorm:"column:is_sms_message;type:bool"`
}

func (MessageTemplate) TableName() string {
	return "tb_message_templates"
}
