package mpesa

type B2PochiService struct {
	httpClient *HTTPClient
	auth       *Auth
	helpers    *Helpers
	baseURL    string
}

func NewB2PochiService(httpClient *HTTPClient, auth *Auth, config *Config) *B2PochiService {
	return &B2PochiService{
		httpClient: httpClient,
		auth:       auth,
		helpers:    NewHelpers(config),
		baseURL:    httpClient.GetBaseURL(),
	}
}

func (s *B2PochiService) Send(phonenumber string, amount int, remarks string, occasion string, resultURL string, timeoutURL string, shortCodeType string) (map[string]interface{}, error) {
	url := s.baseURL + "/mpesa/b2pochi/v1/paymentrequest"
	body := map[string]interface{}{
		"OriginatorConversationID": s.helpers.GetFormattedTimestamp(),
		"InitiatorName":            s.helpers.GetConfig("initiator_name"),
		"SecurityCredential":       "",
		"CommandID":                "BusinessPayToPochi",
		"Amount":                   amount,
		"PartyA":                   s.helpers.GetConfig("b2c_shortcode"),
		"PartyB":                   s.helpers.PhoneValidator(phonenumber),
		"Remarks":                  remarks,
		"Occasion":                 occasion,
		"ResultURL":                resultURL,
		"QueueTimeOutURL":          timeoutURL,
	}

	token, err := s.auth.GetAccessToken(shortCodeType)
	if err != nil {
		return nil, err
	}

	return s.httpClient.Post(url, body, token)
}
