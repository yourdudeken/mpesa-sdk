package mpesa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	Paybill = "CustomerPayBillOnline"
	Till    = "CustomerBuyGoodsOnline"
)

type Config struct {
	Environment         string
	MpesaConsumerKey    string
	MpesaConsumerSecret string
	B2cConsumerKey      string
	B2cConsumerSecret   string
	Passkey             string
	Shortcode           string
	TillNumber          string
	InitiatorName       string
	InitiatorPassword   string
	B2cShortcode        string
	Callbacks           map[string]string
}

type Mpesa struct {
	config      Config
	baseURL     string
	accessToken string
	tokenExpiry time.Time
}

func New(config Config) *Mpesa {
	baseURL := "https://sandbox.safaricom.co.ke"
	if config.Environment == "production" {
		baseURL = "https://api.safaricom.co.ke"
	}
	return &Mpesa{
		config:  config,
		baseURL: baseURL,
	}
}

func (m *Mpesa) getConfig(key string) string {
	switch key {
	case "environment":
		return m.config.Environment
	case "mpesa_consumer_key":
		return m.config.MpesaConsumerKey
	case "mpesa_consumer_secret":
		return m.config.MpesaConsumerSecret
	case "b2c_consumer_key":
		return m.config.B2cConsumerKey
	case "b2c_consumer_secret":
		return m.config.B2cConsumerSecret
	case "passkey":
		return m.config.Passkey
	case "shortcode":
		return m.config.Shortcode
	case "till_number":
		return m.config.TillNumber
	case "initiator_name":
		return m.config.InitiatorName
	case "initiator_password":
		return m.config.InitiatorPassword
	case "b2c_shortcode":
		return m.config.B2cShortcode
	default:
		return ""
	}
}

func (m *Mpesa) resolveCallbackUrl(paramURL, configKey string) (string, error) {
	var configURL string
	if m.config.Callbacks != nil {
		configURL = m.config.Callbacks[configKey]
	}
	if paramURL != "" {
		return paramURL, nil
	}
	if configURL != "" {
		return configURL, nil
	}
	return "", fmt.Errorf("ensure you have set the %s in the config or passed as a parameter", configKey)
}

func (m *Mpesa) phoneValidator(phoneNo string) string {
	phoneNo = strings.TrimPrefix(phoneNo, "+")
	if strings.HasPrefix(phoneNo, "0") {
		phoneNo = "254" + phoneNo[1:]
	} else if strings.HasPrefix(phoneNo, "7") {
		phoneNo = "254" + phoneNo
	}
	return phoneNo
}

func (m *Mpesa) getFormattedTimestamp() string {
	return time.Now().Format("20060102150405")
}

func (m *Mpesa) lipaNaMpesaPassword() string {
	timestamp := m.getFormattedTimestamp()
	password := m.getConfig("shortcode") + m.getConfig("passkey") + timestamp
	return base64.StdEncoding.EncodeToString([]byte(password))
}

func (m *Mpesa) generateAccessToken(shortCodeType string) (string, error) {
	if m.accessToken != "" && time.Now().Before(m.tokenExpiry) {
		return m.accessToken, nil
	}

	consumerKey := m.getConfig("mpesa_consumer_key")
	consumerSecret := m.getConfig("mpesa_consumer_secret")
	if shortCodeType == "B2C" || shortCodeType == "B2B" {
		consumerKey = m.getConfig("b2c_consumer_key")
		consumerSecret = m.getConfig("b2c_consumer_secret")
	}

	auth := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))

	url := m.baseURL + "/oauth/v1/generate?grant_type=client_credentials"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	accessToken := result["access_token"].(string)
	expiresIn := int(result["expires_in"].(float64))
	m.accessToken = accessToken
	m.tokenExpiry = time.Now().Add(time.Duration(expiresIn-60) * time.Second)

	return accessToken, nil
}

