# IoT SIM Management API

Manage IoT SIM cards and connectivity for IoT devices.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/iot/v1/manage`

## Overview
The IoT SIM Management API enables businesses to manage M2M (Machine-to-Machine) SIM cards for IoT devices. This includes activation, deactivation, status checks, and data usage management.

### Key Features
- IoT SIM activation and deactivation
- Device connectivity management
- Real-time status monitoring
- Data usage tracking
- Automated billing integration
- Remote SIM management

## How It Works
1. Business manages IoT SIM through API
2. SIM activation/deactivation processed
3. Device connectivity established
4. Real-time monitoring enabled
5. Usage tracked and billed
6. Status reports provided

## Use Cases
- Smart meter management
- IoT device connectivity
- M2M communication
- Remote device monitoring
- Smart city infrastructure
- Industrial IoT deployments
- Asset tracking devices

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- IoT SIM card inventory
- Business account setup for IoT
- Device database configured

### Good to Know
IoT SIM cards require special enterprise setup and dedicated connectivity plans.

## Request Body
```json
{
  "InitiatorName": "iotuser",
  "SecurityCredential": "SAFVNChNHfVtXEZMBuVo+a1Hwr+DtrUVN3zVg==",
  "CommandID": "ActivateIOTSIM",
  "ICCID": "8934401122334455667788",
  "IMEI": "356938035643809",
  "DeviceName": "SmartMeter_001",
  "DeviceLocation": "Building A",
  "DataPlan": "10GB",
  "BillingCycle": "Monthly"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| InitiatorName | Initiator username | String | iotuser |
| SecurityCredential | Encrypted password | String | EToK4lNR... |
| CommandID | Operation (ActivateIOTSIM, DeactivateIOTSIM, CheckStatus) | String | ActivateIOTSIM |
| ICCID | Integrated Circuit Card Identifier | String | 8934401122334455667788 |
| IMEI | International Mobile Equipment Identity | String | 356938035643809 |
| DeviceName | Name of IoT device | String | SmartMeter_001 |
| DeviceLocation | Physical device location | String | Building A |
| DataPlan | Monthly data allocation | String | 10GB |
| BillingCycle | Billing frequency | String | Monthly |

## Response Body
```json
{
  "ResponseCode": "0",
  "ResponseDescription": "SIM activated successfully",
  "ICCID": "8934401122334455667788",
  "Status": "Active",
  "ActivationDate": "2024-01-15",
  "DataPlan": "10GB",
  "ExpiryDate": "2025-01-15"
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ResponseCode | Operation result | String | 0 |
| ResponseDescription | Status message | String | SIM activated successfully |
| ICCID | SIM card identifier | String | 8934401122334455667788 |
| Status | Current SIM status | String | Active |
| ActivationDate | Activation date | Date | 2024-01-15 |
| DataPlan | Active data plan | String | 10GB |
| ExpiryDate | Plan expiry date | Date | 2025-01-15 |

## SIM Status Codes

| Status | Description |
|--------|-------------|
| Active | SIM is operational |
| Inactive | SIM is not active |
| Suspended | SIM suspended |
| Expired | SIM plan expired |
| Blocked | SIM blocked |

## Operations Supported

| Operation | Description |
|-----------|-------------|
| ActivateIOTSIM | Activate new IoT SIM |
| DeactivateIOTSIM | Deactivate active SIM |
| CheckStatus | Check SIM status |
| UpdateDataPlan | Update data plan |
| ReportUsage | Get usage report |
| SuspendSIM | Suspend SIM temporarily |

## Error Codes

| errorCode | errorMessage | Mitigation |
|-----------|-------------|------------|
| 400.002.02 | Invalid ICCID format | Verify ICCID number |
| 401.002.01 | Invalid credentials | Regenerate token |
| 404.002.01 | SIM not found | Check ICCID |
| 500.001.1001 | SIM already active | Cannot reactivate |
| 500.003.02 | System is busy | Retry request |

## Data Usage Tracking
Monitor IoT device data consumption:
- Real-time usage reports
- Usage alerts and thresholds
- Auto-throttling options
- Plan upgrade capabilities

## Billing Integration
IoT billing through Daraja:
- Automated monthly billing
- Usage-based charges
- Plan-based subscriptions
- Invoice generation

## Testing
Use Daraja Simulator with test IoT SIM details.

## Go Live
Submit IoT SIM inventory and device infrastructure details.

## Support
- **Chatbot:** Daraja Chatbot
- **Email:** apisupport@safaricom.co.ke
