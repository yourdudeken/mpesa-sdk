# AGENTS.md

## Repository Structure

Multi-language Mpesa Daraja SDK. Each SDK is in `sdk/{Language}/`:
- **PHP**: Laravel package (main focus) - `sdk/PHP/`
- **Node.js**: TypeScript SDK - `sdk/Node.js/`
- **Python**: SDK - `sdk/Python/`
- **Go, Java, C#**: Stub implementations

## Commands

### PHP (Laravel)
```bash
cd sdk/PHP
composer install
composer test          # runs Pest tests
```

### Node.js
```bash
cd sdk/Node.js
npm install
npm run build          # TypeScript compilation
npm test               # Jest tests
```

### Python
```bash
cd sdk/Python
pip install -e .
# No test script defined
```

## Testing

- **PHP**: Uses Pest (PHPUnit wrapper). Test command: `vendor/bin/pest`
- **CI**: Runs on PRs to `main` with PHP 8.2-8.4 and Laravel 10-12 matrix

## Key Files

- `config/mpesa.php` - PHP Laravel config (publish via `php artisan mpesa:install`)
- `src/Mpesa.php` - Main PHP SDK class
- `src/index.ts` - Node.js entrypoint
- `src/yourdudeken_mpesa_sdk/mpesa.py` - Python SDK

## Environment

PHP requires `.env` with `MPESA_CONSUMER_KEY`, `MPESA_CONSUMER_SECRET`, etc.
