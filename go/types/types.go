package types

import "time"

type Environment string

const (
	Sandbox    Environment = "sandbox"
	Production Environment = "production"
)

type RetryConfig struct {
	MaxRetries  int
	BaseDelayMs int
	MaxDelayMs  int
}

type MpesaConfig struct {
	ConsumerKey        string
	ConsumerSecret     string
	Environment        Environment
	Passkey            string
	InitiatorName      string
	InitiatorPassword  string
	SecurityCredential string
	Timeout            time.Duration
	RetryConfig        RetryConfig
}

// ---- Auth ----
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenCache struct {
	Token     string
	ExpiresAt time.Time
}

// ---- STK Push ----
type TransactionType string

const (
	CustomerPayBillOnline  TransactionType = "CustomerPayBillOnline"
	CustomerBuyGoodsOnline TransactionType = "CustomerBuyGoodsOnline"
)

type STKPushRequest struct {
	BusinessShortCode int             `json:"BusinessShortCode"`
	Password          string          `json:"Password"`
	Timestamp         string          `json:"Timestamp"`
	TransactionType   TransactionType `json:"TransactionType"`
	Amount            int             `json:"Amount"`
	PartyA            int             `json:"PartyA"`
	PartyB            int             `json:"PartyB"`
	PhoneNumber       int             `json:"PhoneNumber"`
	CallBackURL       string          `json:"CallBackURL"`
	AccountReference  string          `json:"AccountReference"`
	TransactionDesc   string          `json:"TransactionDesc"`
}

type STKPushResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

type STKQueryRequest struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	CheckoutRequestID string `json:"CheckoutRequestID"`
}

type STKQueryResponse struct {
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResultCode          string `json:"ResultCode"`
	ResultDesc          string `json:"ResultDesc"`
}

type STKCallbackPayload struct {
	Body struct {
		StkCallback struct {
			MerchantRequestID string `json:"MerchantRequestID"`
			CheckoutRequestID string `json:"CheckoutRequestID"`
			ResultCode        int    `json:"ResultCode"`
			ResultDesc        string `json:"ResultDesc"`
			CallbackMetadata  *struct {
				Item []struct {
					Name  string      `json:"Name"`
					Value interface{} `json:"Value"`
				} `json:"Item"`
			} `json:"CallbackMetadata"`
		} `json:"stkCallback"`
	} `json:"Body"`
}

type STKCallbackResult struct {
	Success           bool
	MerchantRequestID string
	CheckoutRequestID string
	ResultCode        int
	ResultDescription string
	Amount            *float64
	ReceiptNumber     *string
	TransactionDate   *string
	PhoneNumber       *string
}

// ---- C2B ----
type ResponseType string

const (
	ResponseCompleted ResponseType = "Completed"
	ResponseCancelled ResponseType = "Cancelled"
)

type C2BCommandID string

const (
	C2BPayBill  C2BCommandID = "CustomerPayBillOnline"
	C2BBuyGoods C2BCommandID = "CustomerBuyGoodsOnline"
)

type C2BRegisterURLRequest struct {
	ShortCode       string       `json:"ShortCode"`
	ResponseType    ResponseType `json:"ResponseType"`
	ConfirmationURL string       `json:"ConfirmationURL"`
	ValidationURL   string       `json:"ValidationURL"`
}

type C2BSimulateRequest struct {
	ShortCode     int          `json:"ShortCode"`
	CommandID     C2BCommandID `json:"CommandID"`
	Amount        int          `json:"Amount"`
	Msisdn        int          `json:"Msisdn"`
	BillRefNumber string       `json:"BillRefNumber,omitempty"`
}

type C2BResponse struct {
	OriginatorCoversationID string `json:"OriginatorCoversationID"`
	ResponseCode            string `json:"ResponseCode"`
	ResponseDescription     string `json:"ResponseDescription"`
}

// ---- B2C ----
type B2CCommandID string

const (
	SalaryPayment    B2CCommandID = "SalaryPayment"
	BusinessPayment  B2CCommandID = "BusinessPayment"
	PromotionPayment B2CCommandID = "PromotionPayment"
)

type B2CRequest struct {
	OriginatorConversationID string       `json:"OriginatorConversationID,omitempty"`
	InitiatorName            string       `json:"InitiatorName"`
	SecurityCredential       string       `json:"SecurityCredential"`
	CommandID                B2CCommandID `json:"CommandID"`
	Amount                   int          `json:"Amount"`
	PartyA                   int          `json:"PartyA"`
	PartyB                   int          `json:"PartyB"`
	Remarks                  string       `json:"Remarks"`
	QueueTimeOutURL          string       `json:"QueueTimeOutURL"`
	ResultURL                string       `json:"ResultURL"`
	Occassion                string       `json:"Occassion,omitempty"`
}

type B2CResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

// ---- B2B ----
type B2BCommandID string

const (
	BusinessPayBill                       B2BCommandID = "BusinessPayBill"
	BusinessBuyGoods                      B2BCommandID = "BusinessBuyGoods"
	MerchantToMerchantTransfer            B2BCommandID = "MerchantToMerchantTransfer"
	MerchantTransferFromMerchantToWorking B2BCommandID = "MerchantTransferFromMerchantToWorking"
	MerchantServicesMMFAccountBalance     B2BCommandID = "MerchantServicesMMFAccountBalance"
	AgencyFloatAdvance                    B2BCommandID = "AgencyFloatAdvance"
)

