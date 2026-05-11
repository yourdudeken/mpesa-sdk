package mpesa

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HTTPClient struct {
	baseURL string
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{baseURL: baseURL}
}

func (c *HTTPClient) Post(url string, body map[string]interface{}, token string) (map[string]interface{}, error) {
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	return result, nil
}

func (c *HTTPClient) Get(url string, token string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	return result, nil
}

func (c *HTTPClient) GetBaseURL() string {
	return c.baseURL
}
