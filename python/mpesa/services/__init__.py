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
from mpesa.client import Mpesa


class STKPushService:
    def __init__(self, client: Mpesa) -> None:
        self._client = client

    def initiate(self, request: STKPushRequest | dict) -> STKPushResponse:
        return self._client.stk_push(request)

    def query(self, request: STKQueryRequest | dict) -> STKQueryResponse:
        return self._client.stk_query(request)


class C2BService:
    def __init__(self, client: Mpesa) -> None:
        self._client = client

    def register_url(self, request: C2BRegisterURLRequest | dict) -> C2BResponse:
        return self._client.c2b_register_url(request)

    def simulate(self, request: C2BSimulateRequest | dict) -> C2BResponse:
        return self._client.c2b_simulate(request)


class B2CService:
    def __init__(self, client: Mpesa) -> None:
        self._client = client

    def send(self, request: B2CRequest | dict) -> B2CResponse:
        return self._client.b2c(request)


class B2BService:
    def __init__(self, client: Mpesa) -> None:
        self._client = client

    def send(self, request: B2BRequest | dict) -> B2BResponse:
        return self._client.b2b(request)


class ReversalService:
    def __init__(self, client: Mpesa) -> None:
        self._client = client

    def reverse(self, request: ReversalRequest | dict) -> ReversalResponse:
        return self._client.reversal(request)


class TransactionStatusService:
    def __init__(self, client: Mpesa) -> None:
        self._client = client

    def query(self, request: TransactionStatusRequest | dict) -> TransactionStatusResponse:
        return self._client.transaction_status(request)


class AccountBalanceService:
    def __init__(self, client: Mpesa) -> None:
        self._client = client

    def query(self, request: AccountBalanceRequest | dict) -> AccountBalanceResponse:
        return self._client.account_balance(request)


class DynamicQRService:
    def __init__(self, client: Mpesa) -> None:
        self._client = client

    def generate(self, request: DynamicQRRequest | dict) -> DynamicQRResponse:
        return self._client.dynamic_qr(request)


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
