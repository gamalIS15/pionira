package handlers

import (
	"github.com/labstack/echo/v4"
	"pionira/cmd/api/requests"
	"pionira/cmd/api/services"
	"pionira/common"
	"pionira/internal/models"
)

func (h *Handler) ChangeUserPassword(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendBadRequestResponse(c, "User not logged in")
	}

	payload := new(requests.ChangePasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	//validation
	validationErrors := h.ValidateBody(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	// compare client data password to hashed password
	if common.ComparePasswordHash(payload.CurrentPassword, user.Password) == false {
		return common.SendBadRequestResponse(c, "The supplied password does not match")
	}

	userService := services.NewUserService(h.DB)
	err := userService.ChangeUserPassword(payload.Password, user)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Password Change Successfully", nil)
}
