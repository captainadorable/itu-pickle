package picker

import (
	"fmt"
	"itu-pickle/config"
	"itu-pickle/utils"
	views "itu-pickle/views/index"
	"strings"
	"time"
)

var Started bool = false
var StopChan chan struct{} = make(chan struct{})

func Start(token string, timeout int, ecrn, scrn []string) {
  if Started {
    config.Logcu.SendComponent(views.FormButtons(Started))
    return
  }

  Started = true
  config.Logcu.SendComponent(views.FormButtons(Started))

  config.Logcu.Log("CRN'ler kontrol ediliyor.")
  if !CheckCrns(ecrn, scrn) {
    Stop()
    return
  }

  config.Logcu.Log(fmt.Sprintf("%d saniye bekleniyor", timeout))
  go func() {
    tick := time.Tick(time.Duration(timeout) * time.Second)

    Loop:
    for {
      select {
      case <-tick:
        utils.Request(ecrn, scrn, token)
        config.Logcu.Log(fmt.Sprintf("%d saniye bekleniyor", timeout))
      case <-StopChan:
        Stop()
        break Loop
      }
    }
  }()
}

func Stop() {
  Started = false
  config.Logcu.SendComponent(views.FormButtons(Started))
  config.Logcu.Log("Durduruldu")
}

func CheckCrns(ecrns, scrns []string) bool {
  crns := append(ecrns, scrns...)

	crnList := utils.FindCrns(crns)
  if crnList == nil {
    config.Logcu.Log("Bir şeyler yanlış gitti.")
    return false
  }

  for _, crn := range crnList {
    config.Logcu.Log(crn)
    if strings.HasPrefix(crn, "Not Found") {
      config.Logcu.Log(fmt.Sprintf("CRN bulunamadı: %s", crn))
      return false
    } 
  }
  return true
}
