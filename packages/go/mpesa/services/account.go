package mpesa

type AccountService struct {
	httpClient *HTTPClient
	auth       *Auth
	helpers    *Helpers
	baseURL    string
}

func NewAccountService(httpClient *HTTPClient, auth *Auth, config *Config) *AccountService {
	return &AccountService{
		httpClient: httpClient,
		auth:       auth,
		helpers:    NewHelpers(config),
		baseURL:    httpClient.GetBaseURL(),
	}
}

func (s *AccountService) Balance(shortcode string, identifierType int, remarks string, resultURL string, timeoutURL string, shortCodeType string) (map[string]interface{}, error) {
	url := s.baseURL + "/mpesa/accountbalance/v1/query"
	body := map[string]interface{}{
		"Initiator":          s.helpers.GetConfig("initiator_name"),
		"SecurityCredential": "",
		"CommandID":          "AccountBalance",
		"PartyA":             shortcode,
		"IdentifierType":     identifierType,
		"Remarks":            remarks,
		"ResultURL":          resultURL,
		"QueueTimeOutURL":    timeoutURL,
	}

	token, err := s.auth.GetAccessToken(shortCodeType)
	if err != nil {
		return nil, err
	}

	return s.httpClient.Post(url, body, token)
}

func (s *AccountService) Status(shortcode string, transactionID string, identifierType int, remarks string, resultURL string, timeoutURL string, shortCodeType string) (map[string]interface{}, error) {
	url := s.baseURL + "/mpesa/transactionstatus/v1/query"
	body := map[string]interface{}{
		"Initiator":          s.helpers.GetConfig("initiator_name"),
		"SecurityCredential": "",
		"CommandID":          "TransactionStatusQuery",
		"TransactionID":      transactionID,
		"PartyA":             shortcode,
		"IdentifierType":     identifierType,
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

func (s *AccountService) Reversal(shortcode string, transactionID string, amount float64, remarks string, resultURL string, timeoutURL string, shortCodeType string) (map[string]interface{}, error) {
	url := s.baseURL + "/mpesa/reversal/v1/request"
	body := map[string]interface{}{
		"Initiator":              s.helpers.GetConfig("initiator_name"),
		"SecurityCredential":     "",
		"CommandID":              "TransactionReversal",
		"TransactionID":          transactionID,
		"Amount":                 amount,
		"ReceiverParty":          shortcode,
		"RecieverIdentifierType": "11",
		"Remarks":                remarks,
		"Occasion":               "",
		"ResultURL":              resultURL,
		"QueueTimeOutURL":        timeoutURL,
	}

	token, err := s.auth.GetAccessToken(shortCodeType)
	if err != nil {
		return nil, err
	}

	return s.httpClient.Post(url, body, token)
}
