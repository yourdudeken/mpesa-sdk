# IMSI API

Retrieve International Mobile Subscriber Identity information for mobile numbers.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/imsi/v1/query`

## Overview
The IMSI (International Mobile Subscriber Identity) API enables retrieval of mobile subscriber information. This is used for number validation and subscriber identity verification through Daraja.

### Key Features
- Mobile number validation
- Subscriber identity verification
- Real-time IMSI lookup
- Integration with payment systems
- Enhanced transaction security

## How It Works
1. Application queries IMSI API with phone number
2. API validates the mobile number format
3. Safaricom validates the subscriber status
4. IMSI information returned
5. Application receives subscriber details
6. Integration with transaction processing

## Use Cases
- Mobile number validation before payments
- Subscriber identity verification
- Account registration validation
- Fraud prevention
- Customer verification
- Mobile wallet integration

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Consumer Key & Consumer Secret
- Active API access to IMSI service

### Good to Know
IMSI is primarily for Safaricom Kenya mobile numbers in format 2547XXXXXXXX.

## Request Body
```json
{
  "PhoneNumber": "254722000000",
  "AccessToken": "YYhZ20EF2nlgD2ekqK1Sy70b3eY"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| PhoneNumber | Mobile number to verify (format: 2547XXXXXXXX) | String | 254722000000 |
| AccessToken | OAuth access token from Authorization API | String | YYhZ20EF2nlgD2ekqK1Sy70b3eY |

## Response Body
```json
{
  "ResponseCode": "0",
  "ResponseDescription": "Success",
  "PhoneNumber": "254722000000",
  "IMSI": "639001000000001",
  "SubscriberStatus": "Active",
  "NetworkOperator": "Safaricom"
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ResponseCode | Query result code | String | 0 |
| ResponseDescription | Status message | String | Success |
| PhoneNumber | Queried phone number | String | 254722000000 |
| IMSI | International Mobile Subscriber Identity | String | 639001000000001 |
| SubscriberStatus | Subscriber status (Active/Inactive) | String | Active |
| NetworkOperator | Mobile network operator | String | Safaricom |

## IMSI Information
The IMSI consists of:
- MCC (Mobile Country Code): 639 for Kenya
- MNC (Mobile Network Code): 01 for Safaricom
- MSIN (Mobile Subscription Identification Number): Unique subscriber identifier

## Error Codes

| errorCode | errorMessage | Mitigation |
|-----------|-------------|------------|
| 400.002.02 | Invalid Phone Number | Use format 2547XXXXXXXX |
| 401.002.01 | Invalid Access Token | Regenerate token |
| 404.002.01 | Number not found | Verify number exists |
| 500.003.02 | Service unavailable | Retry request |
| 500.001.1001 | Subscriber not found | Check number validity |

## Use Cases
- Pre-payment validation
- Customer KYC processes
- Fraud detection systems
- Account linking
- Payment gateway integration

## Testing
Test with Daraja provided test numbers.

## Support
- **Chatbot:** Daraja Chatbot
- **Email:** apisupport@safaricom.co.ke
