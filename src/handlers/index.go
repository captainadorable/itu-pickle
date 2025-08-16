package handlers

import (
	"itu-pickle/models"
	"itu-pickle/picker"
	"itu-pickle/utils"
	index "itu-pickle/views/index"

	"github.com/labstack/echo/v4"
)

func HandleIndex(c echo.Context) error {
	return index.Index(picker.Started, utils.Logcu.Messages, models.UserData{LoggedIn: false}).Render(c.Request().Context(), c.Response())
}

