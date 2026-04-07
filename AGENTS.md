# AGENTS.md

## Repository Structure

Multi-language Mpesa Daraja SDK in `packages/{Language}/`:
- **PHP**: Laravel package - `packages/php/` (main focus, tested in CI)
- **Node.js**: TypeScript SDK - `packages/node/`
- **Python**: SDK - `packages/python/`
- **Java**: Stub - `packages/java/`
- **C#**: Stub - `packages/dotnet/`
- **Go**: Stub - `packages/go/`

## Package Names

| Language | Package Name |
|----------|--------------|
| Node.js | @yourdudeken/mpesa-sdk |
| Python | yourdudeken-mpesa-sdk |
| PHP | yourdudeken/mpesa-sdk |
| Java | com.yourdudeken.mpesa |
| C# | Yourdudeken.Mpesa |
| Go | github.com/yourdudeken/mpesa-sdk |

## Directory Structure

### PHP (`packages/php/`)
```
php/
├── src/
│   ├── Config/MpesaConfig.php
│   ├── Services/
│   │   ├── STKPushService.php
│   │   ├── B2CService.php
│   │   ├── B2BService.php
│   │   ├── C2BService.php
│   │   ├── AccountService.php
│   │   └── B2PochiService.php
│   ├── Http/Client.php
│   ├── Helpers/Signature.php
│   ├── Exceptions/MpesaException.php
│   ├── Facades/Mpesa.php
│   ├── Console/InstallMpesaPackage.php
│   └── certificates/
├── config/mpesa.php
├── tests/
├── composer.json
└── phpunit.xml
```

### Node.js (`packages/node/`)
```
node/
├── src/
│   ├── core/
│   │   ├── MpesaClient.ts
│   │   └── Config.ts
│   ├── services/
│   ├── http/
│   ├── types/
│   ├── utils/
│   ├── errors/
│   └── certificates/
├── tests/
├── package.json
└── tsconfig.json
```

### Python (`packages/python/`)
```
python/
├── src/yourdudeken_mpesa_sdk/
│   ├── client.py
│   ├── config.py
│   ├── services/
│   ├── http/
│   ├── certificates/
│   └── __init__.py
├── tests/
├── pyproject.toml
└── setup.py
```

### Java (`packages/java/`)
```
java/
├── src/main/java/com/yourdudeken/mpesa/
│   ├── MpesaClient.java
│   ├── config/
│   ├── auth/
│   ├── services/
│   ├── http/
│   ├── exceptions/
│   └── models/
├── src/test/
├── pom.xml
└── src/main/resources/
```

### C# (`packages/dotnet/`)
```
dotnet/
├── src/Yourdudeken.Mpesa/
│   ├── MpesaClient.cs
│   ├── Config/
│   ├── Auth/
│   ├── Services/
│   ├── Http/
│   └── Exceptions/
├── tests/
└── MpesaSdk.csproj
```

### Go (`packages/go/`)
```
go/
├── mpesa/
│   ├── client.go
│   ├── config.go
│   ├── services/
│   └── http/
├── go.mod
└── certificates/
```

## Commands

### PHP (Laravel)
```bash
cd packages/php
composer install
composer test          # runs Pest tests
```

### Node.js
```bash
cd packages/node
npm install
npm run build          # TypeScript compilation
npm test               # Jest tests
```

### Python
```bash
cd packages/python
pip install -e .
```

## Testing

- **PHP**: Uses Pest (PHPUnit wrapper). Test command: `vendor/bin/pest`
- **CI**: Runs on PRs to `main` with PHP 8.2-8.4 and Laravel 10-12 matrix

## Key Files

- `packages/php/src/Mpesa.php` - Main PHP SDK class
- `packages/php/config/mpesa.php` - Laravel config (publish via `php artisan mpesa:install`)
- `packages/node/src/core/MpesaClient.ts` - Node.js entrypoint
- `packages/python/src/yourdudeken_mpesa_sdk/client.py` - Python SDK

## Environment

PHP requires `.env` with `MPESA_CONSUMER_KEY`, `MPESA_CONSUMER_SECRET`, etc.