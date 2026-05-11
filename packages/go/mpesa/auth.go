package mpesa

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type tokenEntry struct {
	token  string
	expiry time.Time
}

type Auth struct {
	config     *Config
	baseURL    *string
	tokenCache map[string]tokenEntry
}

func NewAuth(config *Config, baseURL *string) *Auth {
	return &Auth{
		config:     config,
		baseURL:    baseURL,
		tokenCache: make(map[string]tokenEntry),
	}
}

func (a *Auth) GetAccessToken(shortCodeType string) (string, error) {
	if entry, ok := a.tokenCache[shortCodeType]; ok {
		if entry.token != "" && time.Now().Before(entry.expiry) {
			return entry.token, nil
		}
	}
	return a.generateAccessToken(shortCodeType)
}

func (a *Auth) generateAccessToken(shortCodeType string) (string, error) {
	consumerKey := a.config.MpesaConsumerKey
	consumerSecret := a.config.MpesaConsumerSecret

	if shortCodeType == "B2C" || shortCodeType == "B2B" {
		if a.config.B2cConsumerKey != "" && a.config.B2cConsumerSecret != "" {
			consumerKey = a.config.B2cConsumerKey
			consumerSecret = a.config.B2cConsumerSecret
		}
	}

	auth := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))

	url := *a.baseURL + "/oauth/v1/generate?grant_type=client_credentials"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	token := result["access_token"].(string)
	expiresIn := int(result["expires_in"].(float64))
	entry := tokenEntry{
		token:  token,
		expiry: time.Now().Add(time.Duration(expiresIn-60) * time.Second),
	}
	a.tokenCache[shortCodeType] = entry

	return token, nil
}

func (a *Auth) ClearToken() {
	a.tokenCache = make(map[string]tokenEntry)
}
