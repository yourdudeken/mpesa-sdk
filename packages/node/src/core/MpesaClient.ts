import axios, { AxiosInstance } from 'axios';
import * as crypto from 'crypto';
import * as fs from 'fs';
import * as path from 'path';

export interface MpesaConfig {
  environment: 'sandbox' | 'production';
  mpesaConsumerKey: string;
  mpesaConsumerSecret: string;
  b2cConsumerKey?: string;
  b2cConsumerSecret?: string;
  passkey: string;
  shortcode: string;
  tillNumber?: string;
  initiatorName: string;
  initiatorPassword: string;
  b2cShortcode?: string;
  callbacks: {
    c2bValidationUrl?: string;
    c2bConfirmationUrl?: string;
    b2cResultUrl?: string;
    b2cTimeoutUrl?: string;
    b2pochiResultUrl?: string;
    b2pochiTimeoutUrl?: string;
    callbackUrl?: string;
    statusResultUrl?: string;
    statusTimeoutUrl?: string;
    balanceResultUrl?: string;
    balanceTimeoutUrl?: string;
    reversalResultUrl?: string;
    reversalTimeoutUrl?: string;
    b2bResultUrl?: string;
    b2bTimeoutUrl?: string;
  };
}

export interface StkPushParams {
  phonenumber: string;
  amount: number;
  accountNumber: string;
  callbackUrl?: string;
  transactionType?: 'CustomerPayBillOnline' | 'CustomerBuyGoodsOnline';
  shortCodeType?: 'C2B' | 'B2C' | 'B2B';
}

export interface B2cParams {
  phonenumber: string;
  commandId: 'SalaryPayment' | 'BusinessPayment' | 'PromotionPayment';
  amount: number;
  remarks: string;
  resultUrl?: string;
  timeoutUrl?: string;
  shortCodeType?: 'C2B' | 'B2C' | 'B2B';
}

export interface B2bParams {
  receiverShortcode: string;
  commandId: 'BusinessPayBill' | 'MerchantToMerchantTransfer' | 'MerchantTransferFromMerchantToWorking' | 'MerchantServicesMMFAccountTransfer' | 'AgencyFloatAdvance';
  amount: number;
  remarks: string;
  accountNumber?: string;
  resultUrl?: string;
  timeoutUrl?: string;
  shortCodeType?: 'C2B' | 'B2C' | 'B2B';
}

export interface C2bRegisterParams {
  shortcode: string;
  confirmUrl?: string;
  validateUrl?: string;
  shortCodeType?: 'C2B' | 'B2C' | 'B2B';
}

export interface C2bSimulateParams {
  phonenumber: string;
  amount: number;
  shortcode: string;
  commandId: 'CustomerPayBillOnline' | 'CustomerBuyGoodsOnline';
  accountNumber?: string;
  shortCodeType?: 'C2B' | 'B2C' | 'B2B';
}

export interface TransactionStatusParams {
  shortcode: string;
  transactionId: string;
  identifierType: 1 | 2 | 4;
  remarks: string;
  resultUrl?: string;
  timeoutUrl?: string;
  shortCodeType?: 'C2B' | 'B2C' | 'B2B';
}

export interface AccountBalanceParams {
  shortcode: string;
  identifierType: 1 | 2 | 4;
  remarks: string;
  resultUrl?: string;
  timeoutUrl?: string;
  shortCodeType?: 'C2B' | 'B2C' | 'B2B';
}

export interface ReversalParams {
  shortcode: string;
  transactionId: string;
  amount: number;
  remarks: string;
  resultUrl?: string;
  timeoutUrl?: string;
  shortCodeType?: 'C2B' | 'B2C' | 'B2B';
}

export interface B2PochiParams {
  phonenumber: string;
  amount: number;
  remarks: string;
  occasion?: string;
  resultUrl?: string;
  timeoutUrl?: string;
  shortCodeType?: 'C2B' | 'B2C' | 'B2B';
}

export class Mpesa {
  private config: MpesaConfig;
  private baseUrl: string;
  private client: AxiosInstance;
  private accessToken: string | null = null;
  private tokenExpiry: number = 0;

  public static PAYBILL = 'CustomerPayBillOnline';
  public static TILL = 'CustomerBuyGoodsOnline';

