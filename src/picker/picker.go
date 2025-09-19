package picker

import (
	"fmt"
	"itu-pickle/client"
	"itu-pickle/utils"
	views "itu-pickle/views/index"
	"time"
)

var Started bool = false
var StopChan chan struct{} = make(chan struct{})

type PickerConfig struct {
	Token string
	Token2 string
	Token3 string
	Crncontrol string
	Ecrn []string
	Scrn []string

	Timeout int
	TargetHour int
	TargetMin int
	TargetSec int

	ApiClient *client.ApiClient
}

func Start(conf PickerConfig) {
	// Form buttons
  if Started {
    utils.SendComponent(views.FormButtons(Started))
    return
  }

  Started = true
  utils.SendComponent(views.FormButtons(Started))

	// Crn control
	if conf.Crncontrol == "on" {
		CrnControl(conf.Ecrn, conf.Scrn)
	} else {
		utils.Log("CRN kontrolü atlanıyor.", "default")
	}

	// Configure target time
	now := time.Now()
	targetTime, tokenTime := ConfigureTime(conf.TargetHour, conf.TargetMin, conf.TargetSec )
	targetTimeDuration := targetTime.Sub(now)
	tokenTimeDuration := tokenTime.Sub(now)

	// If the target time has already passed for today, stop
	if now.After(targetTime) {
		utils.Log(fmt.Sprintf("Hedef saat çoktan geçmiş (%s)", targetTime.Format("15:04:05")), "error")
		Stop()
		return
	}

	// If token is empty, then start token timer
	if conf.Token == "" {
		tokenTimeTimer := time.NewTimer(tokenTimeDuration)
		utils.Log(fmt.Sprintf("Token almak için bekleniyor (son 30 saniyede alınır) (%s)", tokenTimeDuration.String()), "default")

		Loop1:
		for {
			select {
				case <-tokenTimeTimer.C:
					break Loop1
				case <-StopChan:
					Stop()
					return
			}
		}
		
		newToken, err := conf.ApiClient.GetToken()
		if err != nil {
			utils.Log(fmt.Sprintf("Token alınamadı: %v", err), "error")
			Stop()
			return
		}
		newToken.Token = "Bearer "+newToken.Token
		conf.Token = newToken.Token

		utils.Log("Token alındı", "success")
	}

	// Target time
	targetTimeTimer := time.NewTimer(targetTime.Sub(now))
	utils.Log(fmt.Sprintf("Hedef saate kadar bekleniyor (%s)", targetTimeDuration.String()), "default")

	Loop2:
	for {
		select {
			case <-targetTimeTimer.C:
				break Loop2
			case <-StopChan:
        Stop()
				return
		}
	}

	// Start the picker
	utils.Log("Hedef saat geldi. İstek gönderiliyor", "default")
	response, err := conf.ApiClient.Request(conf.Ecrn, conf.Scrn, conf.Token)
	if err != nil {
		utils.Log(fmt.Sprintf("Bir hata oluştu: %v", err), "error")
	} else {
		LogResponse(response)
	}

  utils.Log(fmt.Sprintf("%d saniye bekleniyor", conf.Timeout), "default")
  go func() {
    tick := time.Tick(time.Duration(conf.Timeout) * time.Second)

    Loop:
    for {
      select {
      case <-tick:
				response, err := conf.ApiClient.Request(conf.Ecrn, conf.Scrn, conf.Token)
				if err != nil {
					utils.Log(fmt.Sprintf("Bir hata oluştu: %v", err), "error")
				}
				LogResponse(response)

        utils.Log(fmt.Sprintf("%d saniye bekleniyor", conf.Timeout), "default")
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
  utils.Log("Durduruldu", "default")
}

func CrnControl(ecrn, scrn []string) {
	utils.Log("CRN kontrolü yapılıyor.", "default")

	crns := append(ecrn, scrn...)

	crnList, err := utils.FindCrns(crns)
	if err != nil {
		utils.Log(fmt.Sprintf("Bir şeyler yanlış gitti: %v", err), "error")
		Stop()
	}

	for _, crn := range crnList {
		if crn.Exists {
			utils.Log(fmt.Sprintf("%s: %s", crn.Crn, crn.Code), "success")
		} else {
			utils.Log(fmt.Sprintf("CRN bulunamadı: %s", crn.Crn), "error")
			Stop()
		} 
	}
}

func ConfigureTime(targetHour, targetMinute, targetSecond int) (time.Time, time.Time) {
	now := time.Now()
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), targetHour, targetMinute, targetSecond, 0, now.Location())
	tokenTime := targetTime.Add(time.Second * -30)
	if targetHour == 999 || targetMinute == 999 || targetSecond == 999 {
		targetTime = now
		tokenTime = now
	}

	return targetTime, tokenTime
}

func LogResponse(post client.PostResp) {
	// print response
	utils.Log("-- İşlem başarılı --", "success")
	utils.Log("Alınacak dersler", "success")
	for _, ecrn := range post.EcrnResultList {
		color := "error"
		if ecrn.ResultCode == "successResult" {
			color = "success"
		}
		// Log(fmt.Sprintf("[%s] -> [%t] -> [%d] -> [%s]", ecrn.Crn, ecrn.OperationFinished, ecrn.StatusCode, ecrn.ResultCode), color)
		utils.Log(fmt.Sprintf(utils.ReturnValues[ecrn.ResultCode], ecrn.Crn), color)
	}

	if len(post.ScrnResultList) != 0 {
		utils.Log("Bırakılacak dersler", "success")
		for _, scrn := range post.ScrnResultList {
			color := "error"
			if scrn.ResultCode == "successResult" {
				color = "success"
			}
			// Log(fmt.Sprintf("[%s] -> [%t] -> [%d] -> [%s]", scrn.Crn, scrn.OperationFinished, scrn.StatusCode, scrn.ResultCode), "success")
			utils.Log(fmt.Sprintf(utils.ReturnValues[scrn.ResultCode], scrn.Crn), color)
		}
	}
}
