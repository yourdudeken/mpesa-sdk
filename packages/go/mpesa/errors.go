package mpesa

import "fmt"

type MpesaError struct {
	Message string
	Code    int
}

func (e *MpesaError) Error() string {
	return fmt.Sprintf("MpesaError: %s (code: %d)", e.Message, e.Code)
}

func InvalidTransactionType(transactionType string) *MpesaError {
	return &MpesaError{Message: fmt.Sprintf("Invalid transaction type: %s. Use PAYBILL or TILL.", transactionType)}
}

func MissingAccountReference() *MpesaError {
	return &MpesaError{Message: "An Account Reference is required for All transactions."}
}

func MissingTillNumber() *MpesaError {
	return &MpesaError{Message: "Till number is required for Buy Goods transactions."}
}

func MissingCallbackUrl(key string) *MpesaError {
	return &MpesaError{Message: fmt.Sprintf("Ensure you have set the %s in the config or passed as a parameter", key)}
}

func MissingB2BAccountNumber() *MpesaError {
	return &MpesaError{Message: "Account Number is required for BusinessPayBill CommandID"}
}

func AuthenticationFailed(message string) *MpesaError {
	return &MpesaError{Message: message, Code: 401}
}

func APIError(message string, code int) *MpesaError {
	return &MpesaError{Message: message, Code: code}
}
