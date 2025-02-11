package main

import (
	"itu-pickle/config"
	"itu-pickle/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
  e := echo.New()
  config.Config()

  e.GET("/ws", config.Logcu.WebSocketHandler)

  e.GET("/", handlers.HandleIndex)
  
  e.POST("/start", handlers.HandleStart)
  e.GET("/stop", handlers.HandleStop)

  e.Logger.Fatal(e.Start(":3001"))
}

