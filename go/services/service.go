package services

import (
	"context"
	"fmt"

	"github.com/yourdudeken/mpesa-sdk/go/client"
	svctypes "github.com/yourdudeken/mpesa-sdk/go/services/types"
	"github.com/yourdudeken/mpesa-sdk/go/types"
)

type Service struct {
	client *client.Client
}

func NewService(c *client.Client) *Service {
	return &Service{client: c}
}

func (s *Service) STKPush(ctx context.Context, input svctypes.STKPushInput) (*svctypes.STKPushResult, error) {
	req := types.STKPushRequest{
		BusinessShortCode: input.BusinessShortCode,
		TransactionType:   input.TransactionType,
		Amount:            input.Amount,
		PartyA:            input.PartyA,
		PartyB:            input.PartyB,
		PhoneNumber:       input.PhoneNumber,
		CallBackURL:       input.CallBackURL,
		AccountReference:  input.AccountReference,
		TransactionDesc:   input.TransactionDesc,
	}
	resp, err := s.client.STKPush(ctx, req)
	if err != nil {
		return nil, err
	}
	return &svctypes.STKPushResult{
		CheckoutRequestID:   resp.CheckoutRequestID,
		MerchantRequestID:   resp.MerchantRequestID,
		ResponseCode:        resp.ResponseCode,
		ResponseDescription: resp.ResponseDescription,
		CustomerMessage:     resp.CustomerMessage,
	}, nil
}

func (s *Service) STKQuery(ctx context.Context, input svctypes.STKQueryInput) (*svctypes.STKQueryResult, error) {
	req := types.STKQueryRequest{
		BusinessShortCode: fmt.Sprintf("%d", input.BusinessShortCode),
		CheckoutRequestID: input.CheckoutRequestID,
	}
	resp, err := s.client.STKQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	return &svctypes.STKQueryResult{
		ResponseCode:        resp.ResponseCode,
		ResponseDescription: resp.ResponseDescription,
		MerchantRequestID:   resp.MerchantRequestID,
		CheckoutRequestID:   resp.CheckoutRequestID,
		ResultCode:          resp.ResultCode,
		ResultDesc:          resp.ResultDesc,
	}, nil
}

func (s *Service) C2BRegisterURL(ctx context.Context, input svctypes.C2BRegisterURLInput) (*svctypes.C2BResult, error) {
	req := types.C2BRegisterURLRequest{
		ShortCode:       input.ShortCode,
		ResponseType:    input.ResponseType,
		ConfirmationURL: input.ConfirmationURL,
		ValidationURL:   input.ValidationURL,
	}
	resp, err := s.client.C2BRegisterURL(ctx, req)
	if err != nil {
		return nil, err
	}
	return &svctypes.C2BResult{
		OriginatorConversationID: resp.OriginatorCoversationID,
		ResponseCode:             resp.ResponseCode,
		ResponseDescription:      resp.ResponseDescription,
	}, nil
}

func (s *Service) C2BSimulate(ctx context.Context, input svctypes.C2BSimulateInput) (*svctypes.C2BResult, error) {
	req := types.C2BSimulateRequest{
		ShortCode:     input.ShortCode,
		CommandID:     input.CommandID,
		Amount:        input.Amount,
		Msisdn:        input.Msisdn,
		BillRefNumber: input.BillRefNumber,
	}
	resp, err := s.client.C2BSimulate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &svctypes.C2BResult{
		OriginatorConversationID: resp.OriginatorCoversationID,
		ResponseCode:             resp.ResponseCode,
		ResponseDescription:      resp.ResponseDescription,
	}, nil
}

func (s *Service) B2C(ctx context.Context, input svctypes.B2CInput) (*svctypes.B2CResult, error) {
	req := types.B2CRequest{
		InitiatorName:      input.InitiatorName,
		SecurityCredential: input.SecurityCredential,
		CommandID:          input.CommandID,
		Amount:             input.Amount,
		PartyA:             input.PartyA,
		PartyB:             input.PartyB,
		Remarks:            input.Remarks,
		QueueTimeOutURL:    input.QueueTimeOutURL,
		ResultURL:          input.ResultURL,
		Occassion:          input.Occassion,
	}
	resp, err := s.client.B2C(ctx, req)
	if err != nil {
		return nil, err
	}
	return &svctypes.B2CResult{
		ConversationID:           resp.ConversationID,
		OriginatorConversationID: resp.OriginatorConversationID,
		ResponseCode:             resp.ResponseCode,
		ResponseDescription:      resp.ResponseDescription,
	}, nil
}

