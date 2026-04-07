package mpesa

import "fmt"

type STKPushService struct {
	httpClient *HTTPClient
	auth       *Auth
	helpers    *Helpers
	baseURL    string
}

func NewSTKPushService(httpClient *HTTPClient, auth *Auth, config *Config) *STKPushService {
	return &STKPushService{
		httpClient: httpClient,
		auth:       auth,
		helpers:    NewHelpers(config),
		baseURL:    httpClient.GetBaseURL(),
	}
}

func (s *STKPushService) Push(phonenumber string, amount int, accountNumber string, callbackURL string, transactionType string, shortCodeType string) (map[string]interface{}, error) {
	if accountNumber == "" {
		return nil, MissingAccountReference()
	}

	url := s.baseURL + "/mpesa/stkpush/v1/processrequest"
	data := map[string]interface{}{
		"BusinessShortCode": s.helpers.GetConfig("shortcode"),
		"Password":          s.helpers.LipaNaMpesaPassword(),
		"Timestamp":         s.helpers.GetFormattedTimestamp(),
		"Amount":            amount,
		"PartyA":            s.helpers.PhoneValidator(phonenumber),
		"PartyB":            s.helpers.GetConfig("shortcode"),
		"TransactionType":   transactionType,
		"PhoneNumber":       s.helpers.PhoneValidator(phonenumber),
		"TransactionDesc":   "Payment",
		"AccountReference":  accountNumber,
		"CallBackURL":       callbackURL,
	}

	token, err := s.auth.GetAccessToken(shortCodeType)
	if err != nil {
		return nil, err
	}

	return s.httpClient.Post(url, data, token)
}

func (s *STKPushService) Query(checkoutRequestID string, shortCodeType string) (map[string]interface{}, error) {
	url := s.baseURL + "/mpesa/stkpushquery/v1/query"
	data := map[string]interface{}{
		"BusinessShortCode": s.helpers.GetConfig("shortcode"),
		"Password":          s.helpers.LipaNaMpesaPassword(),
		"Timestamp":         s.helpers.GetFormattedTimestamp(),
		"CheckoutRequestID": checkoutRequestID,
	}

	token, err := s.auth.GetAccessToken(shortCodeType)
	if err != nil {
		return nil, err
	}

	return s.httpClient.Post(url, data, token)
}

func (s *STKPushService) Validate(accountNumber string) error {
	if accountNumber == "" {
		return fmt.Errorf("An Account Reference is required for All transactions.")
	}
	return nil
}
