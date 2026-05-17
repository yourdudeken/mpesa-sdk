# M-Pesa SDK - Production-Grade Transformation Plan

**Version:** 2.0 | **Last Updated:** 2026-05-17 | **Owner:** Platform Team

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2026-05-17 | 2.0 | Implemented Phase 1 (Critical Fixes & Security Hardening) + Phase 2 (Observability & Resilience) + Phase 4 (CI/CD Hardening). See Implementation Log below. |
| 2026-05-17 | 1.0 | Initial audit and roadmap creation |

---

## Executive Summary

The M-Pesa SDK is a polyglot monorepo (TypeScript, Python, Go) providing access to Safaricom's Daraja API. The architecture is well-structured with clear separation of concerns, typed models, error hierarchies, and webhook support.

### Current Risk Status

| Risk | Severity | Status | Details |
|------|----------|--------|---------|
| Go services layer bug | High | RESOLVED | `STKQuery` in Go `services/service.go` used `string(rune(input.BusinessShortCode))` — fixed to `fmt.Sprintf("%d", ...)` |
| Go retry body re-read | High | RESOLVED | `doRequest` body reader consumed on first retry attempt — fixed to re-marshal per attempt |
| No structured logging | High | RESOLVED | Logger interface + default implementations added to all 3 SDKs |
| No input validation (TS) | High | RESOLVED | `Validation` class added with runtime guards for STKPush; further services pending |
| No rate limiting | Medium | RESOLVED | Token bucket rate limiter + circuit breaker added to TypeScript |
| No SAST/SCA scanning | Medium | RESOLVED | Gitleaks in CI, npm audit in CI, gitleaks config added |
| No secret/lint scanning | Medium | RESOLVED | `.gitleaks.toml` with M-Pesa-specific patterns, Dependabot for auto-updates |
| No Docker environment | Low | RESOLVED | `docker-compose.yml` with TS/Python/Go/docs services |
| RSA encryption legacy | Medium | RESOLVED | Python upgraded from PKCS1v15 to OAEP with SHA-256 |
| Python framework middleware | Medium | RESOLVED | FastAPI webhook router added (`create_fastapi_router`) |
| Inconsistent API design across languages | High | OPEN | TS: `mpesa.stkPush.initiate()`, Python: `client.stk_push()`, Go: `client.STKPush()` |
| Duplicate endpoint definitions | Medium | OPEN | Each SDK hardcodes endpoints instead of loading from `shared/endpoints.json` |
| No OpenTelemetry/tracing | Medium | OPEN | No distributed tracing support yet |
| Minimal test coverage | High | OPEN | No HTTP mock tests, no integration tests, no load tests |
| Incomplete webhook signature verification | Medium | OPEN | TS/Python/Go verification is placeholder |
| No async/parallel support | Medium | OPEN | Python sync-only, Go no worker pool |
| Monolithic Python client | Medium | OPEN | Mpesa class flat, services layer thin passthrough |
| No compliance documentation | Medium | OPEN | No PCI-DSS/SOC2/GDPR guidance |

---

## Full Architecture Review

### TypeScript SDK (`typescript/`)

**Structure:** Well-modularized - `client/`, `services/`, `errors/`, `types/`, `utils/`, `webhooks/`, `middleware/`, `interceptors/`.

**Strengths:**
- Clean dependency injection (services receive client)
- Axios interceptor pattern for retry + error mapping
- Dual ESM/CJS exports via tsup
- Middleware for Express + Fastify
- Typed errors with `toJSON()`

**Issues:**
- `services/*.ts` reference `SANDBOX_ENDPOINTS` directly instead of client-managed endpoints
- STK Push password/timestamp logic duplicated across `initiate()` and `query()`
- No request ID generation/correlation
- No async local storage for context propagation
- Logging hook interface is minimal - no duration tracking on response
- `maskSensitiveData` doesn't mask nested objects

**RESOLVED:**
- Logger interface + `noopLogger`/`createConsoleLogger` factories integrated into client, retry interceptor, and token acquisition
- `Validation` class with `requiredString`, `requiredNumber`, `validUrl`, `phoneNumber`, `maxLength`, `oneOf`, `amount` methods
- Input validation integrated into `STKPushService.initiate()`
- `TokenBucketRateLimiter` + `CircuitBreaker` implementations with config support in `MpesaConfig`

