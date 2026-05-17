import asyncio
import logging
import uuid
from typing import Any, Optional

import httpx

from mpesa.environment import ENDPOINTS, get_full_url
from mpesa.exceptions import (
    AuthenticationError,
    APIConnectionError,
    MpesaAPIError,
    RateLimitError,
    TimeoutError,
)
from mpesa.models import (
    AccountBalanceRequest,
    AccountBalanceResponse,
    AccessTokenResponse,
    B2BRequest,
    B2BResponse,
    B2CRequest,
    B2CResponse,
    C2BRegisterURLRequest,
    C2BResponse,
    C2BSimulateRequest,
    DynamicQRRequest,
    DynamicQRResponse,
    MpesaConfig,
    ReversalRequest,
    ReversalResponse,
    STKPushRequest,
    STKPushResponse,
    STKQueryRequest,
    STKQueryResponse,
    TransactionStatusRequest,
    TransactionStatusResponse,
    _get_logger,
)
from mpesa.utils import generate_password, generate_timestamp
from mpesa.utils.circuit_breaker import CircuitBreaker, CircuitBreakerOpenError, CircuitBreakerConfig
from mpesa.utils.rate_limiter import TokenBucketRateLimiter, NoopRateLimiter, RateLimiterConfig

RETRYABLE_STATUS_CODES = {408, 429, 500, 502, 503, 504}


def _generate_request_id() -> str:
    return f"mpesa-{uuid.uuid4().hex[:16]}"


class _AsyncTokenManager:
    def __init__(self, client: httpx.AsyncClient, config: MpesaConfig) -> None:
        self._client = client
        self._config = config
        self._token: Optional[str] = None
        self._expires_at: float = 0.0
        self._lock = asyncio.Lock()
        self._logger = _get_logger(config.logger)

    async def get_token(self) -> str:
        if self._token and asyncio.get_event_loop().time() < self._expires_at:
            return self._token

        async with self._lock:
            if self._token and asyncio.get_event_loop().time() < self._expires_at:
                return self._token

            self._logger.debug("Fetching new access token")

            url = get_full_url(self._config.environment, ENDPOINTS["AUTH"])
            response = await self._client.get(
                url,
                params={"grant_type": "client_credentials"},
                auth=(self._config.consumer_key, self._config.consumer_secret),
            )
            data = response.raise_for_status().json()
            token_data = AccessTokenResponse(**data)
            self._token = token_data.access_token
            self._expires_at = asyncio.get_event_loop().time() + token_data.expires_in - 60

            self._logger.debug("Access token acquired", extra={"expires_in": token_data.expires_in})
            return self._token

    def invalidate(self) -> None:
        self._token = None
        self._expires_at = 0.0
        self._logger.warning("Access token invalidated")


