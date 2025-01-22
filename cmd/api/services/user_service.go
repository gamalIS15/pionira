package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"pionira/cmd/api/requests"
	"pionira/common"
	"pionira/internal/models"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (us UserService) RegisterUser(ur *requests.RegisterAuthRequest) (*models.UserModel, error) {
	//hash password
	hashPassword, err := common.HashPassword(ur.Password)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("hashing password failed")
	}

	createdUser := models.UserModel{
		FirstName: &ur.FirstName,
		LastName:  &ur.LastName,
		Email:     ur.Email,
		Password:  hashPassword,
	}

	//fmt.Println(createdUser)
	result := us.db.Create(&createdUser)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, errors.New("failed to create user")
	}

	return &createdUser, nil
}

func (us UserService) GetUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	result := us.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (us UserService) ChangeUserPassword(newPassword string, user models.UserModel) error {
	hashPassword, err := common.HashPassword(newPassword)
	if err != nil {
		fmt.Println(err)
		return errors.New("password change failed")
	}
	result := us.db.Model(user).Update("Password", hashPassword)
	if result.Error != nil {
		fmt.Println(result.Error)
		return errors.New("failed to update password")
	}
	return nil
}
