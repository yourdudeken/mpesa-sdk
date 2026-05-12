# Getting Started with Safaricom Daraja APIs

## Introduction
Welcome to Safaricom APIs. This guide helps set up and integrate Daraja APIs into your application to start accepting M-Pesa payments.

The APIs are RESTful, using HTTP verbs (GET and POST) with JSON-encoded request parameters and responses. Code samples are available in Curl, Ruby, PHP, Python, NodeJS, and Java.

## Terminologies

| Term | Description |
|------|-------------|
| M-PESA | Mobile money transfer service in Kenya for storing and transferring money via mobile phones |
| Command IDs | Unique commands specifying transaction types |
| Salary Payment | Sends money to both registered and unregistered M-Pesa customers |
| Business Payment | Sends money to registered M-Pesa customers only |
| Promotion Payment | Sends promotional payments to registered M-Pesa customers with congratulatory messages |
| Short Code | Unique number for receiving customer payments (Pay Bill, Buy Goods, or Till Number) |
| Pay Bill | Collects money regularly from customers |
| Buy Goods | Used for retail purchases of goods and services |
| Till Number | Attached to a store for customer payments |
| REST | A software architectural style for creating web services |
| APIs | Functions and procedures for creating applications that access features or data of a service |
| DOCS | Technical documentation for using and integrating with an API |

## Testing on Localhost
Use an HTTP tunneling client like **Ngrok** or **LocalTunnel** to make local services accessible over the internet.

## Going Live
After testing in the sandbox, ensure:
- You have an M-PESA account (PayBill, Till Number, or B2C)
- You have access to the M-PESA Portal and have created an Admin or Business Manager
- Follow the steps on the "Go Live" page
- Contact `m-pesabusiness@safaricom.co.ke` if you lack an admin

## Callback and IP Whitelisting
Whitelist these IPs to accept callbacks from the Safaricom API Gateway:

```
196.201.214.200
196.201.214.206
196.201.213.114
196.201.214.207
196.201.214.208
196.201.213.44
196.201.212.127
196.201.212.138
196.201.212.129
196.201.212.136
196.201.212.74
196.201.212.69
```

## M-Pesa API Certificates
Download certificates for encrypting security credentials:
- **Sandbox**: M-Pesa public key certificate
- **Production**: M-Pesa public key certificate

Used for APIs: B2C, B2B, Transaction Status Query, Reversal.

## How to Generate Security Credentials
1. Write the unencrypted password into a byte array
2. Encrypt the array using the M-Pesa public key certificate:
   - Use the RSA algorithm
   - Apply PKCS #1.5 padding (not OAEP)
3. Convert the resulting encrypted byte array into a string using base64 encoding
4. The resulting base64-encoded string is the security credential

## Creating an HTTP Server
M-Pesa APIs are asynchronous. Responses are sent to the CallBackURL or ResultURL. Deploy an HTTP listener with POST methods to receive responses.

## SSL Certificates
Safaricom offers SSL Certificates for robust server authentication, data encryption, and protection against cyber threats.
