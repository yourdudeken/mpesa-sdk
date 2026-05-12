# Authorization API

Gives a time-bound access token to call allowed APIs.

**Endpoint:** `GET https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials`

## Overview
The Authorization API generates access tokens required for authenticating all Daraja API calls using OAuth 2.0 Authentication.

### Key Features
- OAuth 2.0 Authentication
- Token Expiry: 3600 seconds (1 hour)
- Supports Automated Testing via the Simulator
- Postman Collection support

> **Note:** This API must be called before any other API in the Daraja platform, as all other APIs require an access token for authentication.

## How It Works
1. Developer retrieves Consumer Key and Consumer Secret from the Daraja Portal
2. Developer sends a request to the Authorization API using Basic Authentication
3. API validates credentials and returns an access token
4. The token is then used in subsequent API calls

## Getting Started
### Prerequisites
- Create a Daraja Account on Safaricom Developer Portal
- Create a sandbox app to get API credentials
- Retrieve Consumer Key & Consumer Secret from your sandbox app on My Apps

### Good to Know
- **Token Expiry:** Tokens expire after 3600 seconds (1 hour)

## Integration Steps
### Generate an OAuth Access Token

**Authorization:** Basic Auth (username/password)

**Headers:**
```json
{
  "Authorization": "Basic Q2RtTmJkdDBpQk4xb3FEZkthc200ZGFiZHBLbXRhTm46RExLRzdQQnVuNzIwR1ppbQ=="
}
```

**Params:**
```
grant_type: client_credentials
```

### Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| grant_type | The client_credentials grant type is supported | Query | client_credentials |

### Response Body
```json
{
  "access_token": "c9SQxWWhmdVRlyh0zh8gZDTkubVF",
  "expires_in": 3599
}
```

### Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| access_token | Access token to access the APIs | JSON Response Item | c9SQxWWhmdVRlyh0zh8gZDTkubVF |
| expires_in | Token expiry time in seconds | JSON Response Item | 3599 |

### Error Response Parameter Definition

| Error | Description | Probable Cause | Mitigation |
|-------|-------------|----------------|------------|
| 400.008.02 | Invalid grant type passed | Incorrect grant type | Select grant type as client_credentials |
| 400.008.01 | Invalid authentication type passed | Incorrect Authorisation type | Select authorization type as Basic |

## Testing
### Using the Simulator
Safaricom provides an Authorization API Simulator for testing access token generation via the Daraja Portal.

## Support
### FAQs
- **Why is my access token not working?** Tokens expire in 3600 seconds. Generate a new one if expired.
- **What should I do if I get an invalid grant type error?** Ensure grant_type is set to `client_credentials`.
- **Can I generate multiple tokens?** Yes, but each request invalidates the previous token.

### Chatbot
Developers can get instant responses using the Daraja Chatbot.

### Production Issues & Incident Management
- **Incident Management Page:** Visit the Incident Management page
- **Email:** apisupport@safaricom.co.ke
