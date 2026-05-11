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

type Mpesa struct {
	config  *Config
	baseURL string
	auth    *Auth
}

func NewClient(config *Config) *Mpesa {
	baseURL := "https://sandbox.safaricom.co.ke"
	if config.Environment == "production" {
		baseURL = "https://api.safaricom.co.ke"
	}
	c := &Mpesa{
		config:  config,
		baseURL: baseURL,
	}
	c.auth = NewAuth(config, &c.baseURL)
	return c
}

func (m *Mpesa) resolveCallback(paramURL, configKey string) string {
	if paramURL != "" {
		return paramURL
	}
	if m.config.Callbacks != nil {
		if url, ok := m.config.Callbacks[configKey]; ok {
			return url
		}
	}
	return ""
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
	password := m.config.Shortcode + m.config.Passkey + timestamp
	return base64.StdEncoding.EncodeToString([]byte(password))
}

func (m *Mpesa) securityCredential() string {
	cred, err := m.generateSecurityCredential()
	if err != nil {
		return ""
	}
	return cred
}

func (m *Mpesa) generateSecurityCredential() (string, error) {

	certPath := "certificates/SandboxCertificate.cer"
	if m.config.Environment == "production" {
		certPath = "certificates/ProductionCertificate.cer"
	}

	wd, _ := os.Getwd()
	paths := []string{
		certPath,
		filepath.Join(wd, certPath),
	}

	var pubKeyData []byte
	var err error
	for _, p := range paths {
		pubKeyData, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", fmt.Errorf("failed to read certificate: %v", err)
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

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, []byte(m.config.InitiatorPassword))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (m *Mpesa) mpesaRequest(url string, body map[string]interface{}, shortCodeType string) (map[string]interface{}, error) {
	token, err := m.auth.GetAccessToken(shortCodeType)
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

func (m *Mpesa) Stkpush(phonenumber string, amount int, accountNumber string, callbackURL string, transactionType ...string) (map[string]interface{}, error) {
	if accountNumber == "" {
		return nil, fmt.Errorf("an Account Reference is required for All transactions")
	}

	tType := Paybill
	if len(transactionType) > 0 {
		tType = transactionType[0]
	}

	if tType == Till && m.config.TillNumber == "" {
		return nil, fmt.Errorf("till number is required for Buy Goods transactions")
	}

	url := m.baseURL + "/mpesa/stkpush/v1/processrequest"
	partyB := m.config.Shortcode
	if tType == Till {
		partyB = m.config.TillNumber
	}

	data := map[string]interface{}{
		"BusinessShortCode": m.config.Shortcode,
		"Password":          m.lipaNaMpesaPassword(),
		"Timestamp":         m.getFormattedTimestamp(),
		"Amount":            amount,
		"PartyA":            m.phoneValidator(phonenumber),
		"PartyB":            partyB,
		"TransactionType":   tType,
		"PhoneNumber":       m.phoneValidator(phonenumber),
		"TransactionDesc":   "Payment",
		"AccountReference":  accountNumber,
		"CallBackURL":       m.resolveCallback(callbackURL, "callback_url"),
	}

	return m.mpesaRequest(url, data, "C2B")
}

func (m *Mpesa) Stkquery(checkoutRequestID string, callbackURL string) (map[string]interface{}, error) {
	url := m.baseURL + "/mpesa/stkpushquery/v1/query"
	data := map[string]interface{}{
		"BusinessShortCode": m.config.Shortcode,
		"Password":          m.lipaNaMpesaPassword(),
		"Timestamp":         m.getFormattedTimestamp(),
		"CheckoutRequestID": checkoutRequestID,
	}

	return m.mpesaRequest(url, data, "C2B")
}

func (m *Mpesa) B2c(phonenumber string, commandId string, amount int, remarks string) (map[string]interface{}, error) {
	url := m.baseURL + "/mpesa/b2c/v1/paymentrequest"
	body := map[string]interface{}{
		"InitiatorName":      m.config.InitiatorName,
		"SecurityCredential": m.securityCredential(),
		"CommandID":          commandId,
		"Amount":             amount,
		"PartyA":             m.config.B2cShortcode,
		"PartyB":             m.phoneValidator(phonenumber),
		"Remarks":            remarks,
		"Occassion":          "",
		"ResultURL":          m.resolveCallback("", "b2c_result_url"),
		"QueueTimeOutURL":    m.resolveCallback("", "b2c_timeout_url"),
	}

	return m.mpesaRequest(url, body, "B2C")
}

func (m *Mpesa) Validated_b2c(phonenumber string, commandId string, amount int, remarks string, idNumber string) (map[string]interface{}, error) {
	url := m.baseURL + "/mpesa/b2cvalidate/v2/paymentrequest"
	body := map[string]interface{}{
		"InitiatorName":            m.config.InitiatorName,
		"SecurityCredential":       m.securityCredential(),
		"CommandID":                commandId,
		"Amount":                   amount,
		"PartyA":                   m.config.B2cShortcode,
		"PartyB":                   m.phoneValidator(phonenumber),
		"Remarks":                  remarks,
		"Occassion":                "",
		"OriginatorConversationID": m.getFormattedTimestamp(),
		"IDType":                   "01",
		"IDNumber":                 idNumber,
		"ResultURL":                m.resolveCallback("", "b2c_result_url"),
		"QueueTimeOutURL":          m.resolveCallback("", "b2c_timeout_url"),
	}

	return m.mpesaRequest(url, body, "B2C")
}

func (m *Mpesa) B2b(receiverShortcode string, commandId string, amount int, remarks string, accountNumber string) (map[string]interface{}, error) {
	if commandId == "BusinessPayBill" && accountNumber == "" {
		return nil, fmt.Errorf("Account Number is required for BusinessPayBill CommandID")
	}

	url := m.baseURL + "/mpesa/b2b/v1/paymentrequest"
	body := map[string]interface{}{
		"Initiator":              m.config.InitiatorName,
		"SecurityCredential":     m.securityCredential(),
		"CommandID":              commandId,
		"SenderIdentifierType":   "4",
		"RecieverIdentifierType": "4",
		"Amount":                 amount,
		"PartyA":                 m.config.B2cShortcode,
		"PartyB":                 receiverShortcode,
		"AccountReference":       accountNumber,
		"Remarks":                remarks,
		"ResultURL":              m.resolveCallback("", "b2b_result_url"),
		"QueueTimeOutURL":        m.resolveCallback("", "b2b_timeout_url"),
	}

	return m.mpesaRequest(url, body, "B2B")
}

func (m *Mpesa) C2bregisterURLS(shortcode string, confirmUrl string, validateUrl string) (map[string]interface{}, error) {
	url := m.baseURL + "/mpesa/c2b/v2/registerurl"
	body := map[string]interface{}{
		"ShortCode":       shortcode,
		"ResponseType":    "Completed",
		"ConfirmationURL": m.resolveCallback(confirmUrl, "c2b_confirmation_url"),
		"ValidationURL":   m.resolveCallback(validateUrl, "c2b_validation_url"),
	}

	return m.mpesaRequest(url, body, "C2B")
}

func (m *Mpesa) C2bsimulate(phonenumber string, amount int, shortcode string, commandId string, accountNumber string) (map[string]interface{}, error) {
	url := m.baseURL + "/mpesa/c2b/v2/simulate"
	data := map[string]interface{}{
		"Msisdn":    m.phoneValidator(phonenumber),
		"Amount":    amount,
		"CommandID": commandId,
		"ShortCode": shortcode,
	}

	if commandId == Paybill && accountNumber != "" {
		data["BillRefNumber"] = accountNumber
	}

	return m.mpesaRequest(url, data, "C2B")
}

func (m *Mpesa) AccountBalance(shortcode string, identifierType int, remarks string) (map[string]interface{}, error) {
	url := m.baseURL + "/mpesa/accountbalance/v1/query"
	body := map[string]interface{}{
		"Initiator":          m.config.InitiatorName,
		"SecurityCredential": m.securityCredential(),
		"CommandID":          "AccountBalance",
		"PartyA":             shortcode,
		"IdentifierType":     identifierType,
		"Remarks":            remarks,
		"ResultURL":          m.resolveCallback("", "balance_result_url"),
		"QueueTimeOutURL":    m.resolveCallback("", "balance_timeout_url"),
	}

	return m.mpesaRequest(url, body, "C2B")
}

func (m *Mpesa) TransactionStatus(shortcode string, transactionId string, identifierType int, remarks string) (map[string]interface{}, error) {
	url := m.baseURL + "/mpesa/transactionstatus/v1/query"
	body := map[string]interface{}{
		"Initiator":          m.config.InitiatorName,
		"SecurityCredential": m.securityCredential(),
		"CommandID":          "TransactionStatusQuery",
		"TransactionID":      transactionId,
		"PartyA":             shortcode,
		"IdentifierType":     identifierType,
		"Remarks":            remarks,
		"Occassion":          "",
		"ResultURL":          m.resolveCallback("", "status_result_url"),
		"QueueTimeOutURL":    m.resolveCallback("", "status_timeout_url"),
	}

	return m.mpesaRequest(url, body, "C2B")
}

func (m *Mpesa) Reversal(shortcode string, transactionId string, amount int, remarks string) (map[string]interface{}, error) {
	url := m.baseURL + "/mpesa/reversal/v1/request"
	body := map[string]interface{}{
		"Initiator":              m.config.InitiatorName,
		"SecurityCredential":     m.securityCredential(),
		"CommandID":              "TransactionReversal",
		"TransactionID":          transactionId,
		"Amount":                 amount,
		"ReceiverParty":          shortcode,
		"RecieverIdentifierType": "11",
		"Remarks":                remarks,
		"Occasion":               "",
		"ResultURL":              m.resolveCallback("", "reversal_result_url"),
		"QueueTimeOutURL":        m.resolveCallback("", "reversal_timeout_url"),
	}

	return m.mpesaRequest(url, body, "C2B")
}

func (m *Mpesa) B2pochi(phonenumber string, amount int, remarks string) (map[string]interface{}, error) {
	url := m.baseURL + "/mpesa/b2pochi/v1/paymentrequest"
	body := map[string]interface{}{
		"OriginatorConversationID": m.getFormattedTimestamp(),
		"InitiatorName":            m.config.InitiatorName,
		"SecurityCredential":       m.securityCredential(),
		"CommandID":                "BusinessPayToPochi",
		"Amount":                   amount,
		"PartyA":                   m.config.B2cShortcode,
		"PartyB":                   m.phoneValidator(phonenumber),
		"Remarks":                  remarks,
		"Occassion":                "",
		"ResultURL":                m.resolveCallback("", "b2pochi_result_url"),
		"QueueTimeOutURL":          m.resolveCallback("", "b2pochi_timeout_url"),
	}

	return m.mpesaRequest(url, body, "B2C")
}
