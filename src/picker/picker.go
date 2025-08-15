package picker

import (
	"fmt"
	"itu-pickle/utils"
	views "itu-pickle/views/index"
	"strings"
	"time"
)

var Started bool = false
var StopChan chan struct{} = make(chan struct{})

func Start(token string, timeout int, ecrn, scrn []string, targetHour int, targetMin int, targetSec int) {


  if Started {
    utils.SendComponent(views.FormButtons(Started))
    return
  }

  Started = true
  utils.SendComponent(views.FormButtons(Started))

  utils.Log("CRN'ler kontrol ediliyor.")
  if !CheckCrns(ecrn, scrn) {
    Stop()
    return
  }

	// Configure target time
	now := time.Now()
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), targetHour, targetMin, targetSec, 0, now.Location())
	duration := targetTime.Sub(now)

	// If the target time has already passed for today, schedule it for tomorrow
	if now.After(targetTime) {
		utils.Log(fmt.Sprintf("Hedef saat çoktan geçmiş (%s)", targetTime.Format("15:04:05")))
		Stop()
		return
	}

	targetTimeTimer := time.NewTimer(duration)
	utils.Log(fmt.Sprintf("Hedef saate kadar bekleniyor (%s)", duration.String()))

	Loop:
	for {
		select {
			case <-targetTimeTimer.C:
				break Loop
			case <-StopChan:
        Stop()
				return
		}
	}

	utils.Log("Hedef saate geldi.")

	// Start the picker
  utils.Log(fmt.Sprintf("%d saniye bekleniyor", timeout))
  go func() {
    tick := time.Tick(time.Duration(timeout) * time.Second)

    Loop:
    for {
      select {
      case <-tick:
        utils.Request(ecrn, scrn, token)
        utils.Log(fmt.Sprintf("%d saniye bekleniyor", timeout))
      case <-StopChan:
        Stop()
        break Loop
      }
    }
  }()
}

func Stop() {
  Started = false
  utils.SendComponent(views.FormButtons(Started))
  utils.Log("Durduruldu")
}

func CheckCrns(ecrns, scrns []string) bool {
  fmt.Println(len(ecrns), len(scrns))
  crns := append(ecrns, scrns...)

	crnList := utils.FindCrns(crns)
  if crnList == nil {
    utils.Log("Bir şeyler yanlış gitti.")
    return false
  }

  for _, crn := range crnList {
    utils.Log(crn)
    if strings.HasPrefix(crn, "Not Found") {
      utils.Log(fmt.Sprintf("CRN bulunamadı: %s", crn))
      return false
    } 
  }
  return true
}
