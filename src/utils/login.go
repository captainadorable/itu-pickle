package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"itu-pickle/config"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	startURL = "https://obs.itu.edu.tr/"
	jwtURL   = "https://obs.itu.edu.tr/ogrenci/auth/jwt"
	infoURL = "https://obs.itu.edu.tr/api/ogrenci/KisiselBilgiler/"
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

var User UserData = UserData{}

func Login(username, password string) UserData  {

	// Create cookiejar
	jar, err := cookiejar.New(nil)
	if err != nil {
		Log(fmt.Sprintf("Bir hata oluştu: %v", err))
		return User
	}
	User.Jar = jar
	
	// Create httpclient
	client := &http.Client{
		Jar: jar,
	}

	// Initial request
	req, _ := http.NewRequest("GET", startURL, nil)
	req.Header.Set("User-Agent", config.Agent)
		// TODO: make an element in config to default this header
	initialResp, err := client.Do(req)
	if err != nil {
		Log(fmt.Sprintf("Bir hata oluştu: %v", err))
		return User
	}
	defer initialResp.Body.Close()

	loginPageURL := initialResp.Request.URL.String()

	// Parse the HTML of the login page to find the hidden form values.
	doc, err := goquery.NewDocumentFromReader(initialResp.Body)
	if err != nil {
		Log(fmt.Sprintf("Giriş yapılamadı: %v", err))
		return User
	}

	viewstate, _ := doc.Find("input[name=__VIEWSTATE]").Attr("value")
	viewstategenerator, _ := doc.Find("input[name=__VIEWSTATEGENERATOR]").Attr("value")
	eventvalidation, _ := doc.Find("input[name=__EVENTVALIDATION]").Attr("value")

	if viewstate == "" || eventvalidation == "" {
		Log(fmt.Sprintf("Giriş yapılamadı: %v", err))
		return User
	}

	// Build the form data for the POST request.
	formData := url.Values{
		"__EVENTTARGET":                       {""},
		"__EVENTARGUMENT":                     {""},
		"__VIEWSTATE":                         {viewstate},
		"__VIEWSTATEGENERATOR":                {viewstategenerator},
		"__EVENTVALIDATION":                   {eventvalidation},
		"ctl00$ContentPlaceHolder1$tbUserName": {username},
		"ctl00$ContentPlaceHolder1$tbPassword": {password},
		"ctl00$ContentPlaceHolder1$btnLogin":   {"Giriş / Login"},
	}

	// Send the POST request to log in.
	postResp, err := client.Post(loginPageURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	postResp.Header.Set("User-Agent", config.Agent)
	if err != nil {
		Log(fmt.Sprintf("Giriş yapılamadı: %v", err))
		return User
	}
	defer postResp.Body.Close()

	// Check if we were successfully redirected back to the main site.
	if !strings.Contains(postResp.Request.URL.String(), "obs.itu.edu.tr") {
		Log("Giriş yapılamadı: Not redirected back to the main site.")
		return User
	}


	firstToken := GetToken()
	if firstToken.Token == "" {
		Log("Giriş yapılamadı: Token oluşturulamadı")
		return User
	}
	User.Tokens = append(User.Tokens, firstToken)
	
	apiReq, _ := http.NewRequest("GET", infoURL, nil)
	apiReq.Header.Set("Authorization", "Bearer "+firstToken.Token)
	apiReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	apiReq.Header.Set("Accept", "application/json, text/plain, */*")
	apiReq.Header.Set("User-Agent", config.Agent)

	apiResp, err := client.Do(apiReq)
	if err != nil {
		Log("Giriş yapılamadı: Api hatası")
		return User
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode != http.StatusOK {
		Log("Giriş yapılamadı: Api hatası")
		return User
	}

	// Read and print the final API response body.
	bodyBytes, err := io.ReadAll(apiResp.Body)
	if err != nil {
		Log("Giriş yapılamadı: Api hatası")
		return User
	}
	var responseData MainApiResponse
	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		Log("Giriş yapılamadı: Api hatası")
		return User
	}

	User.Fullname = responseData.KisiselBilgiler.AdSoyad
	User.LoggedIn = true

	return User
}

func GetToken() Token {
	newToken := Token{}

	if User.Jar == nil {
		Log(fmt.Sprintf("Token alınamadı: Cookie hatası"))
		return newToken
	}

	// Create httpclient
	client := &http.Client{
		Jar: User.Jar,
	}

	jwtResp, err := client.Get(jwtURL)
	jwtResp.Header.Set("User-Agent", config.Agent)
	if err != nil {
		Log(fmt.Sprintf("Token alınamadı: %v", err))
		return newToken
	}
	defer jwtResp.Body.Close()

	if jwtResp.StatusCode != http.StatusOK {
		Log(fmt.Sprintf("Token alınamadı: %v", err))
		return newToken
	}

	// Read the plain text response body to get the token.
	tokenBytes, err := io.ReadAll(jwtResp.Body)
	if err != nil {
		Log(fmt.Sprintf("Token alınamadı: %v", err))
		return newToken
	}
	bearerToken := string(tokenBytes)

	if bearerToken == "" {
		Log(fmt.Sprintf("Token alınamadı: %v", err))
		return newToken
	}

	newToken.Token = bearerToken
	newToken.CreatedAt = time.Now()

	return newToken
}
