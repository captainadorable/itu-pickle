package handlers

import (
	"itu-pickle/utils"
	views "itu-pickle/views/schedule"
	"strings"

	"github.com/labstack/echo/v4"
)

func HandleSchedule(c echo.Context) error {
  return utils.Render(c, views.Schedule(utils.ScheduleResponse{}))
}

func HandleGetSchedule(c echo.Context) error {
	// Parse form values
	token := c.FormValue("token")

  // check values
  token = strings.TrimSpace(token)
  if !strings.HasPrefix(token, "Bearer ") {
    token = "Bearer "+token
  }

  schedule := utils.ScheduleRequest(token)

  
	// Send response
	return utils.Render(c, views.Schedule(schedule))
}
