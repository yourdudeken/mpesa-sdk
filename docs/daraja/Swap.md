# Swap API

Swap between M-Pesa and float/working capital within business accounts.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/swap/v1/transfer`

## Overview
The Swap API enables businesses to transfer funds between their M-Pesa accounts and float/working capital accounts. This facilitates capital management and liquidity optimization for multi-account business operations.

### Key Features
- Inter-account fund transfers
- Working capital management
- Float account operations
- Real-time transfer processing
- Balance optimization
- Automated reconciliation

## How It Works
1. Business initiates swap transfer between accounts
2. Funds requested from source account
3. Target account credited with funds
4. Transfer processed and confirmed
5. Both accounts updated immediately
6. Transaction logged for audit trail

## Use Cases
- Working capital redistribution
- Float account management
- Balance optimization
- Multi-location fund management
- Liquidity management
- Operating account balancing

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Multiple business accounts
- Float/working capital account
- Business Admin/Manager operators

### Good to Know
Swap is internal transfer between business accounts without going through M-Pesa customers.

## Request Body
```json
{
  "InitiatorName": "swapuser",
  "SecurityCredential": "SAFVNChNHfVtXEZMBuVo+a1Hwr+DtrUVN3zVg==",
  "CommandID": "SwapFunds",
  "Amount": "50000",
  "SourceAccount": "600000",
  "TargetAccount": "600001",
  "Remarks": "Working capital swap"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| InitiatorName | Initiator username | String | swapuser |
| SecurityCredential | Encrypted password | String | EToK4lNR... |
| CommandID | Command type (SwapFunds) | String | SwapFunds |
| Amount | Transfer amount in KES | Numeric | 50000 |
| SourceAccount | Account to debit | String | 600000 |
| TargetAccount | Account to credit | String | 600001 |
| Remarks | Transfer notes | String | Working capital swap |

## Response Body
```json
{
  "ResponseCode": "0",
  "ResponseDescription": "Swap transfer completed successfully",
  "TransactionID": "OA90000000",
  "Amount": "50000",
  "SourceAccount": "600000",
  "TargetAccount": "600001",
  "Timestamp": "2024-01-15 14:30:00",
  "NewSourceBalance": "200000",
  "NewTargetBalance": "300000"
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ResponseCode | Transfer result | String | 0 |
| ResponseDescription | Status message | String | Swap transfer completed successfully |
| TransactionID | Unique transfer ID | String | OA90000000 |
| Amount | Amount transferred | Numeric | 50000 |
| SourceAccount | Source account | String | 600000 |
| TargetAccount | Target account | String | 600001 |
| Timestamp | Transfer date/time | DateTime | 2024-01-15 14:30:00 |
| NewSourceBalance | New source balance | Numeric | 200000 |
| NewTargetBalance | New target balance | Numeric | 300000 |

## Transfer Types

| Type | Description |
|------|-------------|
| M2F | M-Pesa to Float |
| F2M | Float to M-Pesa |
| M2M | M-Pesa to M-Pesa (different accounts) |
| WorkingAccount | Transfer to working account |

## Error Codes

| errorCode | errorMessage | Mitigation |
|-----------|-------------|------------|
| 400.002.02 | Invalid account | Verify account numbers |
| 401.002.01 | Invalid credentials | Regenerate token |
| 404.002.01 | Account not found | Check account exists |
| 500.001.1001 | Insufficient funds | Check source balance |
| 500.003.02 | System is busy | Retry request |

## Account Balance Limits

| Limit | Description |
|-------|-------------|
| Minimum transfer | KES 1 |
| Maximum transfer | KES 1,000,000 |
| Daily limit | Unlimited |

## Reconciliation
Swap transfers tracked through:
- Transaction logs
- Account statements
- Reconciliation reports
- Audit trails

## Best Practices
- Plan swaps during off-peak hours
- Maintain minimum balance requirements
- Track all swaps in accounting system
- Regular balance verification
- Monthly reconciliation

## Testing
Use Daraja Simulator with test accounts.

## Support
- **Chatbot:** Daraja Chatbot
- **Email:** apisupport@safaricom.co.ke
