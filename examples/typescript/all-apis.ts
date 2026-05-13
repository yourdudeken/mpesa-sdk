import { Mpesa } from "@yourdudeken/mpesa-sdk";

function resolveShortcode(): number {
  const raw = process.env.MPESA_SHORTCODE;
  if (!raw) throw new Error("MPESA_SHORTCODE is required");
  const parsed = parseInt(raw, 10);
  if (isNaN(parsed) || parsed <= 0) {
    throw new Error(`Invalid MPESA_SHORTCODE: ${raw}`);
  }
  return parsed;
}

const mpesa = new Mpesa({
  consumerKey: process.env.MPESA_CONSUMER_KEY!,
  consumerSecret: process.env.MPESA_CONSUMER_SECRET!,
  environment: (process.env.MPESA_ENV as "sandbox" | "production") ?? "sandbox",
  passkey: process.env.MPESA_PASSKEY!,
  initiatorName: process.env.MPESA_INITIATOR_NAME!,
  initiatorPassword: process.env.MPESA_INITIATOR_PASSWORD!,
  securityCredential: process.env.MPESA_SECURITY_CREDENTIAL!,
});

export async function stkPush() {
  const shortcode = resolveShortcode();
  const response = await mpesa.stkPush.initiate({
    BusinessShortCode: shortcode,
    TransactionType: "CustomerPayBillOnline",
    Amount: 1,
    PartyA: 254722000000,
    PartyB: shortcode,
    PhoneNumber: 254722000000,
    CallBackURL: "https://your-domain.com/api/mpesa/callback",
    AccountReference: "INV-001",
    TransactionDesc: "Payment for invoice 001",
    Password: "",
    Timestamp: "",
  });
  return response;
}

export async function stkQuery(checkoutRequestId: string) {
  const response = await mpesa.stkPush.query({
    BusinessShortCode: String(resolveShortcode()),
    CheckoutRequestID: checkoutRequestId,
    Password: "",
    Timestamp: "",
  });
  return response;
}

export async function c2bRegisterURL() {
  const response = await mpesa.c2b.registerURL({
    ShortCode: String(resolveShortcode()),
    ResponseType: "Completed",
    ConfirmationURL: "https://your-domain.com/api/c2b/confirmation",
    ValidationURL: "https://your-domain.com/api/c2b/validation",
  });
  return response;
}

export async function c2bSimulate() {
  const response = await mpesa.c2b.simulate({
    ShortCode: resolveShortcode(),
    CommandID: "CustomerPayBillOnline",
    Amount: 100,
    Msisdn: 254708374149,
    BillRefNumber: "ACCNO-001",
  });
  return response;
}

export async function b2cPayment() {
  const response = await mpesa.b2c.send({
    InitiatorName: process.env.MPESA_INITIATOR_NAME!,
    SecurityCredential: process.env.MPESA_SECURITY_CREDENTIAL!,
    CommandID: "BusinessPayment",
    Amount: 100,
    PartyA: resolveShortcode(),
    PartyB: 254705912645,
    Remarks: "Salary disbursement",
    QueueTimeOutURL: "https://your-domain.com/api/b2c/queue",
    ResultURL: "https://your-domain.com/api/b2c/result",
    Occassion: "Monthly Salary",
  });
  return response;
}

export async function b2bPayment() {
  const response = await mpesa.b2b.send({
    Initiator: process.env.MPESA_INITIATOR_NAME!,
    SecurityCredential: process.env.MPESA_SECURITY_CREDENTIAL!,
    CommandID: "BusinessPayBill",
    Amount: 5000,
    PartyA: 123456,
    PartyB: 654321,
    Remarks: "Supplier payment",
    QueueTimeOutURL: "https://your-domain.com/api/b2b/queue",
    ResultURL: "https://your-domain.com/api/b2b/result",
    AccountReference: "SUPP-001",
  });
  return response;
}

export async function reverseTransaction(transactionId: string) {
  const response = await mpesa.reversal.reverse({
    Initiator: process.env.MPESA_INITIATOR_NAME!,
    SecurityCredential: process.env.MPESA_SECURITY_CREDENTIAL!,
    CommandID: "TransactionReversal",
    TransactionID: transactionId,
    Amount: 100,
    ReceiverParty: resolveShortcode(),
    QueueTimeOutURL: "https://your-domain.com/api/reversal/queue",
    ResultURL: "https://your-domain.com/api/reversal/result",
    Remarks: "Customer initiated reversal",
  });
  return response;
}

export async function checkTransactionStatus(transactionId: string) {
  const response = await mpesa.transactionStatus.query({
    Initiator: process.env.MPESA_INITIATOR_NAME!,
    SecurityCredential: process.env.MPESA_SECURITY_CREDENTIAL!,
    CommandID: "TransactionStatusQuery",
    TransactionID: transactionId,
    PartyA: resolveShortcode(),
    IdentifierType: 4,
    ResultURL: "https://your-domain.com/api/status/result",
    QueueTimeOutURL: "https://your-domain.com/api/status/queue",
    Remarks: "Status check",
  });
  return response;
}

export async function checkAccountBalance() {
  const response = await mpesa.accountBalance.query({
    Initiator: process.env.MPESA_INITIATOR_NAME!,
    SecurityCredential: process.env.MPESA_SECURITY_CREDENTIAL!,
    CommandID: "AccountBalance",
    PartyA: resolveShortcode(),
    IdentifierType: 4,
    Remarks: "Daily balance check",
    QueueTimeOutURL: "https://your-domain.com/api/balance/queue",
    ResultURL: "https://your-domain.com/api/balance/result",
  });
  return response;
}

export async function generateQR() {
  const response = await mpesa.dynamicQR.generate({
    MerchantName: "Your Business Name",
    RefNo: "INV-2024-001",
    Amount: 1500,
    TrxCode: "BG",
    CPI: String(resolveShortcode()),
    Size: "300",
  });
  return response;
}
