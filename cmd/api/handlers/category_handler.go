package handlers

import (
	"github.com/labstack/echo/v4"
	"pionira/cmd/api/services"
	"pionira/common"
)

func (h *Handler) ListCategories(c echo.Context) error {
	categoryService := services.NewCategoryService(h.DB)

	categories, err := categoryService.List()
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "ok", categories)

}