class AsyncMpesa:
    def __init__(self, config: MpesaConfig | dict[str, Any]) -> None:
        if isinstance(config, dict):
            config = MpesaConfig(**config)

        self._config = config
        self._logger = _get_logger(config.logger)

        cb_cfg = config.circuit_breaker_config or {}
        self._circuit_breaker = CircuitBreaker(
            failure_threshold=cb_cfg.get("failure_threshold", 5),
            success_threshold=cb_cfg.get("success_threshold", 2),
            timeout_ms=cb_cfg.get("timeout_ms", 30000),
        )

        rl_cfg = config.rate_limiter_config
        if rl_cfg:
            self._rate_limiter = TokenBucketRateLimiter(
                tokens_per_second=rl_cfg.get("tokens_per_second", 5),
                burst_size=rl_cfg.get("burst_size", 10),
            )
        else:
            self._rate_limiter = NoopRateLimiter()

        self._client = httpx.AsyncClient(
            base_url=get_full_url(config.environment, ""),
            timeout=config.timeout,
            headers={
                "Content-Type": "application/json",
                "Accept": "application/json",
            },
            event_hooks={
                "request": [self._log_request],
                "response": [self._log_response],
            },
        )
        self._token_manager = _AsyncTokenManager(self._client, config)
        self._logger.info("Async M-Pesa client initialized", extra={
            "environment": config.environment,
            "timeout": config.timeout,
            "max_retries": config.retry_config.max_retries,
        })

    def _log_request(self, request: httpx.Request) -> None:
        self._logger.debug("Outgoing request",
                           extra={"method": request.method, "url": str(request.url)})

    def _log_response(self, response: httpx.Response) -> None:
        self._logger.debug("Response received",
                           extra={"status": response.status_code, "url": str(response.url)})

    async def _request(self, method: str, url: str, json_data: Optional[dict] = None) -> dict:
        request_id = _generate_request_id()

        while not self._rate_limiter.try_acquire():
            await asyncio.sleep(0.01)

        async def do_request() -> dict:
            last_error: Optional[Exception] = None
            for attempt in range(self._config.retry_config.max_retries + 1):
                try:
                    if attempt > 0:
                        self._logger.warning("Retrying request",
                                             extra={"attempt": attempt, "url": url, "request_id": request_id})

                    token = await self._token_manager.get_token()
                    headers = {
                        "Authorization": f"Bearer {token}",
                        "X-Request-ID": request_id,
                    }
                    if self._config.enable_idempotency and method.upper() == "POST":
                        headers["X-Idempotency-Key"] = self._generate_idempotency_key()
                    response = await self._client.request(
                        method=method,
                        url=url,
                        json=json_data,
                        headers=headers,
                    )

                    if response.status_code in RETRYABLE_STATUS_CODES and attempt < self._config.retry_config.max_retries:
                        delay = min(2 ** attempt * 1.0, 30.0)
                        self._logger.warning("Retryable status code, backing off",
                                             extra={"status": response.status_code, "delay": delay, "attempt": attempt, "request_id": request_id})
                        await asyncio.sleep(delay)
                        continue

                    if response.status_code == 401:
                        self._token_manager.invalidate()
                        raise AuthenticationError(
                            "Authentication failed.",
                            status_code=401,
                            request_id=request_id,
                            raw_response=response.text,
                        )

                    if response.status_code == 429:
                        retry_after = int(response.headers.get("Retry-After", "60"))
                        raise RateLimitError(
                            "Rate limit exceeded.",
                            status_code=429,
                            retry_after=retry_after,
                            request_id=request_id,
                            raw_response=response.text,
                        )

                    response.raise_for_status()
                    json_result = response.json()
                    self._logger.debug("Request successful",
                                       extra={"method": method, "url": url, "status": response.status_code, "request_id": request_id})
                    return json_result

                except httpx.TimeoutException as e:
                    last_error = TimeoutError("Request timed out.", cause=e, request_id=request_id)
                    if attempt < self._config.retry_config.max_retries:
                        delay = min(2 ** attempt * 1.0, 30.0)
                        await asyncio.sleep(delay)
                        continue
                    raise last_error

                except httpx.ConnectError as e:
                    last_error = APIConnectionError("Connection failed.", cause=e, request_id=request_id)
                    if attempt < self._config.retry_config.max_retries:
                        delay = min(2 ** attempt * 1.0, 30.0)
                        await asyncio.sleep(delay)
                        continue
                    raise last_error

                except (AuthenticationError, RateLimitError):
                    raise

                except httpx.HTTPStatusError as e:
                    self._logger.error("API error response",
                                       extra={"status": e.response.status_code, "body": e.response.text, "request_id": request_id})
                    raise MpesaAPIError(
                        str(e),
                        status_code=e.response.status_code,
                        request_id=request_id,
                        raw_response=e.response.text,
                    )

            if last_error:
                raise last_error
            raise MpesaAPIError("Request failed after retries.", request_id=request_id)

        return self._circuit_breaker.call(do_request)

    async def _post(self, endpoint_key: str, data: dict) -> dict:
        url = get_full_url(self._config.environment, ENDPOINTS[endpoint_key])
        return await self._request("POST", url, data)

    async def stk_push(self, request: STKPushRequest | dict) -> STKPushResponse:
        if isinstance(request, dict):
            request = STKPushRequest(**request)
        if not request.Password and self._config.passkey:
            timestamp = request.Timestamp or generate_timestamp()
            request.Password = generate_password(request.BusinessShortCode, self._config.passkey, timestamp)
            request.Timestamp = timestamp
        result = await self._post("STK_PUSH", request.model_dump())
        return STKPushResponse(**result)

    async def stk_query(self, request: STKQueryRequest | dict) -> STKQueryResponse:
        if isinstance(request, dict):
            request = STKQueryRequest(**request)
        if not request.Password and self._config.passkey:
            timestamp = request.Timestamp or generate_timestamp()
            request.Password = generate_password(request.BusinessShortCode, self._config.passkey, timestamp)
            request.Timestamp = timestamp
        result = await self._post("STK_QUERY", request.model_dump())
        return STKQueryResponse(**result)

    async def c2b_register_url(self, request: C2BRegisterURLRequest | dict) -> C2BResponse:
        if isinstance(request, dict):
            request = C2BRegisterURLRequest(**request)
        result = await self._post("C2B_REGISTER_URL", request.model_dump())
        return C2BResponse(**result)

    async def c2b_simulate(self, request: C2BSimulateRequest | dict) -> C2BResponse:
        if isinstance(request, dict):
            request = C2BSimulateRequest(**request)
        result = await self._post("C2B_SIMULATE", request.model_dump())
        return C2BResponse(**result)

    async def b2c(self, request: B2CRequest | dict) -> B2CResponse:
        if isinstance(request, dict):
            request = B2CRequest(**request)
        result = await self._post("B2C", request.model_dump())
        return B2CResponse(**result)

    async def b2b(self, request: B2BRequest | dict) -> B2BResponse:
        if isinstance(request, dict):
            request = B2BRequest(**request)
        result = await self._post("B2B", request.model_dump())
        return B2BResponse(**result)

    async def reversal(self, request: ReversalRequest | dict) -> ReversalResponse:
        if isinstance(request, dict):
            request = ReversalRequest(**request)
        result = await self._post("REVERSAL", request.model_dump())
        return ReversalResponse(**result)

    async def transaction_status(self, request: TransactionStatusRequest | dict) -> TransactionStatusResponse:
        if isinstance(request, dict):
            request = TransactionStatusRequest(**request)
        result = await self._post("TRANSACTION_STATUS", request.model_dump())
        return TransactionStatusResponse(**result)

    async def account_balance(self, request: AccountBalanceRequest | dict) -> AccountBalanceResponse:
        if isinstance(request, dict):
            request = AccountBalanceRequest(**request)
        result = await self._post("ACCOUNT_BALANCE", request.model_dump())
        return AccountBalanceResponse(**result)

    async def dynamic_qr(self, request: DynamicQRRequest | dict) -> DynamicQRResponse:
        if isinstance(request, dict):
            request = DynamicQRRequest(**request)
        result = await self._post("DYNAMIC_QR", request.model_dump())
        return DynamicQRResponse(**result)

    async def rotate_credentials(self, consumer_key: str, consumer_secret: str) -> None:
        self._config.consumer_key = consumer_key
        self._config.consumer_secret = consumer_secret
        self._token_manager.invalidate()
        self._logger.info("Credentials rotated")

    async def close(self) -> None:
        await self._client.aclose()

    async def __aenter__(self) -> "AsyncMpesa":
        return self

    async def __aexit__(self, *args: Any) -> None:
        await self.close()