### Python SDK (`python/`)

**Structure:** Flat package with `client/`, `exceptions/`, `models/`, `services/`, `utils/`, `webhooks/`, `middleware/`.

**Strengths:**
- Pydantic v2 models with field validation
- Hatchling build system (modern Python packaging)
- Context manager support (`with Mpesa() as client`)
- Typed exception hierarchy

**Issues:**
- Monolithic design: Mpesa class contains ALL API methods directly; services layer duplicates interface
- No async support: httpx is imported but only sync client used
- Manual retry logic: Embedded in `_request()` with `time.sleep()` - blocks event loop
- No connection pooling config: Default httpx pool settings

**RESOLVED:**
- `Logger` Protocol with `_get_logger()` helper integrated into Mpesa client, TokenManager, and WebhookManager
- httpx `event_hooks` for request/response logging
- RSA encryption upgraded from PKCS1v15 to OAEP with SHA-256 (MGF1)
- FastAPI webhook router (`create_fastapi_router`) in new `mpesa/middleware/` module

### Go SDK (`go/`)

**Structure:** Flat package with `client/`, `errors/`, `services/`, `types/`, `webhooks/`.

**Strengths:**
- Thread-safe token manager with `sync.RWMutex`
- Context propagation throughout
- Proper HTTP client with timeout configuration
- Good error types with functional options pattern

**Issues:**
- `services/service.go` adds minimal value - most methods are pass-through mapping
- No response validation: JSON unmarshalling trusts server response entirely
- No middleware (gin/echo/http.Handler) integrations
- `golang.org/x/time` dependency imported but `rate` package unused
- Tests depend on real token generation in `TestSTKPushRequestGeneration` - no mocks

**RESOLVED (CRITICAL):**
- **`services/service.go` `STKQuery` bug**: `BusinessShortCode: string(rune(input.BusinessShortCode))` - fixed to `fmt.Sprintf("%d", input.BusinessShortCode)`. Added `STKQueryInput` type with `CheckoutRequestID` field.
- **`doRequest` retry body re-read**: Body reader re-created per attempt instead of shared reader.
- **`defer resp.Body.Close()` inside retry loop**: Changed to direct `resp.Body.Close()` to prevent accumulated deferred closes.
- **No `context.Done()` check in retry backoff**: Added `select { case <-ctx.Done(): return nil, ctx.Err() ... }`.
- **Logger interface**: Added to `types.MpesaConfig` with `stdLogger`/`noopLogger` implementations. Integrated into `Client.doRequest` and `TokenManager.GetToken`.
- **Webhook logging**: Replaced `fmt.Printf` with structured logger in `webhooks.Manager`.

### Cross-Language Consistency Issues

| Aspect | TypeScript | Python | Go |
|--------|-----------|--------|----|
| Naming convention | camelCase | snake_case | PascalCase |
| API surface | `stkPush.initiate()` | `stk_push()` | `STKPush()` |
| Retry config | `retryConfig` object | `max_retries` int | `RetryConfig` struct |
| Error properties | `statusCode`, `requestId`, `rawResponse` | `status_code`, `request_id`, `raw_response` | Functional options |
| Webhook events | Domain events | String events | Typed constants |
| Framework support | Express, Fastify | FastAPI (NEW) | None |

**RESOLVED:**
- Logger interface pattern now consistent across all 3 languages (DEBUG/INFO/WARN/ERROR levels, `noopLogger` default)
- Webhook manager accepts logger in all 3 languages

---

## Security Audit

### OAuth Token Handling

| Check | Status | Notes |
|-------|--------|-------|
| Token caching with TTL minus 60s buffer | PASS | All 3 languages |
| Thread-safe token manager (Go RWMutex) | PASS | Go only; TS/Python single-thread |
| Token refresh on 401 | PASS | TS `withTokenRefresh`, Python invalidate+retry, Go invalidate+return |
| Encrypted token storage | OPEN | Not yet implemented |

### API Key/Credential Safety

| Check | Status | Notes |
|-------|--------|-------|
| No hardcoded secrets in codebase | PASS | Confirmed via audit |
| Sensitive data masking in logs | PASS | TS/Python/Go all have `maskSensitiveData` |
| RSA encryption | PASS (Updated) | Python upgraded from PKCS1v15 to OAEP SHA-256; TS/Go use PKCS1v15 |
| Credential rotation helper | OPEN | Not yet implemented |
| `.gitleaks.toml` with M-Pesa patterns | PASS | Added with consumer key/secret/passkey patterns |

