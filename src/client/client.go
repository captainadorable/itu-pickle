package client

import (
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

// ApiClient manages the HTTP client and session state for the ITU website.
type ApiClient struct {
	Client *http.Client
}

// NewApiClient creates a single, reusable client with its own cookie jar.
func NewApiClient() *ApiClient {
	jar, _ := cookiejar.New(nil)
	return &ApiClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
			Jar:     jar,
		},
	}
}

// LoginWithCredentials attempts a full login and returns the user data.
func (api *ApiClient) LoginWithCredentials(username, password string) (models.UserData, error) {
	// Step 1: Get the login page to scrape tokens
	req, _ := http.NewRequest("GET", startURL, nil)
	req.Header.Set("User-Agent", config.Agent)
	initialResp, err := api.Client.Do(req)
	if err != nil {
		return models.UserData{}, fmt.Errorf("initial request failed: %w", err)
	}
	defer initialResp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(initialResp.Body)
	if err != nil {
		return models.UserData{}, fmt.Errorf("failed to parse login page: %w", err)
	}

	viewstate, _ := doc.Find("input[name=__VIEWSTATE]").Attr("value")
	viewstategenerator, _ := doc.Find("input[name=__VIEWSTATEGENERATOR]").Attr("value")
	eventvalidation, _ := doc.Find("input[name=__EVENTVALIDATION]").Attr("value")

	if viewstate == "" || eventvalidation == "" {
		return models.UserData{}, fmt.Errorf("could not find form tokens on login page")
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
		return models.UserData{}, fmt.Errorf("login post request failed: %w", err)
	}
	defer postResp.Body.Close()

	if !strings.Contains(postResp.Request.URL.String(), "https://obs.itu.edu.tr") {
		return models.UserData{}, fmt.Errorf("login failed, not redirected back to main site")
	}

	// Step 3: Get token and user info
	token, err := api.getToken()
	if err != nil {
		return models.UserData{}, fmt.Errorf("failed to get token after login: %w", err)
	}

	userInfo, err := api.getUserInfo(token)
	if err != nil {
		return models.UserData{}, fmt.Errorf("failed to get user info after login: %w", err)
	}
    
	// Step 5: Return the complete user data
	userData := models.UserData{
		Fullname: userInfo.KisiselBilgiler.AdSoyad,
		LoggedIn: true,
		Tokens:   []models.Token{token},
	}
	fmt.Println(postResp.Request.URL.String())
	return userData, nil
}

func (api *ApiClient) getToken() (models.Token, error) {
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

// getUserInfo is also a private helper method.
func (api *ApiClient) getUserInfo(token models.Token) (models.MainApiResponse, error) {
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

	// Read and print the final API response body.
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
