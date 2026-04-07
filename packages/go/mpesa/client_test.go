package mpesa

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfigCreation(t *testing.T) {
	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Passkey:             "test_passkey",
		Shortcode:           "174379",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
		B2cShortcode:        "600000",
		Callbacks: map[string]string{
			"callback_url":         "https://test.com/callback",
			"b2c_result_url":       "https://test.com/b2c_result",
			"b2c_timeout_url":      "https://test.com/b2c_timeout",
			"b2b_result_url":       "https://test.com/b2b_result",
			"b2b_timeout_url":      "https://test.com/b2b_timeout",
			"b2pochi_result_url":   "https://test.com/b2pochi_result",
			"b2pochi_timeout_url":  "https://test.com/b2pochi_timeout",
			"balance_result_url":   "https://test.com/balance_result",
			"balance_timeout_url":  "https://test.com/balance_timeout",
			"status_result_url":    "https://test.com/status_result",
			"status_timeout_url":   "https://test.com/status_timeout",
			"reversal_result_url":  "https://test.com/reversal_result",
			"reversal_timeout_url": "https://test.com/reversal_timeout",
		},
	}

	if config.Environment != "sandbox" {
		t.Errorf("expected sandbox, got %s", config.Environment)
	}
	if config.MpesaConsumerKey != "test_key" {
		t.Errorf("expected test_key, got %s", config.MpesaConsumerKey)
	}
	if config.Shortcode != "174379" {
		t.Errorf("expected 174379, got %s", config.Shortcode)
	}
	if config.B2cShortcode != "600000" {
		t.Errorf("expected 600000, got %s", config.B2cShortcode)
	}
}

func TestMpesaClientCreation(t *testing.T) {
	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Passkey:             "test_passkey",
		Shortcode:           "174379",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
	}

	client := NewClient(config)
	if client == nil {
		t.Error("expected client, got nil")
	}
}

func TestHelpersPhoneValidator(t *testing.T) {
	helpers := NewHelpers(&Config{})

	tests := []struct {
		input    string
		expected string
	}{
		{"+254712345678", "254712345678"},
		{"0712345678", "254712345678"},
		{"712345678", "254712345678"},
		{"254712345678", "254712345678"},
	}

	for _, tt := range tests {
		result := helpers.PhoneValidator(tt.input)
		if result != tt.expected {
			t.Errorf("PhoneValidator(%s) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

func TestHelpersGetFormattedTimestamp(t *testing.T) {
	helpers := NewHelpers(&Config{})
	result := helpers.GetFormattedTimestamp()

	if len(result) != 14 {
		t.Errorf("expected 14 characters, got %d", len(result))
	}
}

func TestHelpersLipaNaMpesaPassword(t *testing.T) {
	config := &Config{
		Shortcode: "174379",
		Passkey:   "test_passkey",
	}
	helpers := NewHelpers(config)
	result := helpers.LipaNaMpesaPassword()

	if result == "" {
		t.Error("expected non-empty password")
	}
}

func TestHelpersGetConfig(t *testing.T) {
	config := &Config{
		Shortcode:     "174379",
		Passkey:       "test_passkey",
		TillNumber:    "123456",
		InitiatorName: "testapi",
	}
	helpers := NewHelpers(config)

	tests := []struct {
		key      string
		expected string
	}{
		{"shortcode", "174379"},
		{"passkey", "test_passkey"},
		{"till_number", "123456"},
		{"initiator_name", "testapi"},
		{"unknown_key", ""},
	}

	for _, tt := range tests {
		result := helpers.GetConfig(tt.key)
		if result != tt.expected {
			t.Errorf("GetConfig(%s) = %s; want %s", tt.key, result, tt.expected)
		}
	}
}

func TestHelpersResolveCallbackURL(t *testing.T) {
	config := &Config{
		Callbacks: map[string]string{
			"callback_url": "https://config.url/callback",
		},
	}
	helpers := NewHelpers(config)

	result := helpers.ResolveCallbackURL("https://param.url/callback", "callback_url")
	if result != "https://param.url/callback" {
		t.Errorf("expected param URL, got %s", result)
	}

	result = helpers.ResolveCallbackURL("", "callback_url")
	if result != "https://config.url/callback" {
		t.Errorf("expected config URL, got %s", result)
	}

	result = helpers.ResolveCallbackURL("", "unknown_url")
	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}
}

func TestNewClient(t *testing.T) {
	config := &Config{
		Environment: "sandbox",
	}
	client := NewClient(config)

	if client.baseURL != "https://sandbox.safaricom.co.ke" {
		t.Errorf("expected sandbox URL, got %s", client.baseURL)
	}

	config.Environment = "production"
	client = NewClient(config)
	if client.baseURL != "https://api.safaricom.co.ke" {
		t.Errorf("expected production URL, got %s", client.baseURL)
	}
}

func TestAuthToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/v1/generate" {
			t.Errorf("expected /oauth/v1/generate, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "test_token_123",
			"expires_in":   "3599",
		})
	}))
	defer server.Close()

	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
	}
	client := NewClient(config)
	client.baseURL = server.URL

	token, err := client.auth.GetAccessToken("C2B")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if token != "test_token_123" {
		t.Errorf("expected test_token_123, got %s", token)
	}
}