type B2BRequest struct {
	Initiator              string       `json:"Initiator"`
	SecurityCredential     string       `json:"SecurityCredential"`
	CommandID              B2BCommandID `json:"CommandID"`
	SenderIdentifierType   int          `json:"SenderIdentifierType"`
	RecieverIdentifierType int          `json:"RecieverIdentifierType"`
	Amount                 int          `json:"Amount"`
	PartyA                 int          `json:"PartyA"`
	PartyB                 int          `json:"PartyB"`
	Requester              int          `json:"Requester,omitempty"`
	AccountReference       string       `json:"AccountReference,omitempty"`
	Remarks                string       `json:"Remarks"`
	QueueTimeOutURL        string       `json:"QueueTimeOutURL"`
	ResultURL              string       `json:"ResultURL"`
	Occassion              string       `json:"Occassion,omitempty"`
}

type B2BResponse struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

// ---- Reversal ----
type ReversalRequest struct {
	Initiator              string `json:"Initiator"`
	SecurityCredential     string `json:"SecurityCredential"`
	CommandID              string `json:"CommandID"`
	TransactionID          string `json:"TransactionID"`
	Amount                 int    `json:"Amount"`
	ReceiverParty          int    `json:"ReceiverParty"`
	RecieverIdentifierType int    `json:"RecieverIdentifierType"`
	QueueTimeOutURL        string `json:"QueueTimeOutURL"`
	ResultURL              string `json:"ResultURL"`
	Remarks                string `json:"Remarks"`
}

type ReversalResponse struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

// ---- Transaction Status ----
type TransactionStatusRequest struct {
	Initiator              string `json:"Initiator"`
	SecurityCredential     string `json:"SecurityCredential"`
	CommandID              string `json:"CommandID"`
	TransactionID          string `json:"TransactionID,omitempty"`
	OriginalConversationID string `json:"OriginalConversationID,omitempty"`
	PartyA                 int    `json:"PartyA"`
	IdentifierType         int    `json:"IdentifierType"`
	ResultURL              string `json:"ResultURL"`
	QueueTimeOutURL        string `json:"QueueTimeOutURL"`
	Remarks                string `json:"Remarks"`
	Occasion               string `json:"Occasion,omitempty"`
}

type TransactionStatusResponse struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

// ---- Account Balance ----
type AccountBalanceRequest struct {
	Initiator          string `json:"Initiator"`
	SecurityCredential string `json:"SecurityCredential"`
	CommandID          string `json:"CommandID"`
	PartyA             int    `json:"PartyA"`
	IdentifierType     int    `json:"IdentifierType"`
	Remarks            string `json:"Remarks"`
	QueueTimeOutURL    string `json:"QueueTimeOutURL"`
	ResultURL          string `json:"ResultURL"`
}

type AccountBalanceResponse struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

// ---- Dynamic QR ----
type TrxCode string

const (
	TrxBuyGoods       TrxCode = "BG"
	TrxWithdrawCash   TrxCode = "WA"
	TrxPaybill        TrxCode = "PB"
	TrxSendMoney      TrxCode = "SM"
	TrxSendToBusiness TrxCode = "SB"
)

type DynamicQRRequest struct {
	MerchantName string  `json:"MerchantName"`
	RefNo        string  `json:"RefNo"`
	Amount       int     `json:"Amount"`
	TrxCode      TrxCode `json:"TrxCode"`
	CPI          string  `json:"CPI"`
	Size         string  `json:"Size"`
}

type DynamicQRResponse struct {
	ResponseCode        string `json:"ResponseCode"`
	RequestID           string `json:"RequestID"`
	ResponseDescription string `json:"ResponseDescription"`
	QRCode              string `json:"QRCode"`
}

// ---- Shared ----
type ResultParameterItem struct {
	Key   string      `json:"Key"`
	Value interface{} `json:"Value"`
}

type ResultParameters struct {
	ResultParameter []ResultParameterItem `json:"ResultParameter"`
}

type ReferenceItem struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type ReferenceData struct {
	ReferenceItem *ReferenceItem `json:"ReferenceItem"`
}

type ResultDetail struct {
	ResultType               int               `json:"ResultType"`
	ResultCode               int               `json:"ResultCode"`
	ResultDesc               string            `json:"ResultDesc"`
	OriginatorConversationID string            `json:"OriginatorConversationID"`
	ConversationID           string            `json:"ConversationID"`
	TransactionID            string            `json:"TransactionID"`
	ResultParameters         *ResultParameters `json:"ResultParameters,omitempty"`
	ReferenceData            *ReferenceData    `json:"ReferenceData,omitempty"`
}

type MpesaResult struct {
	Result ResultDetail `json:"Result"`
}

type AccountInfo struct {
	AccountName      string
	Currency         string
	AvailableBalance float64
	UnclearedFunds   float64
	ReservedFunds    float64
}

type AccountBalanceResult struct {
	WorkingAccount            *AccountInfo
	UtilityAccount            *AccountInfo
	ChargesPaidAccount        *AccountInfo
	OrganizationSettlementAcc *AccountInfo
	FloatAccount              *AccountInfo
}
