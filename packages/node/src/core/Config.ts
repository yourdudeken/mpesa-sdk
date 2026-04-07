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