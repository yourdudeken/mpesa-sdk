# Mpesa SDK

Multi-language SDK for Mpesa Daraja API.

[![PHP](https://img.shields.io/badge/PHP-8.2+-777bb3?style=flat-square&logo=php)](https://packagist.org/packages/yourdudeken/mpesa)
[![Node.js](https://img.shields.io/badge/Node.js-16+-339933?style=flat-square&logo=node.js)](https://www.npmjs.com/package/@yourdudeken/mpesa-sdk)
[![Python](https://img.shields.io/badge/Python-3.8+-3776AB?style=flat-square&logo=python)](https://pypi.org/project/yourdudeken-mpesa-sdk/)
[![Java](https://img.shields.io/badge/Java-11+-b07219?style=flat-square&logo=java)](https://search.maven.org/artifact/com.yourdudeken.mpesa)
[![C#](https://img.shields.io/badge/C%23-.NET%206+-239120?style=flat-square&logo=csharp)](https://www.nuget.org/packages/Yourdudeken.Mpesa)
[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat-square&logo=go)](https://pkg.go.dev/github.com/yourdudeken/mpesa-sdk)

## Supported Languages

| Language | Package | Version |
|----------|---------|---------|
| PHP/Laravel | `yourdudeken/mpesa-sdk` | v1.0.0 |
| Node.js/TypeScript | `@yourdudeken/mpesa-sdk` | v1.0.0 |
| Python | `yourdudeken-mpesa-sdk` | v1.0.0 |
| Java | `com.yourdudeken.mpesa` | v1.0.0 |
| C#/.NET | `Yourdudeken.Mpesa` | v1.0.0 |
| Go | `github.com/yourdudeken/mpesa-sdk` | v1.0.0 |

## Quick Start

### PHP (Laravel)
```bash
composer require yourdudeken/mpesa-sdk
php artisan mpesa:install
```

### Node.js
```bash
npm install @yourdudeken/mpesa-sdk
```

### Python
```bash
pip install yourdudeken-mpesa-sdk
```

## Documentation

- [PHP SDK](./packages/php/README.md)
- [Node.js SDK](./packages/node/README.md)
- [Python SDK](./packages/python/README.md)
- [Java SDK](./packages/java/README.md)
- [C# SDK](./packages/dotnet/README.md)
- [Go SDK](./packages/go/README.md)

## API Reference

All SDKs expose a consistent interface:

```php
// PHP
$mpesa->stkpush($phone, $amount, $account);
$mpesa->b2c($phone, $command, $amount, $remarks);
$mpesa->c2bregisterURLS($shortcode, $confirmUrl, $validateUrl);
$mpesa->transactionStatus($shortcode, $transactionId, $identifierType, $remarks);
```

```typescript
// Node.js
await mpesa.stkpush({ phonenumber, amount, accountNumber });
await mpesa.b2c({ phonenumber, commandId, amount, remarks });
await mpesa.c2bregisterURLS({ shortcode, confirmUrl, validateUrl });
await mpesa.transactionStatus({ shortcode, transactionId, identifierType, remarks });
```

```python
# Python
mpesa.stkpush(phonenumber, amount, account_number)
mpesa.b2c(phonenumber, command_id, amount, remarks)
mpesa.c2bregisterURLS(shortcode, confirm_url, validate_url)
mpesa.transaction_status(shortcode, transaction_id, identifier_type, remarks)
```

## Supported APIs

- **STK Push** - Lipa na Mpesa Online
- **B2C** - Business to Customer
- **B2B** - Business to Business
- **C2B** - Customer to Business
- **B2Pochi** - Business to Pochi
- **Account Balance** - Check till balance
- **Transaction Status** - Query transaction
- **Reversal** - Reverse transaction

## License

MIT License - see [LICENSE.md](./LICENSE.md)