  constructor(config: MpesaConfig) {
    this.config = config;
    this.baseUrl = config.environment === 'sandbox'
      ? 'https://sandbox.safaricom.co.ke'
      : 'https://api.safaricom.co.ke';
    
    this.client = axios.create({
      baseURL: this.baseUrl,
    });
  }

  private getConfig(key: string): any {
    const keys = key.split('.');
    let value: any = this.config;
    for (const k of keys) {
      value = value?.[k];
    }
    return value;
  }

  private resolveCallbackUrl(paramUrl: string | undefined, configKey: string): string {
    const configUrl = this.getConfig(`callbacks.${configKey}`);
    if (paramUrl) return paramUrl;
    if (configUrl) return configUrl;
    throw new Error(`Ensure you have set the ${configKey} in the config or passed as a parameter`);
  }

  private phoneValidator(phoneNo: string): string {
    phoneNo = phoneNo.startsWith('+') ? phoneNo.substring(1) : phoneNo;
    phoneNo = phoneNo.startsWith('0') ? '254' + phoneNo.substring(1) : phoneNo;
    phoneNo = phoneNo.startsWith('7') ? '254' + phoneNo : phoneNo;
    return phoneNo;
  }

  private getFormattedTimestamp(): string {
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0');
    const day = String(now.getDate()).padStart(2, '0');
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    const seconds = String(now.getSeconds()).padStart(2, '0');
    return `${year}${month}${day}${hours}${minutes}${seconds}`;
  }

  private lipaNaMpesaPassword(): string {
    const timestamp = this.getFormattedTimestamp();
    const password = this.getConfig('shortcode') + this.getConfig('passkey') + timestamp;
    return Buffer.from(password).toString('base64');
  }

  private async generateAccessToken(shortCodeType: string = 'C2B'): Promise<string> {
    if (this.accessToken && Date.now() < this.tokenExpiry) {
      return this.accessToken;
    }

    const consumerKey = (shortCodeType === 'B2C' || shortCodeType === 'B2B')
      ? this.getConfig('b2cConsumerKey')
      : this.getConfig('mpesaConsumerKey');
    const consumerSecret = (shortCodeType === 'B2C' || shortCodeType === 'B2B')
      ? this.getConfig('b2cConsumerSecret')
      : this.getConfig('mpesaConsumerSecret');

    const auth = Buffer.from(`${consumerKey}:${consumerSecret}`).toString('base64');

    try {
      const response = await this.client.get('/oauth/v1/generate?grant_type=client_credentials', {
        headers: { Authorization: `Basic ${auth}` },
      });
      
      this.accessToken = response.data.access_token;
      this.tokenExpiry = Date.now() + (response.data.expires_in - 60) * 1000;
      return this.accessToken;
    } catch (error: any) {
      throw new Error(`Failed to generate access token: ${error.message}`);
    }
  }

  private generateSecurityCredential(): string {
    const certPath = this.config.environment === 'sandbox'
      ? path.join(__dirname, '../certificates/SandboxCertificate.cer')
      : path.join(__dirname, '../certificates/ProductionCertificate.cer');
    
    const pubkey = fs.readFileSync(certPath);
    const password = this.getConfig('initiatorPassword');
    
    const encrypted = crypto.publicEncrypt(
      { key: pubkey, padding: crypto.constants.RSA_PKCS1_PADDING },
      Buffer.from(password)
    );
    
    return encrypted.toString('base64');
  }

  private async mpesaRequest(url: string, body: any, shortCodeType: string = 'C2B'): Promise<any> {
    const token = await this.generateAccessToken(shortCodeType);
    
    try {
      const response = await this.client.post(url, body, {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });
      return response.data;
    } catch (error: any) {
      throw new Error(`Mpesa request failed: ${error.response?.data?.errorMessage || error.message}`);
    }
  }

  public async stkpush(params: StkPushParams): Promise<any> {
    const { phonenumber, amount, accountNumber, callbackUrl, transactionType = Mpesa.PAYBILL, shortCodeType = 'C2B' } = params;

    if (!accountNumber) {
      throw new Error('An Account Reference is required for All transactions.');
    }

    if (transactionType === Mpesa.TILL && !this.getConfig('tillNumber')) {
      throw new Error('Till number is required for Buy Goods transactions.');
    }

    const url = `${this.baseUrl}/mpesa/stkpush/v1/processrequest`;
    const data = {
      BusinessShortCode: this.getConfig('shortcode'),
      Password: this.lipaNaMpesaPassword(),
      Timestamp: this.getFormattedTimestamp(),
      Amount: amount,
      PartyA: this.phoneValidator(phonenumber),
      PartyB: transactionType === Mpesa.PAYBILL ? this.getConfig('shortcode') : this.getConfig('tillNumber'),
      TransactionType: transactionType,
      PhoneNumber: this.phoneValidator(phonenumber),
      TransactionDesc: 'Payment',
      AccountReference: accountNumber,
      CallBackURL: this.resolveCallbackUrl(callbackUrl, 'callbackUrl'),
    };

    return this.mpesaRequest(url, data, shortCodeType);
  }

