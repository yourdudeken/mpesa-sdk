import time
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
)
from mpesa.utils import generate_password, generate_timestamp

RETRYABLE_STATUS_CODES = {408, 429, 500, 502, 503, 504}


class _TokenManager:
    def __init__(self, client: httpx.Client, config: MpesaConfig) -> None:
        self._client = client
        self._config = config
        self._token: Optional[str] = None
        self._expires_at: float = 0.0

    def get_token(self) -> str:
        if self._token and time.time() < self._expires_at:
            return self._token

        url = get_full_url(self._config.environment, ENDPOINTS["AUTH"])
        response = self._client.get(
            url,
            params={"grant_type": "client_credentials"},
            auth=(self._config.consumer_key, self._config.consumer_secret),
        )
        data = response.raise_for_status().json()
        token_data = AccessTokenResponse(**data)
        self._token = token_data.access_token
        self._expires_at = time.time() + token_data.expires_in - 60
        return self._token

    def invalidate(self) -> None:
        self._token = None
        self._expires_at = 0.0


class Mpesa:
    def __init__(self, config: MpesaConfig | dict[str, Any]) -> None:
        if isinstance(config, dict):
            config = MpesaConfig(**config)

        self._config = config
        self._client = httpx.Client(
            base_url=get_full_url(config.environment, ""),
            timeout=config.timeout,
            headers={
                "Content-Type": "application/json",
                "Accept": "application/json",
            },
        )
        self._token_manager = _TokenManager(self._client, config)

    def _request(self, method: str, url: str, json_data: Optional[dict] = None) -> dict:
        last_error: Optional[Exception] = None

        for attempt in range(self._config.max_retries + 1):
            try:
                token = self._token_manager.get_token()
                response = self._client.request(
                    method=method,
                    url=url,
                    json=json_data,
                    headers={"Authorization": f"Bearer {token}"},
                )

                if response.status_code in RETRYABLE_STATUS_CODES and attempt < self._config.max_retries:
                    time.sleep(min(2**attempt * 1.0, 30.0))
                    continue

                if response.status_code == 401:
                    self._token_manager.invalidate()
                    raise AuthenticationError(
                        "Authentication failed.",
                        status_code=401,
                        raw_response=response.text,
                    )

                if response.status_code == 429:
                    retry_after = int(response.headers.get("Retry-After", "60"))
                    raise RateLimitError(
                        "Rate limit exceeded.",
                        status_code=429,
                        retry_after=retry_after,
                        raw_response=response.text,
                    )

                response.raise_for_status()
                return response.json()

            except httpx.TimeoutException as e:
                last_error = TimeoutError("Request timed out.", cause=e)
                if attempt < self._config.max_retries:
                    time.sleep(min(2**attempt * 1.0, 30.0))
                    continue
                raise last_error

            except httpx.ConnectError as e:
                last_error = APIConnectionError("Connection failed.", cause=e)
                if attempt < self._config.max_retries:
                    time.sleep(min(2**attempt * 1.0, 30.0))
                    continue
                raise last_error

            except (AuthenticationError, RateLimitError):
                raise

            except httpx.HTTPStatusError as e:
                raise MpesaAPIError(
                    str(e),
                    status_code=e.response.status_code,
                    raw_response=e.response.text,
                )

        if last_error:
            raise last_error
        raise MpesaAPIError("Request failed after retries.")

    def _post(self, endpoint_key: str, data: dict) -> dict:
        url = get_full_url(self._config.environment, ENDPOINTS[endpoint_key])
        return self._request("POST", url, data)

    def stk_push(self, request: STKPushRequest | dict) -> STKPushResponse:
        if isinstance(request, dict):
            request = STKPushRequest(**request)
        if not request.Password and self._config.passkey:
            timestamp = request.Timestamp or generate_timestamp()
            request.Password = generate_password(request.BusinessShortCode, self._config.passkey, timestamp)
            request.Timestamp = timestamp
        result = self._post("STK_PUSH", request.model_dump())
        return STKPushResponse(**result)

    def stk_query(self, request: STKQueryRequest | dict) -> STKQueryResponse:
        if isinstance(request, dict):
            request = STKQueryRequest(**request)
        if not request.Password and self._config.passkey:
            timestamp = request.Timestamp or generate_timestamp()
            request.Password = generate_password(request.BusinessShortCode, self._config.passkey, timestamp)
            request.Timestamp = timestamp
        result = self._post("STK_QUERY", request.model_dump())
        return STKQueryResponse(**result)

    def c2b_register_url(self, request: C2BRegisterURLRequest | dict) -> C2BResponse:
        if isinstance(request, dict):
            request = C2BRegisterURLRequest(**request)
        result = self._post("C2B_REGISTER_URL", request.model_dump())
        return C2BResponse(**result)

    def c2b_simulate(self, request: C2BSimulateRequest | dict) -> C2BResponse:
        if isinstance(request, dict):
            request = C2BSimulateRequest(**request)
        result = self._post("C2B_SIMULATE", request.model_dump())
        return C2BResponse(**result)

    def b2c(self, request: B2CRequest | dict) -> B2CResponse:
        if isinstance(request, dict):
            request = B2CRequest(**request)
        result = self._post("B2C", request.model_dump())
        return B2CResponse(**result)

    def b2b(self, request: B2BRequest | dict) -> B2BResponse:
        if isinstance(request, dict):
            request = B2BRequest(**request)
        result = self._post("B2B", request.model_dump())
        return B2BResponse(**result)

    def reversal(self, request: ReversalRequest | dict) -> ReversalResponse:
        if isinstance(request, dict):
            request = ReversalRequest(**request)
        result = self._post("REVERSAL", request.model_dump())
        return ReversalResponse(**result)

    def transaction_status(self, request: TransactionStatusRequest | dict) -> TransactionStatusResponse:
        if isinstance(request, dict):
            request = TransactionStatusRequest(**request)
        result = self._post("TRANSACTION_STATUS", request.model_dump())
        return TransactionStatusResponse(**result)

    def account_balance(self, request: AccountBalanceRequest | dict) -> AccountBalanceResponse:
        if isinstance(request, dict):
            request = AccountBalanceRequest(**request)
        result = self._post("ACCOUNT_BALANCE", request.model_dump())
        return AccountBalanceResponse(**result)

    def dynamic_qr(self, request: DynamicQRRequest | dict) -> DynamicQRResponse:
        if isinstance(request, dict):
            request = DynamicQRRequest(**request)
        result = self._post("DYNAMIC_QR", request.model_dump())
        return DynamicQRResponse(**result)

    def close(self) -> None:
        self._client.close()

    def __enter__(self) -> "Mpesa":
        return self

    def __exit__(self, *args: Any) -> None:
        self.close()
