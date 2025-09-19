package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"itu-pickle/config"
	"itu-pickle/models"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	startURL = "https://obs.itu.edu.tr/"
	ogrenciURL = "https://obs.itu.edu.tr/ogrenci"
	jwtURL   = "https://obs.itu.edu.tr/ogrenci/auth/jwt"
	infoURL  = "https://obs.itu.edu.tr/api/ogrenci/KisiselBilgiler/"
)

type ApiClient struct {
	Client *http.Client
	UserData models.UserData
}

func NewApiClient() *ApiClient {
	jar, _ := cookiejar.New(nil)
	userData := models.UserData{}

	return &ApiClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
			Jar:     jar,
		},
		UserData: userData,
	}
}

func (api *ApiClient) LoginWithCredentials(username, password string) (models.UserData, error) {
	// Step 1: Get the login page to scrape tokens
	req, _ := http.NewRequest("GET", startURL, nil)
	req.Header.Set("User-Agent", config.Agent)
	initialResp, err := api.Client.Do(req)
	if err != nil {
		return models.UserData{}, fmt.Errorf("İstek gönderilemedi: %w", err)
	}
	defer initialResp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(initialResp.Body)
	if err != nil {
		return models.UserData{}, fmt.Errorf("Giriş sayfası ayrıştırılamadı: %w", err)
	}

	viewstate, _ := doc.Find("input[name=__VIEWSTATE]").Attr("value")
	viewstategenerator, _ := doc.Find("input[name=__VIEWSTATEGENERATOR]").Attr("value")
	eventvalidation, _ := doc.Find("input[name=__EVENTVALIDATION]").Attr("value")

	if viewstate == "" || eventvalidation == "" {
		return models.UserData{}, fmt.Errorf("Form tokenleri bulunamadı.")
	}
	
	// Step 2: Post the credentials
	formData := url.Values{
		"__VIEWSTATE":                         {viewstate},
		"__VIEWSTATEGENERATOR":                {viewstategenerator},
		"__EVENTVALIDATION":                   {eventvalidation},
		"ctl00$ContentPlaceHolder1$tbUserName": {username},
		"ctl00$ContentPlaceHolder1$tbPassword": {password},
		"ctl00$ContentPlaceHolder1$btnLogin":   {"Giriş / Login"},
	}

	loginPageURL := initialResp.Request.URL.String()
	postResp, err := api.Client.Post(loginPageURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return models.UserData{}, fmt.Errorf("Giriş isteği başarısız: %w", err)
	}
	defer postResp.Body.Close()

	if !strings.Contains(postResp.Request.URL.String(), "https://obs.itu.edu.tr") {
		return models.UserData{}, fmt.Errorf("Giriş başarısız. Kullanıcı adı veya parola hatalı olabilir.")
	}

	// Step 3: Get token and user info
	token, err := api.GetToken()
	if err != nil {
		return models.UserData{}, fmt.Errorf("Giriş yapıldı ancak token alınamadı: %w", err)
	}

	userInfo, err := api.GetUserInfo(token)
	if err != nil {
		return models.UserData{}, fmt.Errorf("Giriş yapıldı ancak kullanıcı bilgisi alınamadı: %w", err)
	}
    
	// Step 5: Return the complete user data
	userData := models.UserData{
		Fullname: userInfo.KisiselBilgiler.AdSoyad,
		LoggedIn: true,
		Tokens:   []models.Token{token},
	}
	api.UserData = userData
	return userData, nil
}

func (api *ApiClient) GetToken() (models.Token, error) {
	newToken := models.Token{}

	jwtResp, err := api.Client.Get(jwtURL)

	jwtResp.Header.Set("User-Agent", config.Agent)
	if err != nil {
		return newToken, err
	}
	defer jwtResp.Body.Close()

	// Read the plain text response body to get the token.
	tokenBytes, err := io.ReadAll(jwtResp.Body)
	if err != nil {
		return newToken, err
	}
	bearerToken := string(tokenBytes)

	// Because of the redirects status code is 200. I know this is not good.
	if len(bearerToken) > 500 {
		return newToken, fmt.Errorf("Token alınamadı.")
	}

	newToken.Token = bearerToken
	newToken.CreatedAt = time.Now()

	return newToken, nil
}

func (api *ApiClient) GetUserInfo(token models.Token) (models.MainApiResponse, error) {
	var responseData models.MainApiResponse

	apiReq, _ := http.NewRequest("GET", infoURL, nil)
	apiReq.Header.Set("Authorization", "Bearer "+token.Token)
	apiReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	apiReq.Header.Set("Accept", "application/json, text/plain, */*")
	apiReq.Header.Set("User-Agent", config.Agent)

	apiResp, err := api.Client.Do(apiReq)
	if err != nil {
		return responseData, err
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode != http.StatusOK {
		return responseData, fmt.Errorf("Kullanıcı bilgisi alınamadı")
	}

	bodyBytes, err := io.ReadAll(apiResp.Body)
	if err != nil {
		return responseData, err
	}
	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		return responseData, err
	}

	return responseData, nil
}

type CrnResult struct {
	Crn               string `json:"crn"`
	OperationFinished bool   `json:"operationFinished"`
	StatusCode        int    `json:"statusCode"`
	ResultCode        string `json:"resultCode"`
}

type PostResp struct {
	EcrnResultList []CrnResult `json:"ecrnResultList"`
	ScrnResultList []CrnResult `json:"scrnResultList"`
}

func (api *ApiClient) Request(ecrnList, scrnList []string, token string) (PostResp, error) {
	url := config.Url
	postResp := PostResp{}

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
		return postResp, err
	}
	r.Header.Set("Authorization", token)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("User-Agent", config.Agent)

	res, err := api.Client.Do(r)
	if err != nil {
		return postResp, err
	}
	defer res.Body.Close()

	// dump2, err := httputil.DumpResponse(res, true)
	// fmt.Println(string(dump2))

	if res.StatusCode == 401 {
		err = fmt.Errorf("HTTP/2.0 401 Unauthorized")
		return postResp, err
	}
	if res.StatusCode == 501 {
		err = fmt.Errorf("HTTP/2.0 501 Server")
		return postResp, err
	}
	if res.StatusCode != 200 {
		return postResp, err
	}

	derr := json.NewDecoder(res.Body).Decode(&postResp)
	if derr != nil {
		return postResp, err
	}

	return postResp, nil
}
