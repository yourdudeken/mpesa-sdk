# Mpesa C# SDK

A .NET SDK for Mpesa Daraja API.

## Installation

```bash
dotnet add package Yourdudeken.Mpesa
```

## Usage

```csharp
using Yourdudeken.Mpesa;

var config = new MpesaConfig
{
    Environment = "sandbox",
    MpesaConsumerKey = "your_key",
    MpesaConsumerSecret = "your_secret",
    Shortcode = "174379",
    Passkey = "your_passkey",
    InitiatorName = "testapi",
    InitiatorPassword = "your_password"
};

var mpesa = new MpesaClient(config);
var response = await mpesa.StkPush("254712345678", 100, "ORDER123");
```

## License

MIT License
