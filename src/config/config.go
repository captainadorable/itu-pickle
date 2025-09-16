package config

import (
	"flag"
	"strings"
)


var Port string
var Url string

var Agent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"

func Config() {
  flag.StringVar(&Port, "port", ":5454", "--port PORT")
  flag.StringVar(&Url, "url", "https://obs.itu.edu.tr/api/ders-kayit/v21", "--url URL")
  flag.Parse()

  if !strings.HasPrefix(Port, ":") {
    Port = ":" + Port 
  }
}