func (m *Mpesa) generateSecurityCredential() (string, error) {
	certPath := "certificates/SandboxCertificate.cer"
	if m.getConfig("environment") == "production" {
		certPath = "certificates/ProductionCertificate.cer"
	}

	certPath = filepath.Join(filepath.Dir(os.Args[0]), certPath)
	if _, err := os.Stat(certPath); err != nil {
		certPath = certPath
	}

	pubKeyData, err := os.ReadFile(certPath)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(pubKeyData)
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}

	rsaPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("not an RSA public key")
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, []byte(m.getConfig("initiator_password")))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (m *Mpesa) mpesaRequest(url string, body map[string]interface{}, shortCodeType string) (map[string]interface{}, error) {
	token, err := m.generateAccessToken(shortCodeType)
	if err != nil {
		return nil, err
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	return result, nil
}

func (m *Mpesa) Stkpush(params map[string]interface{}) (map[string]interface{}, error) {
	phonenumber := params["phonenumber"].(string)
	amount := int(params["amount"].(float64))
	accountNumber := params["accountNumber"].(string)
	callbackURL, _ := params["callbackUrl"].(string)
	transactionType := Paybill
	if t, ok := params["transactionType"].(string); ok {
		transactionType = t
	}
	shortCodeType := "C2B"
	if s, ok := params["shortCodeType"].(string); ok {
		shortCodeType = s
	}

	if accountNumber == "" {
		return nil, fmt.Errorf("an Account Reference is required for All transactions")
	}

	if transactionType == Till && m.getConfig("till_number") == "" {
		return nil, fmt.Errorf("till number is required for Buy Goods transactions")
	}

	url := m.baseURL + "/mpesa/stkpush/v1/processrequest"
	data := map[string]interface{}{
		"BusinessShortCode": m.getConfig("shortcode"),
		"Password":          m.lipaNaMpesaPassword(),
		"Timestamp":         m.getFormattedTimestamp(),
		"Amount":            amount,
		"PartyA":            m.phoneValidator(phonenumber),
		"PartyB":            map[bool]string{true: m.getConfig("shortcode"), false: m.getConfig("till_number")}[transactionType == Paybill],
		"TransactionType":   transactionType,
		"PhoneNumber":       m.phoneValidator(phonenumber),
		"TransactionDesc":   "Payment",
		"AccountReference":  accountNumber,
		"CallBackURL":       m.resolveCallbackUrl(callbackURL, "callback_url"),
	}

	return m.mpesaRequest(url, data, shortCodeType)
}

func (m *Mpesa) Stkquery(checkoutRequestID string, shortCodeType string) (map[string]interface{}, error) {
	url := m.baseURL + "/mpesa/stkpushquery/v1/query"
	data := map[string]interface{}{
		"BusinessShortCode": m.getConfig("shortcode"),
		"Password":          m.lipaNaMpesaPassword(),
		"Timestamp":         m.getFormattedTimestamp(),
		"CheckoutRequestID": checkoutRequestID,
	}

	if shortCodeType == "" {
		shortCodeType = "C2B"
	}
	return m.mpesaRequest(url, data, shortCodeType)
}

func (m *Mpesa) B2c(params map[string]interface{}) (map[string]interface{}, error) {
	phonenumber := params["phonenumber"].(string)
	commandID := params["commandId"].(string)
	amount := int(params["amount"].(float64))
	remarks := params["remarks"].(string)
	resultURL, _ := params["resultUrl"].(string)
	timeoutURL, _ := params["timeoutUrl"].(string)
	shortCodeType := "B2C"
	if s, ok := params["shortCodeType"].(string); ok {
		shortCodeType = s
	}

	url := m.baseURL + "/mpesa/b2c/v1/paymentrequest"
	body := map[string]interface{}{
		"InitiatorName":      m.getConfig("initiator_name"),
		"SecurityCredential": m.generateSecurityCredential(),
		"CommandID":          commandID,
		"Amount":             amount,
		"PartyA":             m.getConfig("b2c_shortcode"),
		"PartyB":             m.phoneValidator(phonenumber),
		"Remarks":            remarks,
		"Occassion":          "",
		"ResultURL":          m.resolveCallbackUrl(resultURL, "b2c_result_url"),
		"QueueTimeOutURL":    m.resolveCallbackUrl(timeoutURL, "b2c_timeout_url"),
	}

	return m.mpesaRequest(url, body, shortCodeType)
}

func (m *Mpesa) Validated_b2c(params map[string]interface{}) (map[string]interface{}, error) {
	phonenumber := params["phonenumber"].(string)
	commandID := params["commandId"].(string)
	amount := int(params["amount"].(float64))
	remarks := params["remarks"].(string)
	idNumber := params["idNumber"].(string)
	resultURL, _ := params["resultUrl"].(string)
	timeoutURL, _ := params["timeoutUrl"].(string)
	shortCodeType := "B2C"
	if s, ok := params["shortCodeType"].(string); ok {
		shortCodeType = s
	}

	url := m.baseURL + "/mpesa/b2cvalidate/v2/paymentrequest"
	body := map[string]interface{}{
		"InitiatorName":            m.getConfig("initiator_name"),
		"SecurityCredential":       m.generateSecurityCredential(),
		"CommandID":                commandID,
		"Amount":                   amount,
		"PartyA":                   m.getConfig("b2c_shortcode"),
		"PartyB":                   m.phoneValidator(phonenumber),
		"Remarks":                  remarks,
		"Occassion":                "",
		"OriginatorConversationID": m.getFormattedTimestamp(),
		"IDType":                   "01",
		"IDNumber":                 idNumber,
		"ResultURL":                m.resolveCallbackUrl(resultURL, "b2c_result_url"),
		"QueueTimeOutURL":          m.resolveCallbackUrl(timeoutURL, "b2c_timeout_url"),
	}

	return m.mpesaRequest(url, body, shortCodeType)
}

func (m *Mpesa) B2b(params map[string]interface{}) (map[string]interface{}, error) {
	receiverShortcode := params["receiverShortcode"].(string)
	commandID := params["commandId"].(string)
	amount := int(params["amount"].(float64))
	remarks := params["remarks"].(string)
	accountNumber, _ := params["accountNumber"].(string)
	resultURL, _ := params["resultUrl"].(string)
	timeoutURL, _ := params["timeoutUrl"].(string)
	shortCodeType := "B2B"
	if s, ok := params["shortCodeType"].(string); ok {
		shortCodeType = s
	}

	if commandID == "BusinessPayBill" && accountNumber == "" {
		return nil, fmt.Errorf("Account Number is required for BusinessPayBill CommandID")
	}

	url := m.baseURL + "/mpesa/b2b/v1/paymentrequest"
	body := map[string]interface{}{
		"Initiator":              m.getConfig("initiator_name"),
		"SecurityCredential":     m.generateSecurityCredential(),
		"CommandID":              commandID,
		"SenderIdentifierType":   "4",
		"RecieverIdentifierType": "4",
		"Amount":                 amount,
		"PartyA":                 m.getConfig("b2c_shortcode"),
		"PartyB":                 receiverShortcode,
		"AccountReference":       accountNumber,
		"Remarks":                remarks,
		"ResultURL":              m.resolveCallbackUrl(resultURL, "b2b_result_url"),
		"QueueTimeOutURL":        m.resolveCallbackUrl(timeoutURL, "b2b_timeout_url"),
	}

	return m.mpesaRequest(url, body, shortCodeType)
}

func (m *Mpesa) C2bregisterURLS(params map[string]interface{}) (map[string]interface{}, error) {
	shortcode := params["shortcode"].(string)
	confirmURL, _ := params["confirmUrl"].(string)
	validateURL, _ := params["validateUrl"].(string)
	shortCodeType := "C2B"
	if s, ok := params["shortCodeType"].(string); ok {
		shortCodeType = s
	}

	url := m.baseURL + "/mpesa/c2b/v2/registerurl"
	body := map[string]interface{}{
		"ShortCode":       shortcode,
		"ResponseType":    "Completed",
		"ConfirmationURL": m.resolveCallbackUrl(confirmURL, "c2b_confirmation_url"),
		"ValidationURL":   m.resolveCallbackUrl(validateURL, "c2b_validation_url"),
	}

	return m.mpesaRequest(url, body, shortCodeType)
}

func (m *Mpesa) C2bsimulate(params map[string]interface{}) (map[string]interface{}, error) {
	phonenumber := params["phonenumber"].(string)
	amount := int(params["amount"].(float64))
	shortcode := params["shortcode"].(string)
	commandID := params["commandId"].(string)
	accountNumber, _ := params["accountNumber"].(string)
	shortCodeType := "C2B"
	if s, ok := params["shortCodeType"].(string); ok {
		shortCodeType = s
	}

	url := m.baseURL + "/mpesa/c2b/v2/simulate"
	data := map[string]interface{}{
		"Msisdn":    m.phoneValidator(phonenumber),
		"Amount":    amount,
		"CommandID": commandID,
		"ShortCode": shortcode,
	}

	if commandID == Paybill {
		data["BillRefNumber"] = accountNumber
	}

	return m.mpesaRequest(url, data, shortCodeType)
}

func (m *Mpesa) TransactionStatus(params map[string]interface{}) (map[string]interface{}, error) {
	shortcode := params["shortcode"].(string)
	transactionID := params["transactionId"].(string)
	identifierType := int(params["identifierType"].(float64))
	remarks := params["remarks"].(string)
	resultURL, _ := params["resultUrl"].(string)
	timeoutURL, _ := params["timeoutUrl"].(string)
	shortCodeType := "C2B"
	if s, ok := params["shortCodeType"].(string); ok {
		shortCodeType = s
	}

	url := m.baseURL + "/mpesa/transactionstatus/v1/query"
	body := map[string]interface{}{
		"Initiator":          m.getConfig("initiator_name"),
		"SecurityCredential": m.generateSecurityCredential(),
		"CommandID":          "TransactionStatusQuery",
		"TransactionID":      transactionID,
		"PartyA":             shortcode,
		"IdentifierType":     identifierType,
		"Remarks":            remarks,
		"Occassion":          "",
		"ResultURL":          m.resolveCallbackUrl(resultURL, "status_result_url"),
		"QueueTimeOutURL":    m.resolveCallbackUrl(timeoutURL, "status_timeout_url"),
	}

	return m.mpesaRequest(url, body, shortCodeType)
}

func (m *Mpesa) AccountBalance(params map[string]interface{}) (map[string]interface{}, error) {
	shortcode := params["shortcode"].(string)
	identifierType := int(params["identifierType"].(float64))
	remarks := params["remarks"].(string)
	resultURL, _ := params["resultUrl"].(string)
	timeoutURL, _ := params["timeoutUrl"].(string)
	shortCodeType := "C2B"
	if s, ok := params["shortCodeType"].(string); ok {
		shortCodeType = s
	}

	url := m.baseURL + "/mpesa/accountbalance/v1/query"
	body := map[string]interface{}{
		"Initiator":          m.getConfig("initiator_name"),
		"SecurityCredential": m.generateSecurityCredential(),
		"CommandID":          "AccountBalance",
		"PartyA":             shortcode,
		"IdentifierType":     identifierType,
		"Remarks":            remarks,
		"ResultURL":          m.resolveCallbackUrl(resultURL, "balance_result_url"),
		"QueueTimeOutURL":    m.resolveCallbackUrl(timeoutURL, "balance_timeout_url"),
	}

	return m.mpesaRequest(url, body, shortCodeType)
}

func (m *Mpesa) Reversal(params map[string]interface{}) (map[string]interface{}, error) {
	shortcode := params["shortcode"].(string)
	transactionID := params["transactionId"].(string)
	amount := int(params["amount"].(float64))
	remarks := params["remarks"].(string)
	resultURL, _ := params["resultUrl"].(string)
	timeoutURL, _ := params["timeoutUrl"].(string)
	shortCodeType := "C2B"
	if s, ok := params["shortCodeType"].(string); ok {
		shortCodeType = s
	}

	url := m.baseURL + "/mpesa/reversal/v1/request"
	body := map[string]interface{}{
		"Initiator":              m.getConfig("initiator_name"),
		"SecurityCredential":     m.generateSecurityCredential(),
		"CommandID":              "TransactionReversal",
		"TransactionID":          transactionID,
		"Amount":                 amount,
		"ReceiverParty":          shortcode,
		"RecieverIdentifierType": "11",
		"Remarks":                remarks,
		"Occasion":               "",
		"ResultURL":              m.resolveCallbackUrl(resultURL, "reversal_result_url"),
		"QueueTimeOutURL":        m.resolveCallbackUrl(timeoutURL, "reversal_timeout_url"),
	}

	return m.mpesaRequest(url, body, shortCodeType)
}

func (m *Mpesa) B2pochi(params map[string]interface{}) (map[string]interface{}, error) {
	phonenumber := params["phonenumber"].(string)
	amount := int(params["amount"].(float64))
	remarks := params["remarks"].(string)
	occasion, _ := params["occasion"].(string)
	resultURL, _ := params["resultUrl"].(string)
	timeoutURL, _ := params["timeoutUrl"].(string)
	shortCodeType := "B2C"
	if s, ok := params["shortCodeType"].(string); ok {
		shortCodeType = s
	}

	url := m.baseURL + "/mpesa/b2pochi/v1/paymentrequest"
	body := map[string]interface{}{
		"OriginatorConversationID": m.getFormattedTimestamp(),
		"InitiatorName":            m.getConfig("initiator_name"),
		"SecurityCredential":       m.generateSecurityCredential(),
		"CommandID":                "BusinessPayToPochi",
		"Amount":                   amount,
		"PartyA":                   m.getConfig("b2c_shortcode"),
		"PartyB":                   m.phoneValidator(phonenumber),
		"Remarks":                  remarks,
		"Occassion":                occasion,
		"ResultURL":                m.resolveCallbackUrl(resultURL, "b2pochi_result_url"),
		"QueueTimeOutURL":          m.resolveCallbackUrl(timeoutURL, "b2pochi_timeout_url"),
	}

	return m.mpesaRequest(url, body, shortCodeType)
}