func (s *Service) B2B(ctx context.Context, input svctypes.B2BInput) (*svctypes.B2BResult, error) {
	req := types.B2BRequest{
		Initiator:              input.Initiator,
		SecurityCredential:     input.SecurityCredential,
		CommandID:              input.CommandID,
		SenderIdentifierType:   input.SenderIdentifierType,
		RecieverIdentifierType: input.RecieverIdentifierType,
		Amount:                 input.Amount,
		PartyA:                 input.PartyA,
		PartyB:                 input.PartyB,
		Requester:              input.Requester,
		AccountReference:       input.AccountReference,
		Remarks:                input.Remarks,
		QueueTimeOutURL:        input.QueueTimeOutURL,
		ResultURL:              input.ResultURL,
		Occassion:              input.Occassion,
	}
	resp, err := s.client.B2B(ctx, req)
	if err != nil {
		return nil, err
	}
	return &svctypes.B2BResult{
		OriginatorConversationID: resp.OriginatorConversationID,
		ConversationID:           resp.ConversationID,
		ResponseCode:             resp.ResponseCode,
		ResponseDescription:      resp.ResponseDescription,
	}, nil
}

func (s *Service) Reversal(ctx context.Context, input svctypes.ReversalInput) (*svctypes.ReversalResult, error) {
	req := types.ReversalRequest{
		Initiator:          input.Initiator,
		SecurityCredential: input.SecurityCredential,
		CommandID:          "TransactionReversal",
		TransactionID:      input.TransactionID,
		Amount:             input.Amount,
		ReceiverParty:      input.ReceiverParty,
		QueueTimeOutURL:    input.QueueTimeOutURL,
		ResultURL:          input.ResultURL,
		Remarks:            input.Remarks,
	}
	resp, err := s.client.Reversal(ctx, req)
	if err != nil {
		return nil, err
	}
	return &svctypes.ReversalResult{
		OriginatorConversationID: resp.OriginatorConversationID,
		ConversationID:           resp.ConversationID,
		ResponseCode:             resp.ResponseCode,
		ResponseDescription:      resp.ResponseDescription,
	}, nil
}

func (s *Service) TransactionStatus(ctx context.Context, input svctypes.TransactionStatusInput) (*svctypes.TransactionStatusResult, error) {
	req := types.TransactionStatusRequest{
		Initiator:              input.Initiator,
		SecurityCredential:     input.SecurityCredential,
		CommandID:              "TransactionStatusQuery",
		TransactionID:          input.TransactionID,
		OriginalConversationID: input.OriginalConversationID,
		PartyA:                 input.PartyA,
		IdentifierType:         4,
		ResultURL:              input.ResultURL,
		QueueTimeOutURL:        input.QueueTimeOutURL,
		Remarks:                input.Remarks,
	}
	resp, err := s.client.TransactionStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	return &svctypes.TransactionStatusResult{
		OriginatorConversationID: resp.OriginatorConversationID,
		ConversationID:           resp.ConversationID,
		ResponseCode:             resp.ResponseCode,
		ResponseDescription:      resp.ResponseDescription,
	}, nil
}

func (s *Service) AccountBalance(ctx context.Context, input svctypes.AccountBalanceInput) (*svctypes.AccountBalanceResult, error) {
	req := types.AccountBalanceRequest{
		Initiator:          input.Initiator,
		SecurityCredential: input.SecurityCredential,
		CommandID:          "AccountBalance",
		PartyA:             input.PartyA,
		IdentifierType:     4,
		Remarks:            input.Remarks,
		QueueTimeOutURL:    input.QueueTimeOutURL,
		ResultURL:          input.ResultURL,
	}
	resp, err := s.client.AccountBalance(ctx, req)
	if err != nil {
		return nil, err
	}
	return &svctypes.AccountBalanceResult{
		OriginatorConversationID: resp.OriginatorConversationID,
		ConversationID:           resp.ConversationID,
		ResponseCode:             resp.ResponseCode,
		ResponseDescription:      resp.ResponseDescription,
	}, nil
}

func (s *Service) DynamicQR(ctx context.Context, input svctypes.DynamicQRInput) (*svctypes.DynamicQRResult, error) {
	req := types.DynamicQRRequest{
		MerchantName: input.MerchantName,
		RefNo:        input.RefNo,
		Amount:       input.Amount,
		TrxCode:      input.TrxCode,
		CPI:          input.CPI,
		Size:         input.Size,
	}
	resp, err := s.client.DynamicQR(ctx, req)
	if err != nil {
		return nil, err
	}
	return &svctypes.DynamicQRResult{
		ResponseCode:        resp.ResponseCode,
		RequestID:           resp.RequestID,
		ResponseDescription: resp.ResponseDescription,
		QRCode:              resp.QRCode,
	}, nil
}
