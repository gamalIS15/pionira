package services

import (
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"pionira/internal/models"
	"strconv"
	"time"
)

type AppTokenService struct {
	db *gorm.DB
}

func NewAppTokenService(db *gorm.DB) *AppTokenService {
	return &AppTokenService{db: db}
}

func (a *AppTokenService) getToken() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(99999-10000) + 10000
}

func (a *AppTokenService) GeneratePasswordToken(user models.UserModel) (*models.AppToken, error) {
	tokenCreate := models.AppToken{
		TargetId:  user.ID,
		Type:      "reset",
		Token:     strconv.Itoa(a.getToken()),
		Used:      false,
		ExpiredAt: time.Now().Add(time.Hour * 24),
	}

	result := a.db.Create(&tokenCreate)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tokenCreate, nil
}

func (a *AppTokenService) ValidateResetPasswordToken(user models.UserModel, token string) (*models.AppToken, error) {
	var retrievedToken models.AppToken

	result := a.db.Where(&models.AppToken{
		TargetId: user.ID,
		Type:     "reset",
		Token:    token,
	}).First(&retrievedToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Invalid token")
		}
		return nil, result.Error
	}

	if retrievedToken.Used {
		return nil, errors.New("Token already used")
	}

	if retrievedToken.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("Token expired, please refresh the token")
	}

	return &retrievedToken, nil
}

func (a *AppTokenService) InValidateToken(userId uint, appT models.AppToken) {
	a.db.Model(&models.AppToken{}).Where("target_id = ? AND token = ?", userId, appT.Token).Update("used", true)
}
