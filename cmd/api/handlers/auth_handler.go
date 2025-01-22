package handlers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"pionira/cmd/api/requests"
	"pionira/cmd/api/services"
	"pionira/common"
	"pionira/internal/mailer"
	"pionira/internal/models"
)

func (h *Handler) RegisterHandler(c echo.Context) error {
	//Bind Request to Struct
	payload := new(requests.RegisterAuthRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	//validation
	validationErrors := h.ValidateBody(c, *payload)

	userService := services.NewUserService(h.DB)
	_, err := userService.GetUserByEmail(payload.Email)

	if errors.Is(err, gorm.ErrRecordNotFound) == false {
		return common.SendBadRequestResponse(c, "Email already in use")
	}

	registerUser, err := userService.RegisterUser(payload)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	mailData := mailer.EmailData{
		Subject: "Welcome to Pionira",
		Meta: struct {
			FirstName string
			LoginLink string
		}{
			FirstName: *registerUser.FirstName,
			LoginLink: "#",
		},
	}
	err = h.Mailer.Send(payload.Email, "welcome.html", mailData)
	if err != nil {
		h.Logger.Error(err)
	}
	return common.SendSuccessResponse(c, "User Registration Success", registerUser)
}

func (h *Handler) LoginHandler(c echo.Context) error {
	// bind data from client
	payload := new(requests.LoginAuthRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		c.Logger().Error(err)
		return common.SendBadRequestResponse(c, err.Error())
	}

	//validation
	validationErrors := h.ValidateBody(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	// check user if exist
	userService := services.NewUserService(h.DB)
	userRetrieved, err := userService.GetUserByEmail(payload.Email)
	fmt.Println(userRetrieved)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common.SendBadRequestResponse(c, "User Not Found")
	}

	// compare client data password to hashed password
	if common.ComparePasswordHash(payload.Password, userRetrieved.Password) == false {
		return common.SendBadRequestResponse(c, "Invalid credentials")
	}

	//generate Token
	accessToken, refreshToken, err := common.GenerateJWT(*userRetrieved)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	// send login token
	return common.SendSuccessResponse(c, "User Logged In", map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          userRetrieved,
	})
}

func (h *Handler) GetAuthenticatedUser(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendBadRequestResponse(c, "User not logged in")
	}
	return common.SendSuccessResponse(c, "Authenticated user retrieved", user)
}
