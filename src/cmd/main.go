package main

import (
	"fmt"
	"itu-pickle/config"
	"itu-pickle/handlers"
	"itu-pickle/utils"

	"github.com/labstack/echo/v4"
)

func main() {
  e := echo.New()
  config.Config()

	utils.StartLogger()

  e.GET("/ws", utils.WebSocketHandler)

  e.GET("/", handlers.HandleIndex)
  e.GET("/schedule", handlers.HandleSchedule)
  
  e.POST("/start", handlers.HandleStart)
  e.GET("/stop", handlers.HandleStop)
  e.POST("/getSchedule", handlers.HandleGetSchedule)

	e.POST("/login", handlers.HandleLoginPost)

  e.GET("/favicon.ico", func(c echo.Context) error {
    return c.File("./favicon.ico")
  })

  fmt.Println("Sunucu başlatıldı. Arayüz: http://localhost"+config.Port)
  e.Logger.Fatal(e.Start(config.Port))
}