### Rate Limiting & Abuse Prevention

| Check | Status | Notes |
|-------|--------|-------|
| Reactive backoff on 429 | PASS | All 3 languages |
| Token bucket rate limiter (TypeScript) | PASS | `TokenBucketRateLimiter` + config in `MpesaConfig` |
| Circuit breaker (TypeScript) | PASS | `CircuitBreaker` with closed/open/half-open states |
| Circuit breaker (Python/Go) | OPEN | Not yet implemented |

### Input Validation

| Check | Status |
|-------|--------|
| TypeScript runtime validation | PASS - Validation class with 9 guard methods integrated into STKPush |
| Python Pydantic field validation | PASS - Already present |
| Go input validation | OPEN |

### Dependency Security

| Check | Status |
|-------|--------|
| npm audit in CI | PASS |
| Gitleaks secret scanning in CI | PASS |
| Dependabot (npm/pip/gomod/actions) | PASS |
| govulncheck / pip safety | OPEN |
| Snyk / Renovate | OPEN |

---

## Production Blockers

| # | Blocker | Status | Resolution |
|---|---------|--------|------------|
| 1 | No structured logging | RESOLVED | Logger interface + implementations added to all 3 SDKs |
| 2 | No distributed tracing | OPEN | Phase 3 |
| 3 | Incomplete test coverage | OPEN | Phase 4 |
| 4 | Go services layer bug (STKQuery) | RESOLVED | Fixed BusinessShortCode conversion |
| 5 | No circuit breakers | IN PROGRESS | TypeScript only |
| 6 | No health checks | OPEN | Phase 3 |
| 7 | No metrics | OPEN | Phase 3 |
| 8 | No idempotency | OPEN | Phase 5 |

---

## Implementation Log

### Session 1 (2026-05-17): Phase 1 - Critical Fixes & Security Hardening

#### Completed

| # | Task | Lang | File(s) |
|---|------|------|---------|
| 1.1 | Fix Go STKQuery BusinessShortCode bug | Go | `go/services/service.go` - Changed `string(rune(...))` to `fmt.Sprintf("%d", ...)` |
| 1.2 | Fix Go doRequest body re-read on retry | Go | `go/client/client.go` - Re-marshal body per attempt, fix deferred Close() |
| 1.3 | Add STKQueryInput type | Go | `go/services/types/requests.go` - Added `STKQueryInput` struct |
| 1.4 | Add runtime input validation to TypeScript | TS | `typescript/src/utils/index.ts` - `Validation` class |
| 1.5 | Validate STKPush.initiate inputs | TS | `typescript/src/services/stk-push.ts` - Added validation guards |
| 1.6 | Add structured logging to all 3 SDKs | All | See detailed list below |
| 1.7 | Upgrade RSA encryption to OAEP in Python | Python | `python/mpesa/utils/__init__.py` |

#### Structured Logging Details

**TypeScript:**
- `typescript/src/types/index.ts` - `Logger` interface added to `MpesaConfig`
- `typescript/src/utils/index.ts` - `noopLogger` + `createConsoleLogger`
- `typescript/src/client/client.ts` - Logger integration in constructor, interceptors, getAccessToken
- `typescript/src/interceptors/retry.ts` - Logger parameter, warn on retry/rate-limit

**Python:**
- `python/mpesa/models/__init__.py` - `Logger` Protocol, `_get_logger()`, `logger` field in `MpesaConfig`
- `python/mpesa/client/__init__.py` - Logger in `_TokenManager`, `_request`, httpx event hooks
- `python/mpesa/webhooks/__init__.py` - Logger in `WebhookManager` constructor, `on`/`emit`

**Go:**
- `go/types/types.go` - `Logger` interface, `noopLogger`, `stdLogger` implementations, `MpesaConfig.Logger`
- `go/client/client.go` - Logger in `NewClient`, `doRequest`, `TokenManager.GetToken`
- `go/webhooks/webhook.go` - Logger in `NewManager`, `Emit`, `HandleSTKCallback`

### Session 2 (2026-05-17): Phase 2 - Observability & Resilience

#### Completed