  public async stkquery(checkoutRequestId: string, shortCodeType: string = 'C2B'): Promise<any> {
    const url = `${this.baseUrl}/mpesa/stkpushquery/v1/query`;
    const data = {
      BusinessShortCode: this.getConfig('shortcode'),
      Password: this.lipaNaMpesaPassword(),
      Timestamp: this.getFormattedTimestamp(),
      CheckoutRequestID: checkoutRequestId,
    };

    return this.mpesaRequest(url, data, shortCodeType);
  }

  public async b2c(params: B2cParams): Promise<any> {
    const url = `${this.baseUrl}/mpesa/b2c/v1/paymentrequest`;
    const body = {
      InitiatorName: this.getConfig('initiatorName'),
      SecurityCredential: this.generateSecurityCredential(),
      CommandID: params.commandId,
      Amount: params.amount,
      PartyA: this.getConfig('b2cShortcode'),
      PartyB: this.phoneValidator(params.phonenumber),
      Remarks: params.remarks,
      Occassion: '',
      ResultURL: this.resolveCallbackUrl(params.resultUrl, 'b2cResultUrl'),
      QueueTimeOutURL: this.resolveCallbackUrl(params.timeoutUrl, 'b2cTimeoutUrl'),
    };

    return this.mpesaRequest(url, body, params.shortCodeType || 'B2C');
  }

  public async validated_b2c(params: B2cParams & { idNumber: string }): Promise<any> {
    const url = `${this.baseUrl}/mpesa/b2cvalidate/v2/paymentrequest`;
    const body = {
      InitiatorName: this.getConfig('initiatorName'),
      SecurityCredential: this.generateSecurityCredential(),
      CommandID: params.commandId,
      Amount: params.amount,
      PartyA: this.getConfig('b2cShortcode'),
      PartyB: this.phoneValidator(params.phonenumber),
      Remarks: params.remarks,
      Occassion: '',
      OriginatorConversationID: this.getFormattedTimestamp(),
      IDType: '01',
      IDNumber: params.idNumber,
      ResultURL: this.resolveCallbackUrl(params.resultUrl, 'b2cResultUrl'),
      QueueTimeOutURL: this.resolveCallbackUrl(params.timeoutUrl, 'b2cTimeoutUrl'),
    };

    return this.mpesaRequest(url, body, params.shortCodeType || 'B2C');
  }

  public async b2b(params: B2bParams): Promise<any> {
    if (params.commandId === 'BusinessPayBill' && !params.accountNumber) {
      throw new Error('Account Number is required for BusinessPayBill CommandID');
    }

    const url = `${this.baseUrl}/mpesa/b2b/v1/paymentrequest`;
    const body = {
      Initiator: this.getConfig('initiatorName'),
      SecurityCredential: this.generateSecurityCredential(),
      CommandID: params.commandId,
      SenderIdentifierType: '4',
      RecieverIdentifierType: '4',
      Amount: params.amount,
      PartyA: this.getConfig('b2cShortcode'),
      PartyB: params.receiverShortcode,
      AccountReference: params.accountNumber,
      Remarks: params.remarks,
      ResultURL: this.resolveCallbackUrl(params.resultUrl, 'b2bResultUrl'),
      QueueTimeOutURL: this.resolveCallbackUrl(params.timeoutUrl, 'b2bTimeoutUrl'),
    };

    return this.mpesaRequest(url, body, params.shortCodeType || 'B2B');
  }

  public async c2bregisterURLS(params: C2bRegisterParams): Promise<any> {
    const url = `${this.baseUrl}/mpesa/c2b/v2/registerurl`;
    const body = {
      ShortCode: params.shortcode,
      ResponseType: 'Completed',
      ConfirmationURL: this.resolveCallbackUrl(params.confirmUrl, 'c2bConfirmationUrl'),
      ValidationURL: this.resolveCallbackUrl(params.validateUrl, 'c2bValidationUrl'),
    };

    return this.mpesaRequest(url, body, params.shortCodeType || 'C2B');
  }

