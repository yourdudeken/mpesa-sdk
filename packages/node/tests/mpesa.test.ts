import { Mpesa, MpesaConfig } from '../src/index';

describe('Mpesa Node.js SDK', () => {
  const config: MpesaConfig = {
    environment: 'sandbox',
    mpesaConsumerKey: 'test_key',
    mpesaConsumerSecret: 'test_secret',
    passkey: 'test_passkey',
    shortcode: '174379',
    initiatorName: 'testapi',
    initiatorPassword: 'test_password',
    callbacks: {
      callbackUrl: 'https://test.com/callback'
    }
  };

  it('should create Mpesa instance', () => {
    const mpesa = new Mpesa(config);
    expect(mpesa).toBeDefined();
  });

  it('should have static constants', () => {
    expect(Mpesa.PAYBILL).toBe('CustomerPayBillOnline');
    expect(Mpesa.TILL).toBe('CustomerBuyGoodsOnline');
  });
});
