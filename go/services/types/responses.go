package types

type STKPushResult struct {
	CheckoutRequestID   string
	MerchantRequestID   string
	ResponseCode        string
	ResponseDescription string
	CustomerMessage     string
}

type STKQueryResult struct {
	ResponseCode        string
	ResponseDescription string
	MerchantRequestID   string
	CheckoutRequestID   string
	ResultCode          string
	ResultDesc          string
}

type C2BResult struct {
	OriginatorConversationID string
	ResponseCode             string
	ResponseDescription      string
}

type B2CResult struct {
	ConversationID           string
	OriginatorConversationID string
	ResponseCode             string
	ResponseDescription      string
}

type B2BResult struct {
	OriginatorConversationID string
	ConversationID           string
	ResponseCode             string
	ResponseDescription      string
}

type ReversalResult struct {
	OriginatorConversationID string
	ConversationID           string
	ResponseCode             string
	ResponseDescription      string
}

type TransactionStatusResult struct {
	OriginatorConversationID string
	ConversationID           string
	ResponseCode             string
	ResponseDescription      string
}

type AccountBalanceResult struct {
	OriginatorConversationID string
	ConversationID           string
	ResponseCode             string
	ResponseDescription      string
}

type DynamicQRResult struct {
	ResponseCode        string
	RequestID           string
	ResponseDescription string
	QRCode              string
}
