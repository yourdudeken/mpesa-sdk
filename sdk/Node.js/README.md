# Mpesa SDK for Node.js

[![npm](https://img.shields.io/npm/v/@yourdudeken/mpesa-sdk.svg)](https://www.npmjs.com/package/@yourdudeken/mpesa-sdk)
[![License](https://img.shields.io/github/license/yourdudeken/mpesa.svg)](LICENSE.md)

A Node.js SDK for the Mpesa Daraja APIs. This SDK allows you to integrate Mpesa Daraja APIs into your Node.js applications with ease.

## Installation

```bash
npm install @yourdudeken/mpesa-sdk
```

## Usage

```javascript
const Mpesa = require('@yourdudeken/mpesa-sdk').default;

const mpesa = new Mpesa({
  environment: 'sandbox',
  mpesaConsumerKey: 'your_consumer_key',
  mpesaConsumerSecret: 'your_consumer_secret',
  passkey: 'your_passkey',
  shortcode: '174379',
  initiatorName: 'testapi',
  initiatorPassword: 'your_password',
  callbacks: {
    callbackUrl: 'https://your-callback-url.com/callback',
  },
});

async function testStkPush() {
  try {
    const response = await mpesa.stkpush({
      phonenumber: '254712345678',
      amount: 10,
      accountNumber: 'TEST001',
    });
    console.log(response);
  } catch (error) {
    console.error(error);
  }
}

testStkPush();
```

## Supported APIs

- **STK Push** - Lipa na Mpesa Express Online
- **STK Query** - Check transaction status
- **B2C** - Business to Customer
- **B2B** - Business to Business
- **B2Pochi** - Business to Pochi La Biashara
- **C2B** - Customer to Business (Register URL & Simulate)
- **Transaction Status** - Check transaction status
- **Account Balance** - Query account balance
- **Reversal** - Reverse a transaction

## Configuration

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| environment | string | Yes | 'sandbox' or 'production' |
| mpesaConsumerKey | string | Yes | C2B Consumer Key |
| mpesaConsumerSecret | string | Yes | C2B Consumer Secret |
| b2cConsumerKey | string | No | B2C Consumer Key |
| b2cConsumerSecret | string | No | B2C Consumer Secret |
| passkey | string | Yes | Lipa na Mpesa Online Passkey |
| shortcode | string | Yes | Business Shortcode |
| tillNumber | string | No | Till Number |
| initiatorName | string | Yes | Mpesa Initiator Name |
| initiatorPassword | string | Yes | Mpesa Initiator Password |
| b2cShortcode | string | No | B2C Shortcode |
| callbacks | object | No | Callback URLs |

## License

MIT License - see [LICENSE.md](LICENSE.md) for details.