func TestStkPush(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ResponseCode":        "0",
			"ResponseDescription": "Success",
			"MerchantRequestID":   "29115-34620561-1",
			"CheckoutRequestID":   "ws_CO_191220191020363925",
		})
	}))
	defer server.Close()

	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
		Callbacks: map[string]string{
			"callback_url": "https://test.com/callback",
		},
	}
	client := NewClient(config)
	client.baseURL = server.URL

	result, err := client.Stkpush("254712345678", 100, "12345", "https://test.com/callback")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result["ResponseCode"] != "0" {
		t.Errorf("expected ResponseCode 0, got %v", result["ResponseCode"])
	}
}

func TestStkPushThrowsErrorWhenAccountReferenceMissing(t *testing.T) {
	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
	}
	client := NewClient(config)

	_, err := client.Stkpush("254712345678", 100, "", "https://test.com/callback")
	if err == nil {
		t.Error("expected error for missing account reference")
	}
}

func TestStkPushThrowsErrorWhenTillNumberRequired(t *testing.T) {
	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
	}
	client := NewClient(config)

	_, err := client.Stkpush("254712345678", 100, "12345", "https://test.com/callback", "CustomerBuyGoodsOnline")
	if err == nil {
		t.Error("expected error when till number required for TILL transaction")
	}
}

func TestB2C(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ResponseCode":             "0",
			"ResponseDescription":      "Success",
			"ConversationID":           "AG_20231217_201020363925",
			"OriginatorConversationID": "201020363925",
		})
	}))
	defer server.Close()

	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
		B2cShortcode:        "600000",
		Callbacks: map[string]string{
			"b2c_result_url":  "https://test.com/b2c_result",
			"b2c_timeout_url": "https://test.com/b2c_timeout",
		},
	}
	client := NewClient(config)
	client.baseURL = server.URL

	result, err := client.B2c("254712345678", "BusinessPayment", 100, "Test payment")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result["ResponseCode"] != "0" {
		t.Errorf("expected ResponseCode 0, got %v", result["ResponseCode"])
	}
}

func TestB2B(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ResponseCode":   "0",
			"ConversationID": "AG_20231217_201020363925",
		})
	}))
	defer server.Close()

	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
		B2cShortcode:        "600000",
		Callbacks: map[string]string{
			"b2b_result_url":  "https://test.com/b2b_result",
			"b2b_timeout_url": "https://test.com/b2b_timeout",
		},
	}
	client := NewClient(config)
	client.baseURL = server.URL

	result, err := client.B2b("600000", "BusinessPayBill", 100, "Test payment", "12345")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result["ResponseCode"] != "0" {
		t.Errorf("expected ResponseCode 0, got %v", result["ResponseCode"])
	}
}

