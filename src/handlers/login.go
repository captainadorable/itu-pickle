package handlers

import (
	"itu-pickle/utils"
	index "itu-pickle/views/index"

	"github.com/labstack/echo/v4"
)


func HandleLoginPost(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	data := utils.Login(username, password)	
	return utils.Render(c, index.LoginPanel(data.Fullname, data.LoggedIn))
}
