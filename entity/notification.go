package entity

import (
	"gorm.io/gorm"
	"time"
)

// MessageType represents the type of message (email or SMS)
type MessageType string

const (
	EMAIL MessageType = "email"
	SMS   MessageType = "sms"
)

// Notification represents a record of a message sent to a user
type Notification struct {
	gorm.Model
	UserID    uint        `gorm:"column:user_id;index"`
	Type      MessageType `gorm:"column:type;type:varchar(10)"`
	Recipient string      `gorm:"column:recipient;type:varchar(255)"`
	Subject   string      `gorm:"column:subject;type:varchar(255)"`
	Content   string      `gorm:"column:content;type:text"`
	Status    string      `gorm:"column:status;type:varchar(20);default:'sent'"` // sent, failed, delivered
	SentAt    time.Time   `gorm:"column:sent_at"`
}