func TestB2BThrowsErrorWhenAccountNumberMissing(t *testing.T) {
	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
		B2cShortcode:        "600000",
	}
	client := NewClient(config)

	_, err := client.B2b("600000", "BusinessPayBill", 100, "Test payment", "")
	if err == nil {
		t.Error("expected error for missing account number")
	}
}

func TestC2bregisterURLS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ResponseCode":        "0",
			"ResponseDescription": "success",
		})
	}))
	defer server.Close()

	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
		Callbacks: map[string]string{
			"c2b_confirmation_url": "https://test.com/c2b_confirm",
			"c2b_validation_url":   "https://test.com/c2b_validate",
		},
	}
	client := NewClient(config)
	client.baseURL = server.URL

	result, err := client.C2bregisterURLS("600000", "https://test.com/confirm", "https://test.com/validate")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result["ResponseCode"] != "0" {
		t.Errorf("expected ResponseCode 0, got %v", result["ResponseCode"])
	}
}

func TestC2bsimulate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ResponseCode":        "0",
			"ResponseDescription": "Success",
		})
	}))
	defer server.Close()

	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
	}
	client := NewClient(config)
	client.baseURL = server.URL

	result, err := client.C2bsimulate("254712345678", 100, "600000", "CustomerPayBillOnline", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result["ResponseCode"] != "0" {
		t.Errorf("expected ResponseCode 0, got %v", result["ResponseCode"])
	}
}

func TestAccountBalance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ResponseCode":        "0",
			"ResponseDescription": "Success",
		})
	}))
	defer server.Close()

	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
		Callbacks: map[string]string{
			"balance_result_url":  "https://test.com/balance_result",
			"balance_timeout_url": "https://test.com/balance_timeout",
		},
	}
	client := NewClient(config)
	client.baseURL = server.URL

	result, err := client.AccountBalance("600000", 4, "Check balance")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result["ResponseCode"] != "0" {
		t.Errorf("expected ResponseCode 0, got %v", result["ResponseCode"])
	}
}

func TestTransactionStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ResponseCode":        "0",
			"ResponseDescription": "Success",
		})
	}))
	defer server.Close()

	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
		Callbacks: map[string]string{
			"status_result_url":  "https://test.com/status_result",
			"status_timeout_url": "https://test.com/status_timeout",
		},
	}
	client := NewClient(config)
	client.baseURL = server.URL

	result, err := client.TransactionStatus("600000", "123456789", 1, "Check status")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result["ResponseCode"] != "0" {
		t.Errorf("expected ResponseCode 0, got %v", result["ResponseCode"])
	}
}

func TestReversal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ResponseCode":        "0",
			"ResponseDescription": "Success",
		})
	}))
	defer server.Close()

	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
		Callbacks: map[string]string{
			"reversal_result_url":  "https://test.com/reversal_result",
			"reversal_timeout_url": "https://test.com/reversal_timeout",
		},
	}
	client := NewClient(config)
	client.baseURL = server.URL

	result, err := client.Reversal("600000", "123456789", 100, "Reverse transaction")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result["ResponseCode"] != "0" {
		t.Errorf("expected ResponseCode 0, got %v", result["ResponseCode"])
	}
}

func TestB2Pochi(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ResponseCode":   "0",
			"ConversationID": "AG_20231217_201020363925",
		})
	}))
	defer server.Close()

	config := &Config{
		Environment:         "sandbox",
		MpesaConsumerKey:    "test_key",
		MpesaConsumerSecret: "test_secret",
		Shortcode:           "174379",
		Passkey:             "test_passkey",
		InitiatorName:       "testapi",
		InitiatorPassword:   "test_password",
		B2cShortcode:        "600000",
		Callbacks: map[string]string{
			"b2pochi_result_url":  "https://test.com/b2pochi_result",
			"b2pochi_timeout_url": "https://test.com/b2pochi_timeout",
		},
	}
	client := NewClient(config)
	client.baseURL = server.URL

	result, err := client.B2pochi("254712345678", 100, "Pochi payment")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result["ResponseCode"] != "0" {
		t.Errorf("expected ResponseCode 0, got %v", result["ResponseCode"])
	}
}
