package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"itu-pickle/config"
	"net/http"
	"strings"
)

type CrnResult struct {
	Crn               string `json:"crn"`
	OperationFinished bool   `json:"operationFinished"`
	StatusCode        int    `json:"statusCode"`
	ResultCode        string `json:"resultCode"`
}

type Post struct {
	EcrnResultList []CrnResult `json:"ecrnResultList"`
  ScrnResultList []CrnResult `json:"scrnResultList"`
}

func Request(ecrnList, scrnList []string, token string) {
	url := config.Url

	ecrnString := ""
	for _, crn := range ecrnList {
		ecrnString += fmt.Sprintf(`"%s",`, crn)
	}
	ecrnString = strings.TrimSuffix(ecrnString, ",")

	scrnString := ""
	for _, crn := range scrnList {
		scrnString += fmt.Sprintf(`"%s",`, crn)
	}
	scrnString = strings.TrimSuffix(scrnString, ",")

	body := []byte(fmt.Sprintf(`{
    "ECRN": [
      %s
    ],
    "SCRN": [
      %s
    ]
  }`, ecrnString, scrnString))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
    config.Logcu.Log(fmt.Sprintf("Request hatası: %v", err))
    return
	}
	r.Header.Add("Authorization", token)
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
    config.Logcu.Log(fmt.Sprintf("Request hatası: %v", err))
    return
	}
	defer res.Body.Close()

  // dump2, err := httputil.DumpResponse(res, true)
  // config.Logcu.Log(string(dump2))

	config.Logcu.Log(res.Status)
  if res.StatusCode == 401 {
    config.Logcu.Log("Token hatalı")
    return
  }
  if res.StatusCode == 501 {
    config.Logcu.Log("API hatası")
    return
  }
  if res.StatusCode != 200 {
    config.Logcu.Log("Bir hata oluştu")
    return
  }

	post := &Post{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
    config.Logcu.Log("Decode hatası")
    return
	}

	// print response
  config.Logcu.Log("<--SONUÇLAR-->")
  config.Logcu.Log("Alınacak dersler")
	for _, ecrn := range post.EcrnResultList {
		// Write those in 1 line
		config.Logcu.Log(fmt.Sprintf("[%s] -> [%t] -> [%d] -> [%s]", ecrn.Crn, ecrn.OperationFinished, ecrn.StatusCode, ecrn.ResultCode))
		config.Logcu.Log(fmt.Sprintf(returnValues[ecrn.ResultCode], ecrn.Crn))
	}

  config.Logcu.Log("Bırakılacak dersler")
	for _, scrn := range post.ScrnResultList {
		// Write those in 1 line
		config.Logcu.Log(fmt.Sprintf("[%s] -> [%t] -> [%d] -> [%s]", scrn.Crn, scrn.OperationFinished, scrn.StatusCode, scrn.ResultCode))
		config.Logcu.Log(fmt.Sprintf(returnValues[scrn.ResultCode], scrn.Crn))
	}
}
