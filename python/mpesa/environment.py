import json
from pathlib import Path
from typing import Literal

_SHARED_ENDPOINTS_PATH = Path(__file__).resolve().parent.parent.parent.parent / "shared" / "endpoints.json"

SANDBOX_BASE_URL = "https://sandbox.safaricom.co.ke"
PRODUCTION_BASE_URL = "https://api.safaricom.co.ke"

_ENDPOINT_KEYS = [
    "auth", "stk_push", "stk_query", "c2b_register_url", "c2b_simulate",
    "b2c", "b2b", "reversal", "transaction_status", "account_balance", "dynamic_qr",
]

ENDPOINTS: dict[str, str] = {}

if _SHARED_ENDPOINTS_PATH.exists():
    with open(_SHARED_ENDPOINTS_PATH) as f:
        data = json.load(f)
    sandbox_eps = data["environments"]["sandbox"]["endpoints"]
    for key in _ENDPOINT_KEYS:
        ENDPOINTS[key.upper()] = sandbox_eps[key]
else:
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
