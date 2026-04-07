# Mpesa SDK for C#

[![NuGet](https://img.shields.io/nuget/v/Yourdudeken.Mpesa.svg)](https://www.nuget.org/packages/Yourdudeken.Mpesa/)
[![License](https://img.shields.io/github/license/yourdudeken/mpesa.svg)](LICENSE.md)

A C# SDK for the Mpesa Daraja APIs. This SDK allows you to integrate Mpesa Daraja APIs into your .NET applications with ease.

## Installation

Install via NuGet Package Manager:

```bash
Install-Package Yourdudeken.Mpesa
```

Or via dotnet CLI:

```bash
dotnet add package Yourdudeken.Mpesa
```

## Usage

```csharp
using Yourdudeken.Mpesa;

var config = new MpesaConfig
{
    Environment = "sandbox",
    MpesaConsumerKey = "your_consumer_key",
    MpesaConsumerSecret = "your_consumer_secret",
    Passkey = "your_passkey",
    Shortcode = "174379",
    InitiatorName = "testapi",
    InitiatorPassword = "your_password",
    Callbacks = new Dictionary<string, string>
    {
        { "callback_url", "https://your-callback-url.com/callback" }
    }
};

var mpesa = new Mpesa(config);

var response = await mpesa.Stkpush(new Dictionary<string, object>
{
    { "phonenumber", "254712345678" },
    { "amount", 10 },
    { "accountNumber", "TEST001" }
});

Console.WriteLine(response);
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
| Environment | string | Yes | "sandbox" or "production" |
| MpesaConsumerKey | string | Yes | C2B Consumer Key |
| MpesaConsumerSecret | string | Yes | C2B Consumer Secret |
| B2cConsumerKey | string | No | B2C Consumer Key |
| B2cConsumerSecret | string | No | B2C Consumer Secret |
| Passkey | string | Yes | Lipa na Mpesa Online Passkey |
| Shortcode | string | Yes | Business Shortcode |
| TillNumber | string | No | Till Number |
| InitiatorName | string | Yes | Mpesa Initiator Name |
| InitiatorPassword | string | Yes | Mpesa Initiator Password |
| B2cShortcode | string | No | B2C Shortcode |
| Callbacks | Dictionary | No | Callback URLs |

## License

MIT License - see [LICENSE.md](LICENSE.md) for details.