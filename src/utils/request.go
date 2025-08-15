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
    ], "SCRN": [ %s ]
  }`, ecrnString, scrnString))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		Log(fmt.Sprintf("Request hatası: %v", err))
		return
	}
	r.Header.Add("Authorization", token)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("User-Agent", config.Agent)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		Log(fmt.Sprintf("Request hatası: %v", err))
		return
	}
	defer res.Body.Close()

	// dump2, err := httputil.DumpResponse(res, true)
	// Log(string(dump2))

	Log(res.Status)
	if res.StatusCode == 401 {
		Log("Token hatalı")
		return
	}
	if res.StatusCode == 501 {
		Log("API hatası")
		return
	}
	if res.StatusCode != 200 {
		Log("Bir hata oluştu")
		return
	}

	post := &Post{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		Log("Decode hatası")
		return
	}

	// print response
	Log("<--SONUÇLAR-->")
	Log("Alınacak dersler")
	for _, ecrn := range post.EcrnResultList {
		// Write those in 1 line
		Log(fmt.Sprintf("[%s] -> [%t] -> [%d] -> [%s]", ecrn.Crn, ecrn.OperationFinished, ecrn.StatusCode, ecrn.ResultCode))
		Log(fmt.Sprintf(returnValues[ecrn.ResultCode], ecrn.Crn))
	}

	Log("Bırakılacak dersler")
	for _, scrn := range post.ScrnResultList {
		// Write those in 1 line
		Log(fmt.Sprintf("[%s] -> [%t] -> [%d] -> [%s]", scrn.Crn, scrn.OperationFinished, scrn.StatusCode, scrn.ResultCode))
		Log(fmt.Sprintf(returnValues[scrn.ResultCode], scrn.Crn))
	}
}

type KayitSinifResultList struct {
	Crn             string `json:"crn"`
	BransKodu       string `json:"bransKodu"`
	DersKodu        string `json:"dersKodu"`
	DersAdiTR       string `json:"dersAdiTR"`
	YerZamanBilgiEN string `json:"yerZamanBilgiEN"`
}

type ScheduleResponse struct {
	Results []KayitSinifResultList `json:"kayitSinifResultList"`
}

func ScheduleRequest(token string) ScheduleResponse {
	url := "https://obs.itu.edu.tr/api/ogrenci/sinif/KayitliSinifListesi/774" //775
	resp := &ScheduleResponse{}

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log(fmt.Sprintf("Request hatası: %v", err))
		return *resp
	}
	r.Header.Add("Authorization", token)
	r.Header.Add("Accepts", "application/json")
	r.Header.Add("User-Agent", config.Agent)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println("request error", err)
		return *resp
	}
	defer res.Body.Close()

	derr := json.NewDecoder(res.Body).Decode(resp)
	if derr != nil {
		fmt.Println("decode error")
		return *resp
	}

	return *resp
}
