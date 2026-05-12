import hashlib
import hmac
from typing import Any, Callable

from mpesa.models import STKCallbackPayload


WebhookHandler = Callable[[str, Any], None]


class WebhookManager:
    def __init__(self) -> None:
        self._handlers: dict[str, list[WebhookHandler]] = {}

    def on(self, event_type: str, handler: WebhookHandler) -> None:
        if event_type not in self._handlers:
            self._handlers[event_type] = []
        self._handlers[event_type].append(handler)

    def off(self, event_type: str, handler: WebhookHandler) -> None:
        if event_type in self._handlers:
            self._handlers[event_type] = [h for h in self._handlers[event_type] if h != handler]

    def emit(self, event_type: str, payload: Any) -> None:
        handlers = self._handlers.get(event_type, [])
        for handler in handlers:
            try:
                handler(event_type, payload)
            except Exception as e:
                import logging
                logging.getLogger("mpesa").error(f"Webhook handler error: {e}")

    def parse_stk_callback(self, body: dict) -> dict:
        payload = STKCallbackPayload(**body)
        callback = payload.Body.stkCallback
        result = {
            "success": callback.ResultCode == 0,
            "merchant_request_id": callback.MerchantRequestID,
            "checkout_request_id": callback.CheckoutRequestID,
            "result_code": callback.ResultCode,
            "result_description": callback.ResultDesc,
        }

        if callback.CallbackMetadata:
            for item in callback.CallbackMetadata.Item:
                if item.Name == "Amount":
                    result["amount"] = float(item.Value) if item.Value else None
                elif item.Name == "MpesaReceiptNumber":
                    result["receipt_number"] = str(item.Value) if item.Value else None
                elif item.Name == "TransactionDate":
                    result["transaction_date"] = str(item.Value) if item.Value else None
                elif item.Name == "PhoneNumber":
                    result["phone_number"] = str(item.Value) if item.Value else None

        return result

    def parse_c2b_validation_response(self, accept: bool = True) -> dict:
        if accept:
            return {"ResultCode": "0", "ResultDesc": "Accepted"}
        return {"ResultCode": "C2B00011", "ResultDesc": "Rejected"}

    def verify_signature(self, payload: str, signature: str, secret: str) -> bool:
        expected = hmac.new(
            secret.encode(), payload.encode(), hashlib.sha256
        ).hexdigest()
        return hmac.compare_digest(expected, signature)


__all__ = [
    "WebhookManager",
    "WebhookHandler",
]
