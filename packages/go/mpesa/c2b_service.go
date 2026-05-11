package mpesa

type C2BService struct {
	httpClient *HTTPClient
	auth       *Auth
	helpers    *Helpers
	baseURL    string
}

func NewC2BService(httpClient *HTTPClient, auth *Auth, config *Config) *C2BService {
	return &C2BService{
		httpClient: httpClient,
		auth:       auth,
		helpers:    NewHelpers(config),
		baseURL:    httpClient.GetBaseURL(),
	}
}

func (s *C2BService) RegisterURLS(shortcode string, confirmURL string, validateURL string, shortCodeType string) (map[string]interface{}, error) {
	url := s.baseURL + "/mpesa/c2b/v2/registerurl"
	body := map[string]interface{}{
		"ShortCode":       shortcode,
		"ResponseType":    "Completed",
		"ConfirmationURL": confirmURL,
		"ValidationURL":   validateURL,
	}

	token, err := s.auth.GetAccessToken(shortCodeType)
	if err != nil {
		return nil, err
	}

	return s.httpClient.Post(url, body, token)
}

func (s *C2BService) Simulate(phonenumber string, amount int, shortcode string, commandId string, accountNumber string, shortCodeType string) (map[string]interface{}, error) {
	url := s.baseURL + "/mpesa/c2b/v2/simulate"
	data := map[string]interface{}{
		"Msisdn":    s.helpers.PhoneValidator(phonenumber),
		"Amount":    amount,
		"CommandID": commandId,
		"ShortCode": shortcode,
	}

	if accountNumber != "" {
		data["BillRefNumber"] = accountNumber
	}

	token, err := s.auth.GetAccessToken(shortCodeType)
	if err != nil {
		return nil, err
	}

	return s.httpClient.Post(url, data, token)
}
