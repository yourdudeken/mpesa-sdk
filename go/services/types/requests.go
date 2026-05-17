package types

import (
	"github.com/yourdudeken/mpesa-sdk/go/types"
)

type STKPushInput struct {
	BusinessShortCode int
	TransactionType   types.TransactionType
	Amount            int
	PartyA            int
	PartyB            int
	PhoneNumber       int
	CallBackURL       string
	AccountReference  string
	TransactionDesc   string
}

type STKQueryInput struct {
	BusinessShortCode int
	CheckoutRequestID string
}

type C2BRegisterURLInput struct {
	ShortCode       string
	ResponseType    types.ResponseType
	ConfirmationURL string
	ValidationURL   string
}

type C2BSimulateInput struct {
	ShortCode     int
	CommandID     types.C2BCommandID
	Amount        int
	Msisdn        int
	BillRefNumber string
}

type B2CInput struct {
	InitiatorName      string
	SecurityCredential string
	CommandID          types.B2CCommandID
	Amount             int
	PartyA             int
	PartyB             int
	Remarks            string
	QueueTimeOutURL    string
	ResultURL          string
	Occassion          string
}

type B2BInput struct {
	Initiator              string
	SecurityCredential     string
	CommandID              types.B2BCommandID
	SenderIdentifierType   int
	RecieverIdentifierType int
	Amount                 int
	PartyA                 int
	PartyB                 int
	Requester              int
	AccountReference       string
	Remarks                string
	QueueTimeOutURL        string
	ResultURL              string
	Occassion              string
}

type ReversalInput struct {
	Initiator          string
	SecurityCredential string
	TransactionID      string
	Amount             int
	ReceiverParty      int
	QueueTimeOutURL    string
	ResultURL          string
	Remarks            string
}

type TransactionStatusInput struct {
	Initiator              string
	SecurityCredential     string
	TransactionID          string
	OriginalConversationID string
	PartyA                 int
	ResultURL              string
	QueueTimeOutURL        string
	Remarks                string
}

type AccountBalanceInput struct {
	Initiator          string
	SecurityCredential string
	PartyA             int
	Remarks            string
	QueueTimeOutURL    string
	ResultURL          string
}

type DynamicQRInput struct {
	MerchantName string
	RefNo        string
	Amount       int
	TrxCode      types.TrxCode
	CPI          string
	Size         string
}
