package models

type UserModel struct {
	BaseModel
	FirstName *string `gorm:"type:varchar(200)" json:"first_name"`
	LastName  *string `gorm:"type:varchar(200)" json:"last_name"`
	Email     string  `gorm:"type:varchar(200)" json:"email"`
	Gender    *string `gorm:"type:varchar(50)" json:"gender"`
	Password  string  `gorm:"type:varchar(200);not null" json:"-"`
}

func (r UserModel) TableName() string {
	return "users"
}
