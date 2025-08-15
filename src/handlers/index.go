package handlers

import (
	"itu-pickle/picker"
	"itu-pickle/utils"
	index "itu-pickle/views/index"

	"github.com/labstack/echo/v4"
)


func HandleIndex(c echo.Context) error {
  return utils.Render(c, index.Index(picker.Started, utils.Logcu.Messages, utils.User.Fullname, utils.User.LoggedIn))
}

