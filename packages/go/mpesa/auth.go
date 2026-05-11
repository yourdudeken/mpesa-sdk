package mpesa

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Auth struct {
	config      *Config
	baseURL     *string
	token       string
	tokenExpiry time.Time
}

func NewAuth(config *Config, baseURL *string) *Auth {
	return &Auth{
		config:  config,
		baseURL: baseURL,
	}
}

func (a *Auth) GetAccessToken(shortCodeType string) (string, error) {
	if a.token != "" && time.Now().Before(a.tokenExpiry) {
		return a.token, nil
	}
	return a.generateAccessToken(shortCodeType)
}

func (a *Auth) generateAccessToken(shortCodeType string) (string, error) {
	consumerKey := a.config.MpesaConsumerKey
	consumerSecret := a.config.MpesaConsumerSecret

	if shortCodeType == "B2C" || shortCodeType == "B2B" {
		consumerKey = a.config.B2cConsumerKey
		consumerSecret = a.config.B2cConsumerSecret
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

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	a.token = result["access_token"].(string)
	expiresIn := int(result["expires_in"].(float64))
	a.tokenExpiry = time.Now().Add(time.Duration(expiresIn-60) * time.Second)

	return a.token, nil
}

func (a *Auth) ClearToken() {
	a.token = ""
}