  public async c2bsimulate(params: C2bSimulateParams): Promise<any> {
    const url = `${this.baseUrl}/mpesa/c2b/v2/simulate`;
    const data = params.commandId === Mpesa.PAYBILL ? {
      Msisdn: this.phoneValidator(params.phonenumber),
      Amount: params.amount,
      BillRefNumber: params.accountNumber,
      CommandID: params.commandId,
      ShortCode: params.shortcode,
    } : {
      Msisdn: this.phoneValidator(params.phonenumber),
      Amount: params.amount,
      CommandID: params.commandId,
      ShortCode: params.shortcode,
    };

    return this.mpesaRequest(url, data, params.shortCodeType || 'C2B');
  }

  public async transactionStatus(params: TransactionStatusParams): Promise<any> {
    const url = `${this.baseUrl}/mpesa/transactionstatus/v1/query`;
    const body = {
      Initiator: this.getConfig('initiatorName'),
      SecurityCredential: this.generateSecurityCredential(),
      CommandID: 'TransactionStatusQuery',
      TransactionID: params.transactionId,
      PartyA: params.shortcode,
      IdentifierType: params.identifierType,
      Remarks: params.remarks,
      Occassion: '',
      ResultURL: this.resolveCallbackUrl(params.resultUrl, 'statusResultUrl'),
      QueueTimeOutURL: this.resolveCallbackUrl(params.timeoutUrl, 'statusTimeoutUrl'),
    };

    return this.mpesaRequest(url, body, params.shortCodeType || 'C2B');
  }

  public async accountBalance(params: AccountBalanceParams): Promise<any> {
    const url = `${this.baseUrl}/mpesa/accountbalance/v1/query`;
    const body = {
      Initiator: this.getConfig('initiatorName'),
      SecurityCredential: this.generateSecurityCredential(),
      CommandID: 'AccountBalance',
      PartyA: params.shortcode,
      IdentifierType: params.identifierType,
      Remarks: params.remarks,
      ResultURL: this.resolveCallbackUrl(params.resultUrl, 'balanceResultUrl'),
      QueueTimeOutURL: this.resolveCallbackUrl(params.timeoutUrl, 'balanceTimeoutUrl'),
    };

    return this.mpesaRequest(url, body, params.shortCodeType || 'C2B');
  }

  public async reversal(params: ReversalParams): Promise<any> {
    const url = `${this.baseUrl}/mpesa/reversal/v1/request`;
    const body = {
      Initiator: this.getConfig('initiatorName'),
      SecurityCredential: this.generateSecurityCredential(),
      CommandID: 'TransactionReversal',
      TransactionID: params.transactionId,
      Amount: params.amount,
      ReceiverParty: params.shortcode,
      RecieverIdentifierType: '11',
      Remarks: params.remarks,
      Occasion: '',
      ResultURL: this.resolveCallbackUrl(params.resultUrl, 'reversalResultUrl'),
      QueueTimeOutURL: this.resolveCallbackUrl(params.timeoutUrl, 'reversalTimeoutUrl'),
    };

    return this.mpesaRequest(url, body, params.shortCodeType || 'C2B');
  }

  public async b2pochi(params: B2PochiParams): Promise<any> {
    const url = `${this.baseUrl}/mpesa/b2pochi/v1/paymentrequest`;
    const body = {
      OriginatorConversationID: this.getFormattedTimestamp(),
      InitiatorName: this.getConfig('initiatorName'),
      SecurityCredential: this.generateSecurityCredential(),
      CommandID: 'BusinessPayToPochi',
      Amount: params.amount,
      PartyA: this.getConfig('b2cShortcode'),
      PartyB: this.phoneValidator(params.phonenumber),
      Remarks: params.remarks,
      Occassion: params.occasion || '',
      ResultURL: this.resolveCallbackUrl(params.resultUrl, 'b2pochiResultUrl'),
      QueueTimeOutURL: this.resolveCallbackUrl(params.timeoutUrl, 'b2pochiTimeoutUrl'),
    };

    return this.mpesaRequest(url, body, params.shortCodeType || 'B2C');
  }
}

export default Mpesa;