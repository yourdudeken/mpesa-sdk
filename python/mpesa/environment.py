from typing import Literal

SANDBOX_BASE_URL = "https://sandbox.safaricom.co.ke"
PRODUCTION_BASE_URL = "https://api.safaricom.co.ke"

ENDPOINTS = {
    "AUTH": "/oauth/v1/generate",
    "STK_PUSH": "/mpesa/stkpush/v1/processrequest",
    "STK_QUERY": "/mpesa/stkpushquery/v1/query",
    "C2B_REGISTER_URL": "/mpesa/c2b/v2/registerurl",
    "C2B_SIMULATE": "/mpesa/c2b/v2/simulate",
    "B2C": "/mpesa/b2c/v3/paymentrequest",
    "B2B": "/mpesa/b2b/v1/paymentrequest",
    "REVERSAL": "/mpesa/reversal/v1/request",
    "TRANSACTION_STATUS": "/mpesa/transactionstatus/v1/query",
    "ACCOUNT_BALANCE": "/mpesa/accountbalance/v1/query",
    "DYNAMIC_QR": "/mpesa/qrcode/v1/generate",
}

Environment = Literal["sandbox", "production"]


def get_base_url(environment: Environment) -> str:
    return SANDBOX_BASE_URL if environment == "sandbox" else PRODUCTION_BASE_URL


def get_full_url(environment: Environment, endpoint_path: str) -> str:
    return f"{get_base_url(environment)}{endpoint_path}"
