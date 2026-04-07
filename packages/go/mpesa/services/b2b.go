package mpesa

type B2BService struct {
	httpClient *HTTPClient
	auth       *Auth
	helpers    *Helpers
	baseURL    string
}

func NewB2BService(httpClient *HTTPClient, auth *Auth, config *Config) *B2BService {
	return &B2BService{
		httpClient: httpClient,
		auth:       auth,
		helpers:    NewHelpers(config),
		baseURL:    httpClient.GetBaseURL(),
	}
}

func (s *B2BService) Send(receiverShortcode string, commandId string, amount int, remarks string, accountNumber string, resultURL string, timeoutURL string, shortCodeType string) (map[string]interface{}, error) {
	if commandId == "BusinessPayBill" && accountNumber == "" {
		return nil, MissingB2BAccountNumber()
	}

	url := s.baseURL + "/mpesa/b2b/v1/paymentrequest"
	body := map[string]interface{}{
		"Initiator":              s.helpers.GetConfig("initiator_name"),
		"SecurityCredential":     "",
		"CommandID":              commandId,
		"SenderIdentifierType":   "4",
		"RecieverIdentifierType": "4",
		"Amount":                 amount,
		"PartyA":                 s.helpers.GetConfig("b2c_shortcode"),
		"PartyB":                 receiverShortcode,
		"AccountReference":       accountNumber,
		"Remarks":                remarks,
		"ResultURL":              resultURL,
		"QueueTimeOutURL":        timeoutURL,
	}

	token, err := s.auth.GetAccessToken(shortCodeType)
	if err != nil {
		return nil, err
	}

	return s.httpClient.Post(url, body, token)
}