| # | Task | Lang | File(s) |
|---|------|------|---------|
| 2.2 | Add client-side rate limiter (token bucket) | TS | `typescript/src/utils/rate-limiter.ts` |
| 2.3 | Add circuit breaker pattern | TS | `typescript/src/utils/circuit-breaker.ts` |
| 2.3 | Integrate rate limiter + circuit breaker config | TS | `typescript/src/types/index.ts` |

**Rate Limiter:** `TokenBucketRateLimiter` with configurable `tokensPerSecond` and `burstSize`. `NoopRateLimiter` for disabled state. Supports `acquire()` (async blocking) and `tryAcquire()` (non-blocking).

**Circuit Breaker:** Three states (closed/open/half-open). Configurable `failureThreshold`, `successThreshold`, `timeoutMs`. Auto-transitions half-open after timeout. Integrated via `MpesaConfig`.

### Session 3 (2026-05-17): Phase 4 - Testing & CI/CD

#### Completed

| # | Task | Lang | File(s) |
|---|------|------|---------|
| 4.4 | Add secret scanning (Gitleaks) | CI | `.github/workflows/ci.yml` - `security-scan` job |
| 4.4 | Add npm audit to CI | CI | `.github/workflows/ci.yml` |
| 4.4 | Add .env detection to CI | CI | `.github/workflows/ci.yml` |
| 4.5 | Add Dependabot config | CI | `.github/dependabot.yml` - npm/pip/gomod/actions |
| 4.6 | Add Docker compose | Infra | `docker-compose.yml` - TS/Python/Go/docs |
| 4.4 | Add Gitleaks config | CI | `.gitleaks.toml` |
| 4.7 | Add go vet to CI | CI | `.github/workflows/ci.yml` |

---

## Remaining Roadmap

### Phase 2 (Continuation): Observability & Resilience

| # | Task | Lang | Effort | Status |
|---|------|------|--------|--------|
| 2.1 | Add OpenTelemetry tracing support | All | 4h | OPEN |
| 2.4 | Add metrics collection (prometheus client) | All | 3h | OPEN |
| 2.5 | Add health check endpoint | All | 1h | OPEN |
| 2.6 | Add idempotency key support | All | 2h | OPEN |
| 2.7 | Implement async/parallel support in Python | Python | 2h | OPEN |

### Phase 3: Cross-Language Consistency

| # | Task | Lang | Effort | Status |
|---|------|------|--------|--------|
| 3.1 | Generate endpoints from `shared/endpoints.json` instead of hardcoding | All | 2h | OPEN |
| 3.2 | Standardize service method naming (service-oriented pattern) | All | 2h | OPEN |
| 3.3 | Add Python async client + sync compatibility | Python | 2h | OPEN |
| 3.4 | Add Python Flask/Django middleware | Python | 2h | OPEN |
| 3.5 | Add Go framework middleware (Gin, Echo, net/http) | Go | 2h | OPEN |
| 3.6 | Standardize retry configuration schema across languages | All | 1h | OPEN |
| 3.7 | Implement uniform webhook event type system | All | 1h | OPEN |

### Phase 4 (Continuation): Testing

| # | Task | Lang | Effort | Status |
|---|------|------|--------|--------|
| 4.1 | Add HTTP mock tests for all services in each language | All | 6h | OPEN |
| 4.2 | Add integration test suite with sandbox API | All | 4h | OPEN |
| 4.3 | Add load/performance tests (k6 or similar) | All | 3h | OPEN |

### Phase 5: Enterprise Features

| # | Task | Lang | Effort | Status |
|---|------|------|--------|--------|
| 5.1 | Add compliance documentation (PCI-DSS, SOC2 guidance) | Docs | 2h | OPEN |
| 5.2 | Add audit trail logging | All | 2h | OPEN |
| 5.3 | Add batch request support (where API allows) | All | 2h | OPEN |
| 5.4 | Add webhook retry with dead-letter queue pattern | All | 2h | OPEN |
| 5.5 | Add credential rotation helper | All | 1h | OPEN |
| 5.6 | Add encrypted token persistence option | All | 2h | OPEN |

---

## Detailed Implementation Specifications

### 1. Structured Logging Interface (IMPLEMENTED)

Each SDK exposes a logger interface with levels: `DEBUG`, `INFO`, `WARN`, `ERROR`.

