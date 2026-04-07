package mpesa

type B2CService struct {
	httpClient *HTTPClient
	auth       *Auth
	helpers    *Helpers
	baseURL    string
}

func NewB2CService(httpClient *HTTPClient, auth *Auth, config *Config) *B2CService {
	return &B2CService{
		httpClient: httpClient,
		auth:       auth,
		helpers:    NewHelpers(config),
		baseURL:    httpClient.GetBaseURL(),
	}
}

func (s *B2CService) Send(phonenumber string, commandId string, amount int, remarks string, resultURL string, timeoutURL string, shortCodeType string) (map[string]interface{}, error) {
	url := s.baseURL + "/mpesa/b2c/v1/paymentrequest"
	body := map[string]interface{}{
		"InitiatorName":      s.helpers.GetConfig("initiator_name"),
		"SecurityCredential": "",
		"CommandID":          commandId,
		"Amount":             amount,
		"PartyA":             s.helpers.GetConfig("b2c_shortcode"),
		"PartyB":             s.helpers.PhoneValidator(phonenumber),
		"Remarks":            remarks,
		"Occassion":          "",
		"ResultURL":          resultURL,
		"QueueTimeOutURL":    timeoutURL,
	}

	token, err := s.auth.GetAccessToken(shortCodeType)
	if err != nil {
		return nil, err
	}

	return s.httpClient.Post(url, body, token)
}
