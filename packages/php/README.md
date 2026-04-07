# Mpesa PHP SDK

A Laravel package for Mpesa Daraja API.

## Installation

```bash
composer require yourdudeken/mpesa-sdk
php artisan mpesa:install
```

## Configuration

```php
// config/mpesa.php
return [
    'environment' => env('MPESA_ENVIRONMENT', 'sandbox'),
    'mpesa_consumer_key' => env('MPESA_CONSUMER_KEY'),
    'mpesa_consumer_secret' => env('MPESA_CONSUMER_SECRET'),
    'passkey' => env('SAFARICOM_PASSKEY'),
    'shortcode' => env('MPESA_BUSINESS_SHORTCODE'),
    'initiator_name' => env('MPESA_INITIATOR_NAME'),
    'initiator_password' => env('MPESA_INITIATOR_PASSWORD'),
    // ... more config
];
```

## Usage

```php
use Yourdudeken\Mpesa\Facades\Mpesa;

// STK Push
$response = Mpesa::stkpush(
    phonenumber: '254712345678',
    amount: 100,
    account_number: 'ORDER123'
);

// B2C
$response = Mpesa::b2c(
    phonenumber: '254712345678',
    command_id: 'BusinessPayment',
    amount: 1000,
    remarks: 'Payment'
);
```

## API Reference

- `stkpush()` - Lipa na Mpesa Online
- `stkquery()` - Query STK Push status
- `b2c()` - Business to Customer
- `b2b()` - Business to Business
- `c2bregisterURLS()` - Register C2B URLs
- `c2bsimulate()` - Simulate C2B
- `transactionStatus()` - Query transaction
- `accountBalance()` - Check balance
- `reversal()` - Reverse transaction

## License

MIT License
