package config

import "itu-pickle/logger"


var Logcu *logger.Logger


func Config() {
  Logcu = logger.NewLogger()
}
