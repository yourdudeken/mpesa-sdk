package main

import (
	"fmt"
	"os"

	"github.com/yourdudeken/mpesa-sdk/mpesa"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	config := &mpesa.Config{
		Environment:         env("MPESA_ENV", "sandbox"),
		MpesaConsumerKey:    os.Getenv("MPESA_CONSUMER_KEY"),
		MpesaConsumerSecret: os.Getenv("MPESA_CONSUMER_SECRET"),
		Passkey:             os.Getenv("MPESA_PASSKEY"),
		Shortcode:           env("MPESA_SHORTCODE", "174379"),
		InitiatorName:       env("MPESA_INITIATOR_NAME", "testapi"),
		InitiatorPassword:   os.Getenv("MPESA_INITIATOR_PASSWORD"),
		B2cShortcode:        os.Getenv("MPESA_B2C_SHORTCODE"),
		TillNumber:          os.Getenv("MPESA_TILL_NUMBER"),
		B2cConsumerKey:      os.Getenv("MPESA_B2C_CONSUMER_KEY"),
		B2cConsumerSecret:   os.Getenv("MPESA_B2C_CONSUMER_SECRET"),
		Callbacks: map[string]string{
			"callback_url":         env("MPESA_CALLBACK_URL", ""),
			"b2c_result_url":       env("MPESA_B2C_RESULT_URL", ""),
			"b2c_timeout_url":      env("MPESA_B2C_TIMEOUT_URL", ""),
			"b2b_result_url":       env("MPESA_B2B_RESULT_URL", ""),
			"b2b_timeout_url":      env("MPESA_B2B_TIMEOUT_URL", ""),
			"b2pochi_result_url":   env("MPESA_B2POCHI_RESULT_URL", ""),
			"b2pochi_timeout_url":  env("MPESA_B2POCHI_TIMEOUT_URL", ""),
			"c2b_validation_url":   env("MPESA_C2B_VALIDATION_URL", ""),
			"c2b_confirmation_url": env("MPESA_C2B_CONFIRMATION_URL", ""),
			"balance_result_url":   env("MPESA_BALANCE_RESULT_URL", ""),
			"balance_timeout_url":  env("MPESA_BALANCE_TIMEOUT_URL", ""),
			"status_result_url":    env("MPESA_STATUS_RESULT_URL", ""),
			"status_timeout_url":   env("MPESA_STATUS_TIMEOUT_URL", ""),
			"reversal_result_url":  env("MPESA_REVERSAL_RESULT_URL", ""),
			"reversal_timeout_url": env("MPESA_REVERSAL_TIMEOUT_URL", ""),
		},
	}

	client := mpesa.NewClient(config)

	// STK Push – Lipa na Mpesa Online
	stkResult, err := client.Stkpush("254712345678", 10, "INV-001", "")
	if err != nil {
		fmt.Println("STK Push error:", err)
		return
	}
	fmt.Println("STK Push:", stkResult)

	// B2C – Business to Customer
	b2cResult, err := client.B2c("254712345678", "BusinessPayment", 500, "Salary payment")
	if err != nil {
		fmt.Println("B2C error:", err)
		return
	}
	fmt.Println("B2C:", b2cResult)

	// B2B – Business to Business
	b2bResult, err := client.B2b("600000", "BusinessPayBill", 1000, "Invoice payment", "INV-001")
	if err != nil {
		fmt.Println("B2B error:", err)
		return
	}
	fmt.Println("B2B:", b2bResult)

	// C2B – Register URLs
	registerResult, err := client.C2bregisterURLS(env("MPESA_SHORTCODE", "174379"), "", "")
	if err != nil {
		fmt.Println("C2B Register error:", err)
		return
	}
	fmt.Println("C2B Register:", registerResult)

	// C2B – Simulate payment
	simulateResult, err := client.C2bsimulate("254712345678", 100, env("MPESA_SHORTCODE", "174379"), mpesa.Paybill, "")
	if err != nil {
		fmt.Println("C2B Simulate error:", err)
		return
	}
	fmt.Println("C2B Simulate:", simulateResult)

	// Account Balance
	balanceResult, err := client.AccountBalance(env("MPESA_SHORTCODE", "174379"), 4, "Daily balance check")
	if err != nil {
		fmt.Println("Account Balance error:", err)
		return
	}
	fmt.Println("Account Balance:", balanceResult)

	// Transaction Status
	txResult, err := client.TransactionStatus(env("MPESA_SHORTCODE", "174379"), "OER7Q9I2PC", 1, "Transaction check")
	if err != nil {
		fmt.Println("Transaction Status error:", err)
		return
	}
	fmt.Println("Transaction Status:", txResult)

	// Reversal
	revResult, err := client.Reversal(env("MPESA_SHORTCODE", "174379"), "OER7Q9I2PC", 10, "Customer refund")
	if err != nil {
		fmt.Println("Reversal error:", err)
		return
	}
	fmt.Println("Reversal:", revResult)

	// B2 Pochi
	pochiResult, err := client.B2pochi("254712345678", 200, "Pochi payment")
	if err != nil {
		fmt.Println("B2 Pochi error:", err)
		return
	}
	fmt.Println("B2 Pochi:", pochiResult)
}
