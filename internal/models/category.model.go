package models

type CategoryModel struct {
	Name     string `json:"name" gorm:"unique; not null;type:varchar(200)"`
	Slug     string `json:"slug" gorm:"not null; varchar(200)"`
	IsCustom bool   `json:"is_custom" gorm:"not null; default:false"`
}

func (c CategoryModel) TableName() string {
	return "categories"
}
