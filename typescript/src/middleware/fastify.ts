import type { WebhookManager, WebhookEvent } from "../webhooks/index.js";

export interface FastifyMpesaWebhookOptions {
  webhookManager: WebhookManager;
  path?: string;
  verifySignature?: boolean;
  secret?: string;
}

export function createFastifyPlugin(
  options: FastifyMpesaWebhookOptions,
): (fastify: any) => Promise<void> {
  return async (fastify: any) => {
    const path = options.path ?? "/mpesa/webhook";

    fastify.post(path, async (request: any, reply: any) => {
      const body = request.body;

      if (options.verifySignature && options.secret) {
        const signature = request.headers["x-mpesa-signature"] as string;
        if (!signature) {
          return reply.status(401).send({ error: "Missing signature" });
        }
      }

      let event: WebhookEvent;

      if (body?.Body?.stkCallback) {
        event = { type: "stk:callback", payload: body };
      } else if (body?.Result?.ResultParameters?.ResultParameter) {
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
      } else if (body?.TransactionType) {
        event = { type: "c2b:validation", payload: body };
      } else {
        return reply.status(400).send({ error: "Unknown webhook event type" });
      }

      await options.webhookManager.handleEvent(event);
      return reply.status(200).send({ received: true });
    });
  };
}
