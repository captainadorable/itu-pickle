package main

import (
	"fmt"
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

  e.GET("/favicon.ico", func(c echo.Context) error {
    return c.File("./favicon.ico")
  })

  fmt.Println("Sunucu başlatıldı. Arayüz: http://localhost"+config.Port)
  e.Logger.Fatal(e.Start(config.Port))
}

