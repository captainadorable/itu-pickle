package models

import (
	"net/http/cookiejar"
	"time"
)

// A simple struct to hold the JWT response
type JwtResponse struct {
	Token string `json:"token"`
}

// MainApiResponse corresponds to the overall JSON structure
type MainApiResponse struct {
	KisiselBilgiler KisiselBilgiler `json:"kisiselBilgiler"`
	StatusCode      int             `json:"statusCode"`
	ResultCode      string          `json:"resultCode"`
	ResultMessage   string          `json:"resultMessage"`
}
type KisiselBilgiler struct {
	AdSoyad            string `json:"adSoyad"`
	// KimlikNo           string `json:"kimlikNo"`
	// Cinsiyet           string `json:"cinsiyet"`
	// FakulteTR          string `json:"fakulteTR"`
	// FakulteEN          string `json:"fakulteEN"`
	// BolumAdiTR         string `json:"bolumAdiTR"`
	// BolumAdiEN         string `json:"bolumAdiEN"`
	// CinsiyetKodu       string `json:"cinsiyetKodu"`
	// Telefon            string `json:"telefon"`
	// EPosta             string `json:"ePosta"`
	// ItuEPosta          string `json:"ituePosta"`
	// IkincilFakulteTR   string `json:"ikincilFakulteTR"`
	// IkincilFakulteEN   string `json:"ikincilFakulteEN"`
	// IkincilBolumAdiTR  string `json:"ikincilBolumAdiTR"`
	// IkincilBolumAdiEN  string `json:"ikincilBolumAdiEN"`
	// IkincilProgramiVar bool   `json:"ikincilProgramiVar"`
}

type Token struct {
	Token string
	CreatedAt time.Time
}

type UserData struct {
	Jar *cookiejar.Jar 
	Fullname string
	Tokens []Token
	LoggedIn bool
}

