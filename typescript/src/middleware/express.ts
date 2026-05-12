import type { WebhookManager, WebhookEvent } from "../webhooks/index.js";

export interface MpesaWebhookOptions {
  webhookManager: WebhookManager;
  path?: string;
  verifySignature?: boolean;
  secret?: string;
}

export function createExpressMiddleware(
  options: MpesaWebhookOptions,
) {
  return async (req: any, res: any, next: any): Promise<void> => {
    if (options.path && req.path !== options.path) {
      return next();
    }

    if (req.method !== "POST") {
      return next();
    }

    try {
      if (options.verifySignature && options.secret) {
        const signature = req.headers["x-mpesa-signature"] as string;
        if (!signature) {
          res.status(401).json({ error: "Missing signature" });
          return;
        }
      }

      const body = req.body;
      let event: WebhookEvent;

      if (body.Body?.stkCallback) {
        event = { type: "stk:callback", payload: body };
      } else if (body.Result?.ResultParameters?.ResultParameter) {
        const resultParams = body.Result.ResultParameters.ResultParameter;
        const hasAccountBalance = resultParams.some(
          (p: { Key: string }) => p.Key === "AccountBalance",
        );
        const hasTransactionStatus = resultParams.some(
          (p: { Key: string }) => p.Key === "TransactionStatus",
        );

        if (hasAccountBalance) {
          event = { type: "account:balance", payload: body };
        } else if (hasTransactionStatus) {
          event = { type: "transaction:status", payload: body };
        } else {
          event = { type: "b2c:result", payload: body };
        }
      } else if (body.TransactionType) {
        event = { type: "c2b:validation", payload: body };
      } else {
        res.status(400).json({ error: "Unknown webhook event type" });
        return;
      }

      await options.webhookManager.handleEvent(event);
      res.status(200).json({ received: true });
    } catch (error) {
      next(error);
    }
  };
}
