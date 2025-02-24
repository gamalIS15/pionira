package models

import "time"

type AppToken struct {
	BaseModel
	TargetId  uint      `json:"target_id" gorm:"index;not null"`
	Token     string    `json:"-" gorm:"not null; type:varchar(255)"`
	Type      string    `json:"type" gorm:"index;not null;type:varchar(255)"`
	Used      bool      `json:"-" gorm:"index;not null;type:bool"`
	ExpiredAt time.Time `json:"-" gorm:"index;not null"`
}

func (a AppToken) TableName() string {
	return "app_token"
}
