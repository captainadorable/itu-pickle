package config

import (
	"flag"
	"itu-pickle/logger"
	"strings"
)


var Logcu *logger.Logger
var Port string
var Url string

func Config() {
  flag.StringVar(&Port, "port", ":3001", "--port PORT")
  flag.StringVar(&Url, "url", "https://obs.itu.edu.tr/api/ders-kayit/v21", "--url URL")
  flag.Parse()

  if !strings.HasPrefix(Port, ":") {
    Port = ":" + Port 
  }

  Logcu = logger.NewLogger()
}
