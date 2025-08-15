package handlers

import (
	"fmt"
	"itu-pickle/utils"
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

	targethour := c.FormValue("targethour")
	targetminute := c.FormValue("targetminute")
	targetsecond := c.FormValue("targetsecond")

	fmt.Println("QWE", targethour, targetminute, targetsecond)

  // check values
  token = strings.TrimSpace(token)
  if !strings.HasPrefix(token, "Bearer ") {
    token = "Bearer "+token
  }

  ecrn = strings.TrimSpace(ecrn)
  if !isValidPattern(ecrn) {
    utils.Log("[HATA]: Alınacak dersler hatalı girildi: CRN1 CRN2 CRN3")
    return c.String(http.StatusBadRequest, "Timeout is not a number.")
  }

  scrn = strings.TrimSpace(scrn)
  if !isValidPattern(scrn) && len(scrn) != 0 {
    utils.Log("[HATA]: Bırakılacak dersler hatalı girildi: CRN1 CRN2 CRN3")
    return c.String(http.StatusBadRequest, "Timeout is not a number.")
  }

  targethour_parsed, err := strconv.Atoi(targethour)
  if err != nil {
    utils.Log("[HATA]: Hedef saat bir sayı olmalı.")
    return c.String(http.StatusBadRequest, "Target time is not a number.")
  }
  targetminute_parsed, err := strconv.Atoi(targetminute)
  if err != nil {
    utils.Log("[HATA]: Hedef dakika bir sayı olmalı.")
    return c.String(http.StatusBadRequest, "Target minute is not a number.")
  }
  targetsecond_parsed, err := strconv.Atoi(targetsecond)
  if err != nil {
    utils.Log("[HATA]: Hedef saniye bir sayı olmalı.")
    return c.String(http.StatusBadRequest, "Target second is not a number.")
  }


  timeout_parsed, err := strconv.Atoi(timeout)
  if err != nil {
    utils.Log("[HATA]: Süre bir sayı olmalı.")
    return c.String(http.StatusBadRequest, "Timeout is not a number.")
  }

  ecrnList := strings.Split(ecrn, " ")

  scrnList := strings.Split(scrn, " ")
  if len(scrn) == 0 {
    scrnList = []string{}
  }
    
  picker.Start(token, timeout_parsed, ecrnList, scrnList, targethour_parsed, targetminute_parsed, targetsecond_parsed)

	// Send response
	return c.String(200, "Started")
}

func HandleStop(c echo.Context) error {
  picker.StopChan <- struct{}{}
  return c.String(200, "Stopped") 
}

func isValidPattern(s string) bool {
	// Regular expression to match "somevalue somevalue somevalue,..."
	re := regexp.MustCompile(`^(\w+ )*\w+$`)
	return re.MatchString(s)
}
