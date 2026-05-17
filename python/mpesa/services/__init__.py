from typing import Any, Callable, Optional

from mpesa.models import (
    AccountBalanceRequest,
    AccountBalanceResponse,
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

PostFn = Callable[[str, dict], dict]


class STKPushService:
    def __init__(self, post: PostFn, config: MpesaConfig) -> None:
        self._post = post
        self._config = config

    def initiate(self, request: STKPushRequest | dict) -> STKPushResponse:
        if isinstance(request, dict):
            request = STKPushRequest(**request)
        if not request.Password and self._config.passkey:
            timestamp = request.Timestamp or generate_timestamp()
            request.Password = generate_password(request.BusinessShortCode, self._config.passkey, timestamp)
            request.Timestamp = timestamp
        result = self._post("STK_PUSH", request.model_dump())
        return STKPushResponse(**result)

    def query(self, request: STKQueryRequest | dict) -> STKQueryResponse:
        if isinstance(request, dict):
            request = STKQueryRequest(**request)
        if not request.Password and self._config.passkey:
            timestamp = request.Timestamp or generate_timestamp()
            request.Password = generate_password(request.BusinessShortCode, self._config.passkey, timestamp)
            request.Timestamp = timestamp
        result = self._post("STK_QUERY", request.model_dump())
        return STKQueryResponse(**result)


class C2BService:
    def __init__(self, post: PostFn) -> None:
        self._post = post

    def register_url(self, request: C2BRegisterURLRequest | dict) -> C2BResponse:
        if isinstance(request, dict):
            request = C2BRegisterURLRequest(**request)
        result = self._post("C2B_REGISTER_URL", request.model_dump())
        return C2BResponse(**result)

    def simulate(self, request: C2BSimulateRequest | dict) -> C2BResponse:
        if isinstance(request, dict):
            request = C2BSimulateRequest(**request)
        result = self._post("C2B_SIMULATE", request.model_dump())
        return C2BResponse(**result)


class B2CService:
    def __init__(self, post: PostFn) -> None:
        self._post = post

    def send(self, request: B2CRequest | dict) -> B2CResponse:
        if isinstance(request, dict):
            request = B2CRequest(**request)
        result = self._post("B2C", request.model_dump())
        return B2CResponse(**result)


class B2BService:
    def __init__(self, post: PostFn) -> None:
        self._post = post

    def send(self, request: B2BRequest | dict) -> B2BResponse:
        if isinstance(request, dict):
            request = B2BRequest(**request)
        result = self._post("B2B", request.model_dump())
        return B2BResponse(**result)


class ReversalService:
    def __init__(self, post: PostFn) -> None:
        self._post = post

    def reverse(self, request: ReversalRequest | dict) -> ReversalResponse:
        if isinstance(request, dict):
            request = ReversalRequest(**request)
        result = self._post("REVERSAL", request.model_dump())
        return ReversalResponse(**result)


class TransactionStatusService:
    def __init__(self, post: PostFn) -> None:
        self._post = post

    def query(self, request: TransactionStatusRequest | dict) -> TransactionStatusResponse:
        if isinstance(request, dict):
            request = TransactionStatusRequest(**request)
        result = self._post("TRANSACTION_STATUS", request.model_dump())
        return TransactionStatusResponse(**result)


class AccountBalanceService:
    def __init__(self, post: PostFn) -> None:
        self._post = post

    def query(self, request: AccountBalanceRequest | dict) -> AccountBalanceResponse:
        if isinstance(request, dict):
            request = AccountBalanceRequest(**request)
        result = self._post("ACCOUNT_BALANCE", request.model_dump())
        return AccountBalanceResponse(**result)


class DynamicQRService:
    def __init__(self, post: PostFn) -> None:
        self._post = post

    def generate(self, request: DynamicQRRequest | dict) -> DynamicQRResponse:
        if isinstance(request, dict):
            request = DynamicQRRequest(**request)
        result = self._post("DYNAMIC_QR", request.model_dump())
        return DynamicQRResponse(**result)


__all__ = [
    "STKPushService",
    "C2BService",
    "B2CService",
    "B2BService",
    "ReversalService",
    "TransactionStatusService",
    "AccountBalanceService",
    "DynamicQRService",
]
