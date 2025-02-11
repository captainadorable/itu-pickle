package handlers

import (
	"itu-pickle/config"
	"itu-pickle/picker"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func HandleStart(c echo.Context) error {
	// Parse form values
	token := c.FormValue("token")
	ecrn := c.FormValue("ecrn")
	scrn := c.FormValue("scrn")
	timeout := c.FormValue("timeout")

  // check values
  if !strings.HasPrefix(token, "Bearer ") {
    config.Logcu.Log("[HATA]: Token hatalı girildi: Bearer <token>")
    return c.String(http.StatusBadRequest, "Invalid token.")
  }

  if !isValidPattern(ecrn) {
    config.Logcu.Log("[HATA]: Alınacak dersler hatalı girildi: CRN1,CRN2,CRN3")
    return c.String(http.StatusBadRequest, "Timeout is not a number.")
  }

  if !isValidPattern(scrn) {
    config.Logcu.Log("[HATA]: Bırakılacak dersler hatalı girildi: CRN1,CRN2,CRN3")
    return c.String(http.StatusBadRequest, "Timeout is not a number.")
  }

  timeout_parsed, err := strconv.Atoi(timeout)
  if err != nil {
    config.Logcu.Log("[HATA]: Süre bir sayı olmalı.")
    return c.String(http.StatusBadRequest, "Timeout is not a number.")
  }

  picker.Start(token, timeout_parsed, strings.Split(ecrn, ","), strings.Split(scrn, ","))

	// Send response
	return c.String(200, "Started")
}

func HandleStop(c echo.Context) error {
  picker.StopChan <- struct{}{}
  return c.String(200, "Stopped") 
}

func isValidPattern(s string) bool {
	// Regular expression to match "somevalue,somevalue,somevalue,..."
	re := regexp.MustCompile(`^(\w+,)*\w+$`)
	return re.MatchString(s)
}