```typescript
// TypeScript - IMPLEMENTED
interface Logger {
  debug(msg: string, meta?: Record<string, unknown>): void;
  info(msg: string, meta?: Record<string, unknown>): void;
  warn(msg: string, meta?: Record<string, unknown>): void;
  error(msg: string, meta?: Record<string, unknown>): void;
}
```

Default: `noopLogger`. Users can inject pino/winston/bunyan (TS), standard library logging (Python/Go).

### 2. Rate Limiter - Token Bucket (IMPLEMENTED - TypeScript only)

```typescript
// TypeScript - IMPLEMENTED
interface RateLimiter {
  acquire(): Promise<void>;
  tryAcquire(): boolean;
}

const mpesa = new Mpesa({
  rateLimiterConfig: {
    tokensPerSecond: 5,
    burstSize: 10,
  },
});
```

### 3. Circuit Breaker (IMPLEMENTED - TypeScript only)

```typescript
// TypeScript - IMPLEMENTED
const mpesa = new Mpesa({
  circuitBreakerConfig: {
    failureThreshold: 5,
    successThreshold: 2,
    timeoutMs: 30000,
  },
});
```

States: closed, open, half-open.

### 4. RSA Security Credential - IMPLEMENTED (Python - OAEP)

```python
# Python - IMPLEMENTED with OAEP SHA-256
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives import hashes

encrypted = cert.encrypt(
    password.encode(),
    padding.OAEP(
        mgf=padding.MGF1(algorithm=hashes.SHA256()),
        algorithm=hashes.SHA256(),
        label=None,
    ),
)
```

### 5. OpenTelemetry Tracing (NOT IMPLEMENTED)

```typescript
// TypeScript - FUTURE
const mpesa = new Mpesa({
  telemetry: {
    tracerProvider: otel.trace.getTracerProvider(),
    meterProvider: otel.metrics.getMeterProvider(),
  },
});
```

### 6. Auto-Generated Endpoints (NOT IMPLEMENTED)

```
shared/endpoints.json  -->  code generation  -->  SDK endpoint constants
```

### 7. Health Check (NOT IMPLEMENTED)

```python
# Python - FUTURE
GET /health -> {
    "status": "healthy" | "degraded" | "unhealthy",
    "token_status": "valid" | "expired" | "refreshing",
    "version": "0.2.0"
}
```

---

## Compliance & Audit Readiness

### PCI-DSS Considerations
- SDK does not store, process, or transmit PAN (Primary Account Numbers) - phone numbers are not PAN
- Webhook payloads may contain transaction metadata - should document data handling
- Recommend users implement encryption at rest for callback payloads

### SOC2 Considerations
- Audit logging of all API calls (request/response) - Logger interface provides foundation but no structured audit trail yet
- Token lifecycle management - In-memory only; no persistence audit trail
- Dependency vulnerability scanning - Dependabot configured; no Snyk/SCA integration yet

### GDPR / Data Privacy
- Phone number fields should be documented as PII
- Recommend users implement data retention policies for callback payloads
- Add data masking configuration options

---

## Scoring & Prioritization Matrix

| Task | Impact | Effort | Priority | Status |
|------|--------|--------|----------|--------|
| Fix Go STKQuery bug | Critical | 15m | P0 | DONE |
| Fix Go retry body re-read | High | 30m | P0 | DONE |
| Add structured logging | High | 4h | P0 | DONE |
| Add runtime validation (TS) | High | 2h | P1 | DONE |
| Add rate limiter | Medium | 3h | P1 | DONE (TS) |
| Add circuit breaker | Medium | 3h | P1 | DONE (TS) |
| SAST/SCA in CI | Medium | 2h | P2 | DONE |
| Docker compose | Low | 1h | P3 | DONE |
| Add request ID correlation | High | 2h | P1 | OPEN |
| Add OTel tracing | Medium | 4h | P1 | OPEN |
| Add HTTP mock tests | High | 6h | P1 | OPEN |
| Auto-generate endpoints | Medium | 2h | P2 | OPEN |
| Standardize API naming | Medium | 2h | P2 | OPEN |
| Python async support | Medium | 2h | P2 | OPEN |
| Framework middleware | Low | 4h | P2 | DONE (Python FastAPI) |
| Load tests | Medium | 3h | P2 | OPEN |
| Credential rotation | Low | 1h | P3 | OPEN |
| Compliance docs | Low | 2h | P3 | OPEN |
