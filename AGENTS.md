# Agent Guide: Mpesa SDK

High-signal guidance for working in this multi-language SDK repository.

## Project Structure
Polyglot repo — SDKs in `packages/node/` (TypeScript), `packages/python/` (Setuptools), `packages/go/` (Go module). Each is standalone; no cross-package dependencies or workspace orchestration.

## Developer Commands

### Node.js (`packages/node/`)
- **Build**: `npm run build` (runs `tsc`).
- **Test**: `npm test` (runs `jest`).
- **Entrypoint**: `src/index.ts`.
- **Note**: `package.json` devDependencies include `jest`, `ts-jest`, and `@types/jest`.

### Python (`packages/python/`)
- **Install**: `pip install -e .` (uses `setup.py`).
- **Test**: `pytest`.
- **Entrypoint**: `src/mpesa/__init__.py` and `src/mpesa/client.py`.
- **Dependencies**: `requests` and `cryptography` (both in `setup.py` and `pyproject.toml`).

### Go (`packages/go/`)
- **Test**: `go test ./mpesa/...` (21 tests, all pass).
- **Module**: `github.com/yourdudeken/mpesa-sdk`
- **Constructor**: `NewClient(config *Config)`.
- **Methods**: positional params, e.g. `client.Stkpush(phonenumber, amount, accountNumber, callbackURL)`.

## Naming Conventions Across Languages

| API | Node (camelCase) | Python (snake_case) | Go (CamelCase) |
|-----|------------------|----------------------|----------------|
| STK Push | `stkpush` | `stkpush` | `Stkpush` |
| STK Query | `stkquery` | `stkquery` | `Stkquery` |
| B2C | `b2c` | `b2c` | `B2c` |
| Validated B2C | `validated_b2c` | `validated_b2c` | `Validated_b2c` |
| B2B | `b2b` | `b2b` | `B2b` |
| C2B Register | `c2bregisterURLS` | `c2b_register_urls` | `C2bregisterURLS` |
| C2B Simulate | `c2bsimulate` | `c2bsimulate` | `C2bsimulate` |
| Account Balance | `accountBalance` | `account_balance` | `AccountBalance` |
| Transaction Status | `transactionStatus` | `transaction_status` | `TransactionStatus` |
| Reversal | `reversal` | `reversal` | `Reversal` |
| B2 Pochi | `b2pochi` | `b2pochi` | `B2pochi` |

Python naming follows snake_case consistently (including `c2b_register_urls`).

## Key Patterns
- **Phone validation**: strip `+`, replace leading `0` with `254`, prepend `254` to numbers starting with `7`. Consistent across all SDKs.
- **Static constants**: `PAYBILL = 'CustomerPayBillOnline'`, `TILL = 'CustomerBuyGoodsOnline'`. Same value in all languages.
- **Callback config keys**: All snake_case (`callback_url`, `b2c_result_url`, `status_result_url`, etc.) — consistent across SDKs.
- **Token auth**: OAuth client_credentials grant. B2C/B2B operations use separate `b2cConsumerKey`/`b2cConsumerSecret` if configured.
- **Security credential**: RSA-encrypted initiator password using Mpesa certificate (PKCS1v15 padding). Certificate paths differ per package.

## Certificate Locations
- Node: `packages/node/src/certificates/`
- Python: `packages/python/src/mpesa/certificates/`
- Go: `packages/go/` (root, not in subdirectory)

## Notable Gaps
- No CI workflows exist (`.github/workflows/` is empty).
- `examplea/` directory exists but is empty.
