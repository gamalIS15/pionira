package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"pionira/common"
	"pionira/internal/models"
	"strings"
)

type AppMiddleware struct {
	Logger echo.Logger
	DB     *gorm.DB
}

func (am *AppMiddleware) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		c.Response().Header().Add("Vary", "Authorization")
		authHeader := c.Request().Header.Get("Authorization")

		if strings.HasPrefix(authHeader, "Bearer ") == false {
			return common.SendUnauthorizedResponse(c, "Please provide Bearer token")
		}

		authHeaderSplit := strings.Split(authHeader, " ")

		accessToken := authHeaderSplit[1]
		claims, err := common.ParseJWTToken(accessToken)
		if err != nil {
			return common.SendUnauthorizedResponse(c, err.Error())
		}

		if common.IsClaimExpired(claims) {
			return common.SendUnauthorizedResponse(c, "Token is expired")
		}

		var user models.UserModel
		result := am.DB.First(&user, claims.ID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return common.SendUnauthorizedResponse(c, "Invalid access token")
		}

		if result.Error != nil {
			return common.SendUnauthorizedResponse(c, "Invalid access token")
		}

		c.Set("user", user)
		return next(c)
	}
}
