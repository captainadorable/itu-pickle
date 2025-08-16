package main

import (
	"fmt"
	"itu-pickle/client"
	"itu-pickle/config"
	"itu-pickle/handlers"
	"itu-pickle/utils"

	"github.com/labstack/echo/v4"
)

func main() {
  e := echo.New()
  config.Config()

	utils.StartLogger()

	apiClient := client.NewApiClient()


  e.GET("/ws", utils.WebSocketHandler)

  e.GET("/", handlers.HandleIndex)
  e.GET("/schedule", handlers.HandleSchedule)
  
  e.POST("/start", handlers.HandleStart)
  e.GET("/stop", handlers.HandleStop)
  e.POST("/getSchedule", handlers.HandleGetSchedule)

	e.POST("/login", func(c echo.Context) error {
		return handlers.HandleLoginPost(c, apiClient)
	})

  e.GET("/favicon.ico", func(c echo.Context) error {
    return c.File("./favicon.ico")
  })

  fmt.Println("Sunucu başlatıldı. Arayüz: http://localhost"+config.Port)
  e.Logger.Fatal(e.Start(config.Port))
}

