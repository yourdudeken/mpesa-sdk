# M-Pesa Express Query API

Check the status of a Lipa Na M-Pesa Online Payment.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/stkpushquery/v1/query`

## Overview
Use this API to check the status of a Lipa Na M-Pesa Online Payment.

## Request Body
```json
{
  "BusinessShortCode": "174379",
  "Password": "MTc0Mzc5YmZiMjc5TliZGJjZjE1OGU5N2RkNzFhNDY3Y2QyZTBjODkzMDU5YjEwZjc4ZTZiNzJhZGExZWQyYzkxOTIwMTYwMjE2MTY1NjI3",
  "Timestamp": "20160216165627",
  "CheckoutRequestID": "ws_CO_260520211133524545"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| BusinessShortCode | Organization's shortcode (Paybill or Buygoods - 5 to 7 digits) | Numeric | 654321 |
| Password | Base64 encoded string of Shortcode+Passkey+Timestamp | String | base64.encode(Shortcode+Passkey+Timestamp) |
| Timestamp | Timestamp in format YYYYMMDDHHmmss | Timestamp | YYYYMMDDHHmmss |
| CheckoutRequestID | Global unique identifier of the processed checkout transaction | String | ws_CO_DMZ_123212312_2342347678234 |

## Response Body
```json
{
  "ResponseCode": "0",
  "ResponseDescription": "The service request has been accepted successfully",
  "MerchantRequestID": "22205-34066-1",
  "CheckoutRequestID": "ws_CO_13012021093521236557",
  "ResultCode": "0",
  "ResultDesc": "The service request is processed successfully."
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| MerchantRequestID | Global unique identifier for submitted payment request | String | 16813-1590513-1 |
| CheckoutRequestID | Global unique identifier of processed checkout transaction | String | ws_CO_DMZ_123212312_2342347678234 |
| ResponseCode | Status code (0 = successful submission) | Numeric | 0 |
| ResultDesc | Status message from API, maps to ResultCode value | String | 0: The service request is processed successfully. 1032: Request canceled by the user |
| ResponseDescription | Acknowledgment message from API | String | The service request has been accepted successfully |
| ResultCode | Status code (0 = successful processing) | Numeric | 0, 1032 |

## Error Response Parameter Definition

| Error | Description | Mitigation |
|-------|-------------|------------|
| 404.001.04 | Invalid Authentication Header | All M-PESA APIs use POST except Authorization (GET) |
| 400.002.05 | Invalid Request Payload | Submit correct request payload as per documentation |
| 400.003.01 | Invalid Access Token | Regenerate token before expiry |
