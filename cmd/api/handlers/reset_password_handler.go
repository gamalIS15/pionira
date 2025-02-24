package handlers

import (
	"encoding/base64"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/url"
	"pionira/cmd/api/requests"
	"pionira/cmd/api/services"
	"pionira/common"
	"pionira/internal/mailer"
)

func (h *Handler) ForgotPasswordHandler(c echo.Context) error {
	payload := new(requests.ForgotPasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	//validation
	validationError := h.ValidateBody(c, *payload)

	if validationError != nil {
		return common.SendFailedValidationResponse(c, validationError)
	}

	userService := services.NewUserService(h.DB)
	appTokenService := services.NewAppTokenService(h.DB)

	retrievedUser, err := userService.GetUserByEmail(payload.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendNotFoundResponse(c, "Record Not Found")
		}

		return common.SendInternalServerErrorResponse(c, "Error occurred while fetching user")
	}

	token, err := appTokenService.GeneratePasswordToken(*retrievedUser)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Error occurred while generating token")
	}

	encodedEmail := base64.RawURLEncoding.EncodeToString([]byte(retrievedUser.Email))

	frontendUrl, err := url.Parse(payload.FrontEndUrl)
	if err != nil {
		return common.SendBadRequestResponse(c, "Error occurred while parsing frontend url")
	}

	query := url.Values{}
	query.Set("email", encodedEmail)
	query.Set("token", token.Token)

	frontendUrl.RawQuery = query.Encode()

	mailData := mailer.EmailData{
		Subject: "Request Password Reset",
		Meta: struct {
			Token       string
			FrontendUrl string
		}{
			Token:       token.Token,
			FrontendUrl: frontendUrl.String(),
		},
	}
	err = h.Mailer.Send(payload.Email, "forgot-password.html", mailData)
	if err != nil {
		h.Logger.Error(err)
	}

	return common.SendSuccessResponse(c, "Forgot password email sent", nil)
}

func (h *Handler) ResetPasswordHandler(c echo.Context) error {
	payload := new(requests.ResetPasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	//validation
	validationError := h.ValidateBody(c, *payload)

	if validationError != nil {
		return common.SendFailedValidationResponse(c, validationError)
	}

	email, err := base64.RawURLEncoding.DecodeString(payload.Meta)
	if err != nil {
		return common.SendBadRequestResponse(c, "Error occurred while decoding email")
	}

	appTokenService := services.NewAppTokenService(h.DB)
	userService := services.NewUserService(h.DB)

	retrievedUser, err := userService.GetUserByEmail(string(email))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendNotFoundResponse(c, "Invalid password reset token")
		}

		return common.SendInternalServerErrorResponse(c, "Error occurred while fetching user")
	}

	appToken, err := appTokenService.ValidateResetPasswordToken(*retrievedUser, payload.Token)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	err = userService.ChangeUserPassword(payload.Password, *retrievedUser)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	appTokenService.InValidateToken(retrievedUser.ID, *appToken)
	return common.SendSuccessResponse(c, "Reset password successfully", nil)
}
