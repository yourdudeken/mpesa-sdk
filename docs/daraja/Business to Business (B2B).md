# Business to Business (B2B) API

This API enables you to pay for goods and services directly from your business account to a till number, merchant store number or Merchant HO.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/b2b/v1/paymentrequest`

## Overview
This API enables you to pay for goods and services directly from your business account to a till number, merchant store number or Merchant HO. You can also use this API to pay a merchant on behalf of a consumer/requestor.

The transaction moves money from your MMF/Working account to the recipient's merchant account.

## Request Body
```json
{
  "Initiator": "API_Usename",
  "SecurityCredential": "FKXl/KPzT8hFOnozI+unz7mXDgTRbrlrZ+C1Vblxpbz7jliLAFa0E/...../uO4gzUkABQuCxAeq+0Hd0A==",
  "Command ID": "BusinessBuyGoods",
  "SenderIdentifierType": "4",
  "RecieverIdentifierType": "4",
  "Amount": "239",
  "PartyA": "123456",
  "PartyB": "000000",
  "AccountReference": "353353",
  "Requester": "254700000000",
  "Remarks": "OK",
  "QueueTimeOutURL": "https://mydomain.com/b2b/businessbuygoods/queue/",
  "ResultURL": "https://mydomain.com/b2b/businessbuygoods/result/"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| CommandID | Use `BusinessPayBill` only | String | BusinessPayBill |
| Initiator | M-Pesa API operator username (needs Org Business Pay Bill API initiator role) | String | Username |
| SecurityCredential | Encrypted password of the M-Pesa API operator | String | 32SzVdmCvjpmQfw3... |
| PartyA | Your shortcode (money deducted from here) | Number | 123454 |
| SenderIdentifierType | Type of shortcode deducting money (4 only) | Number | 4 |
| PartyB | The shortcode to which money will be moved | Number | 000000 |
| RecieverIdentifierType | Type of shortcode credited (4 only) | Number | 4 |
| Requester | Optional consumer's mobile number on whose behalf you're paying | Mobile | 254700000000 |
| Amount | The transaction amount | Number | 300 |
| AccountReference | Account number to associate with payment (up to 13 chars) | String | ACC#03929/4yu |
| Remarks | Additional information (up to 100 chars) | String | Any string |
| QueueTimeOutURL | URL for timeout notification | URL | https://ip:port/path |
| ResultURL | URL for transaction results | URL | https://ip:port/path |
| Occassion | Additional information (up to 100 chars) | String | Any string |

## Response Body
```json
{
  "OriginatorConversationID": "5118-111210482-1",
  "ConversationID": "AG_20230420_2010759fd5662ef6d054",
  "ResponseCode": "0",
  "ResponseDescription": "Accept the service request successfully."
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| OriginatorConversationID | Unique request identifier by Daraja | String | 5118-111210482-1 |
| ConversationID | Unique request identifier by M-Pesa | String | AG_20230420_2010759fd5662ef6d054 |
| ResponseCode | Status code (0 = successful submission) | String | 0 |
| ResponseDescription | Descriptive message of submission status | String | Accept the service request successfully. |

## Successful Result Body
```json
{
  "Result": {
    "ResultType": "0",
    "ResultCode": "0",
    "ResultDesc": "The service request is processed successfully",
    "OriginatorConversationID": "626f6ddf-ab37-4650-b882-b1de92ec9aa4",
    "ConversationID": "12345677dfdf89099B3",
    "TransactionID": "QKA81LK5CY",
    "ResultParameters": {
      "ResultParameter": [
        { "Key": "DebitAccountBalance", "Value": "{Amount={CurrencyCode=KES, MinimumAmount=618683, BasicAmount=6186.83}}" },
        { "Key": "Amount", "Value": "190.00" },
        { "Key": "DebitPartyAffectedAccountBalance", "Value": "Working Account|KES|346568.83|6186.83|340382.00|0.00" },
        { "Key": "TransCompletedTime", "Value": "20221110110717" },
        { "Key": "DebitPartyCharges", "Value": "" },
        { "Key": "ReceiverPartyPublicName", "Value": "000000– Biller Companty" },
        { "Key": "Currency", "Value": "KES" },
        { "Key": "InitiatorAccountCurrentBalance", "Value": "{Amount={CurrencyCode=KES, MinimumAmount=618683, BasicAmount=6186.83}}" }
      ]
    },
    "ReferenceData": {
      "ReferenceItem": [
        { "Key": "BillReferenceNumber", "Value": "19008" },
        { "Key": "QueueTimeoutURL", "Value": "https://mydomain.com/b2b/businessbuygoods/queue/" }
      ]
    }
  }
}
```

## Successful Result Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| Result | Root parameter enclosing the result message | JSON Object | {"Result":{}} |
| ResultType | Status code (0 = sent to listener) | Number | 0 |
| ResultCode | Transaction result code (0 = success) | Number | 0 |
| ResultDesc | Descriptive message for the result | String | The service request is processed successfully. |
| OriginatorConversationId | Unique request identifier by API gateway | String | AG_2376487236_126732989KJHJKH |
| ConversationId | Unique request identifier by M-Pesa | String | 236543-276372-2 |
| TransactionID | Unique M-PESA transaction ID | String | LHG31AA5TX |
| ResultParameters | JSON object with more transaction details | JSON Object | |
| Amount | Transaction amount | Number | 100 |
| TransactionCompletedTime | 14-digit completion timestamp | Number | 20171206163233 |
| ReceiverPartyPublicName | Public name of credit party | String | 600000 - saf test org |
| DebitPartyCharges | Transaction fee deducted | Number | 1 |
| Currency | Currency code | String | KES |
| DebitPartyAffectedAccountBalance | Organization's account balance after debit | String | Working Account\|KES\|500000.00\|... |
| DebitAccountCurrentBalance | Organization's account balance | String | {Amount={...}} |
| InitiatorAccountCurrentBalance | Organization's account balance | String | {Amount={...}} |
| ReferenceData | JSON object with reference data | JSON Object | |

## Unsuccessful Result Body
```json
{
  "Result": {
    "ResultType": 0,
    "ResultCode": 2001,
    "ResultDesc": "The initiator information is invalid.",
    "OriginatorConversationID": "12337-23509183-5",
    "ConversationID": "AG_20200120_0000657265d5fa9ae5c0",
    "TransactionID": "OAK0000000",
    "ResultParameters": {
      "ResultParameter": { "Key": "BOCompletedTime", "Value": 20200120164825 }
    },
    "ReferenceData": {
      "ReferenceItem": { "Key": "QueueTimeoutURL", "Value": "https://mydomain.com/b2b/businessbuygoods/queue/" }
    }
  }
}
```

## Failed Result Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| Result | Root parameter enclosing result | JSON Object | {"Result":{}} |
| ConversationId | Global unique identifier from M-Pesa | String | 236543-276372-2 |
| OriginatorConversationId | Global unique identifier from API proxy | String | AG_2376487236_126732989KJHJKH |
| ResultDesc | Status message from API | String | The initiator information is invalid. |
| ResultType | Status code (usually 0) | Number | 0 |
| ResultCode | Numeric status code (0 = success) | Number | 2001 |
| ResultParameters | JSON object with transaction details | JSON Object | |
| TransactionID | Unique M-PESA transaction ID | String | OAK0000000 |
| ReferenceData | JSON object with reference data | JSON Object | |
