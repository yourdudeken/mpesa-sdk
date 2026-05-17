from mpesa.webhooks import WebhookManager


__all__ = [
    "create_fastapi_router",
]


def create_fastapi_router(webhook_manager: WebhookManager, secret: str = "") -> "APIRouter":
    try:
        from fastapi import APIRouter, HTTPException, Request
    except ImportError:
        raise ImportError("fastapi is required. Install with: pip install yourdudeken-mpesa-sdk[fastapi]")

    router = APIRouter()

    @router.post("/mpesa/webhook")
    async def handle_webhook(request: Request):
        body = await request.json()

        if secret:
            signature = request.headers.get("x-mpesa-signature", "")
            if not signature:
                raise HTTPException(status_code=401, detail="Missing signature")

        if body.get("Body", {}).get("stkCallback"):
            result = webhook_manager.parse_stk_callback(body)
            webhook_manager.emit("stk:callback", result)
        elif body.get("Result", {}).get("ResultParameters", {}).get("ResultParameter"):
            params = body["Result"]["ResultParameters"]["ResultParameter"]
            has_balance = any(p.get("Key") == "AccountBalance" for p in params)
            has_status = any(p.get("Key") == "TransactionStatus" for p in params)

            if has_balance:
                webhook_manager.emit("account:balance", body)
            elif has_status:
                webhook_manager.emit("transaction:status", body)
            else:
                webhook_manager.emit("b2c:result", body)
        elif body.get("TransactionType"):
            webhook_manager.emit("c2b:validation", body)
        else:
            raise HTTPException(status_code=400, detail="Unknown webhook event type")

        return {"received": True}

    return router
