package handlers

import (
	"itu-pickle/client" // Import our new client package
	"itu-pickle/models"
	"itu-pickle/utils"
	index "itu-pickle/views/index"

	"github.com/labstack/echo/v4"
)

func HandleLoginPost(c echo.Context, apiClient *client.ApiClient) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	userData, err := apiClient.LoginWithCredentials(username, password)
	if err != nil {
		utils.Log(err.Error())
		return index.LoginPanel(models.UserData{LoggedIn: false}).Render(c.Request().Context(), c.Response())
	}

	return index.LoginPanel(userData).Render(c.Request().Context(), c.Response())
}
