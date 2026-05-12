import { MpesaApiClient } from "./client/client.js";
import {
  STKPushService,
  C2BService,
  B2CService,
  B2BService,
  ReversalService,
  TransactionStatusService,
  AccountBalanceService,
  DynamicQRService,
} from "./services/index.js";
import { WebhookManager } from "./webhooks/index.js";
import type { MpesaConfig } from "./types/index.js";

export class Mpesa {
  public readonly stkPush: STKPushService;
  public readonly c2b: C2BService;
  public readonly b2c: B2CService;
  public readonly b2b: B2BService;
  public readonly reversal: ReversalService;
  public readonly transactionStatus: TransactionStatusService;
  public readonly accountBalance: AccountBalanceService;
  public readonly dynamicQR: DynamicQRService;
  public readonly webhooks: WebhookManager;
  public readonly client: MpesaApiClient;

  constructor(config: MpesaConfig) {
    this.client = new MpesaApiClient(config);
    this.stkPush = new STKPushService(this.client);
    this.c2b = new C2BService(this.client);
    this.b2c = new B2CService(this.client);
    this.b2b = new B2BService(this.client);
    this.reversal = new ReversalService(this.client);
    this.transactionStatus = new TransactionStatusService(this.client);
    this.accountBalance = new AccountBalanceService(this.client);
    this.dynamicQR = new DynamicQRService(this.client);
    this.webhooks = new WebhookManager({
      passkey: config.passkey,
    });
  }
}

export { MpesaApiClient } from "./client/client.js";
export * from "./services/index.js";
export * from "./types/index.js";
export * from "./errors/index.js";
export * from "./webhooks/index.js";
export * from "./middleware/index.js";
export * from "./utils/index.js";
