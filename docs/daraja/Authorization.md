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

**Method:** GET

**Authorization:** Basic Auth (Base64 encode Consumer Key:Consumer Secret)

**URL:** `https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials`

**Headers:**
```
Authorization: Basic {Base64(ConsumerKey:ConsumerSecret)}
```

**Example with curl:**
```bash
curl -X GET "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials" \
  -H "Authorization: Basic Q2RtTmJkdDBpQk4xb3FEZkthc200ZGFiZHBLbXRhTm46RExLRzdQQnVuNzIwR1ppbQ=="
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

## Production Endpoints
**Token Endpoint:** `https://api.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials`

Token expiry is 3600 seconds for both sandbox and production.

## Code Examples

### Python
```python
import requests
import base64

consumer_key = "YOUR_CONSUMER_KEY"
consumer_secret = "YOUR_CONSUMER_SECRET"
auth_url = "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"

credentials = f"{consumer_key}:{consumer_secret}"
encoded = base64.b64encode(credentials.encode()).decode()

headers = {"Authorization": f"Basic {encoded}"}
response = requests.get(auth_url, headers=headers)
token = response.json()["access_token"]
```

### JavaScript/Node.js
```javascript
const axios = require('axios');
const base64 = require('base-64');

const consumerKey = "YOUR_CONSUMER_KEY";
const consumerSecret = "YOUR_CONSUMER_SECRET";
const credentials = `${consumerKey}:${consumerSecret}`;
const encoded = base64.encode(credentials);

const config = {
  headers: {"Authorization": `Basic ${encoded}`}
};

axios.get("https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials", config)
  .then(response => console.log(response.data.access_token));
```

## Token Validation

Tokens are valid for 3600 seconds (1 hour) from generation. When a token expires:
- All API calls using the expired token will return 401 Unauthorized
- A new token must be generated using the Authorization API
- Implement token caching in your application to minimize API calls

## Security Best Practices
- Store Consumer Key and Consumer Secret securely (environment variables, vaults)
- Do not hardcode credentials in source code
- Rotate credentials periodically
- Use HTTPS for all token requests
- Implement token refresh logic before expiry (optional: refresh at 55 minutes)

## Support
### FAQs
- **Why is my access token not working?** Tokens expire in 3600 seconds. Generate a new one if expired.
- **What should I do if I get an invalid grant type error?** Ensure grant_type is set to `client_credentials`.
- **Can I generate multiple tokens?** Yes, but each request generates a unique token. Older tokens remain valid until expiry.
- **How do I handle token expiration?** Implement token refresh logic or cache tokens until ~1 minute before expiry.

### Chatbot
Developers can get instant responses using the Daraja Chatbot.

### Production Issues & Incident Management
- **Incident Management Page:** Visit the Incident Management page
- **Email:** apisupport@safaricom.co.ke
- **Response Time:** Support team responds within business hours
