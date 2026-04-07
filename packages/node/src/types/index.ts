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