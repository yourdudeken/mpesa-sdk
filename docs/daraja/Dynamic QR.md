# Dynamic QR API

Generates a dynamic M-PESA QR Code.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/qrcode/v1/generate`

## Overview
Use this API to generate a Dynamic QR code that enables Safaricom M-PESA customers who have My Safaricom App or M-PESA app, to scan a QR code to capture till number and amount, then authorize to pay for goods and services at select LIPA NA M-PESA (LNM) merchant outlets.

## Request Body
```json
{
  "MerchantName": "TEST SUPERMARKET",
  "RefNo": "Invoice Test",
  "Amount": 1,
  "TrxCode": "BG",
  "CPI": "373132",
  "Size": "300"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| MerchantName | Name of the Company/M-Pesa Merchant Name | String | "TEST-Supermarket" |
| RefNo | Transaction Reference | String | "xewr34fer4t" |
| Amount | Total amount for the sale/transaction | Numeric | 2000 |
| TrxCode | Transaction Type | String | BG, WA, PB, SM, SB |
| CPI | Credit Party Identifier (Mobile Number, Business Number, Agent Till, Paybill, or Merchant Buy Goods) | String | "174379" |
| Size | Size of QR code image in pixels (square) | String | "300" |

### TrxCode Values

| Code | Description |
|------|-------------|
| BG | Pay Merchant (Buy Goods) |
| WA | Withdraw Cash at Agent Till |
| PB | Paybill or Business number |
| SM | Send Money (Mobile number) |
| SB | Sent to Business (Business number CPI in MSISDN format) |

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
| ResponseCode | Used to return the Transaction Type | String | Alpha-numeric string |
| ResponseDescription | Status of the transaction | String | QR Code Successfully Generated |
| QRCode | QR Code Image/Data/String (base64) | String | Alpha-numeric string containing the QR Code |

## Error Response Parameter Definition

| Error | Description | Mitigation |
|-------|-------------|------------|
| 404.001.04 | Invalid Authentication Header | All M-PESA APIs use POST except Authorization (GET) |
| 400.002.05 | Invalid Request Payload | Submit correct request payload as per documentation |
| 400.003.01 | Invalid Access Token | Regenerate token before expiry |
