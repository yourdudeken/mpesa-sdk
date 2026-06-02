# Dynamic QR API

Generates a dynamic M-PESA QR Code for customer payments.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/qrcode/v1/generate`

## Overview
The Dynamic QR API enables generation of M-Pesa QR codes that customers can scan using the Safaricom My App or M-Pesa App. Customers scan the QR code, capture till number and amount, then authorize payment for goods and services at LIPA NA M-PESA (LNM) merchant outlets.

### Key Features
- Dynamic QR code generation
- Real-time customizable QR codes
- Base64 encoded QR image response
- Support for multiple transaction types
- Instant payment initiation via scan

## How It Works
1. Merchant generates dynamic QR code with transaction details
2. QR code displayed at point of sale or sent to customer
3. Customer scans QR with M-Pesa or Safaricom app
4. Customer reviews transaction details (till number, amount)
5. Customer enters M-Pesa PIN to authorize
6. Payment processed and confirmed
7. Merchant receives transaction notification

## Use Cases
- Retail store POS payments
- Restaurant bill payments
- E-commerce checkout
- Event/ticket sales
- Utility bill payments
- Marketplace transactions
- Contactless payment promotion

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Consumer Key & Consumer Secret
- Merchant till number or paybill code
- Business registered for M-Pesa

## Request Body
```json
{
  "MerchantName": "TEST SUPERMARKET",
  "RefNo": "Invoice-12345",
  "Amount": 2000,
  "TrxCode": "BG",
  "CPI": "373132",
  "Size": "300"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample | Required |
|------|-------------|------|--------|----------|
| MerchantName | Name of the Company/M-Pesa Merchant | String | "TEST-Supermarket" | Yes |
| RefNo | Transaction Reference/Invoice Number | String | "INV-2024-001" | Yes |
| Amount | Total amount for the transaction in KES | Numeric | 2000 | Yes |
| TrxCode | Transaction Type Code | String | BG, WA, PB, SM, SB | Yes |
| CPI | Credit Party Identifier (phone, business number, till, or paybill) | String | "254722000000" or "373132" | Yes |
| Size | Size of QR code image in pixels (square, 100-500) | String | "300" | No |

### TrxCode Values

| Code | Description | CPI Format | Use Case |
|------|-------------|-----------|----------|
| BG | Pay Merchant (Buy Goods) | Till number or Merchant number | Retail purchases |
| WA | Withdraw Cash at Agent Till | Agent till number | ATM-style withdrawals |
| PB | Paybill or Business Number | Business shortcode (4-8 digits) | Bill payments |
| SM | Send Money (Mobile number) | MSISDN format (2547XXXXXXXX) | Person-to-person transfers |
| SB | Send to Business (Business number) | Business number in MSISDN format | Business transfers |

## Response Body
```json
{
  "ResponseCode": "AG_20191219_000043fdf61864fe9ff5",
  "RequestID": "16738-27456357-1",
  "ResponseDescription": "QR Code Successfully Generated.",
  "QRCode": "iVBORw0KGgoAAAANSUhEUgAAASwAAAEsCAIAAAD2HxkiAAAHtEl..."
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ResponseCode | Unique response identifier from M-Pesa | String | AG_20191219_000043fdf61864fe9ff5 |
| RequestID | Unique request identifier from Daraja | String | 16738-27456357-1 |
| ResponseDescription | Status of the transaction | String | QR Code Successfully Generated |
| QRCode | QR Code Image as Base64 encoded string | String | iVBORw0KGgoAAAANSUhEUgAAASwAAAEsCAIAAAD... |

## QR Code Display

The QR code returned is Base64 encoded. To display:

### HTML/Web
```html
<img src="data:image/png;base64,{QRCode}" alt="M-Pesa QR Code" />
```

### React
```jsx
<img src={`data:image/png;base64,${qrCode}`} alt="M-Pesa QR" />
```

### Mobile Apps
Decode Base64 and display as PNG image

## Error Codes

| Error | Description | HTTP | Mitigation |
|-------|-------------|------|-----------|
| 404.001.04 | Invalid Authentication Header | 401 | Use POST method; ensure Authorization header |
| 400.002.05 | Invalid Request Payload | 400 | Submit correct payload as per documentation |
| 400.003.01 | Invalid Access Token | 401 | Regenerate token; ensure not expired |
| 400.002.02 | Invalid CPI format | 400 | Verify CPI matches transaction type |
| 400.002.03 | Invalid Amount | 400 | Ensure amount is numeric and positive |

## Transaction Limits

| Limit | Value |
|-------|-------|
| Minimum Amount | KES 1 |
| Maximum Amount | KES 250,000 |
| QR Code Size | 100-500 pixels |

## Best Practices
- Generate new QR code for each transaction
- Display clear instructions for customers to scan
- Include merchant name and amount for customer clarity
- Use appropriate QR size (minimum 200px for scanning reliability)
- Store QR codes temporarily (valid for one transaction)
- Validate transaction amount before QR generation

## Testing
1. Create sandbox app on Daraja Portal
2. Generate access token via Authorization API
3. Call Dynamic QR API with test data
4. Use simulator to scan and test payment

## Production Deployment
Same endpoints work for production with live merchant credentials.

## Support
- **Chatbot:** Daraja Chatbot
- **Email:** apisupport@safaricom.co.ke
