package main

import (
	"log"
	"pionira/common"
	"pionira/internal/models"
)

func main() {
	db, err := common.NewMysql()
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		&models.UserModel{},
		&models.AppToken{},
		&models.CategoryModel{})

	if err != nil {
		panic(err)
	}
	log.Println("Database migration completed")

}
