import axios from 'axios';
import { Mpesa, MpesaConfig } from '../src/index';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('Mpesa Node.js SDK', () => {
  const config: MpesaConfig = {
    environment: 'sandbox',
    mpesaConsumerKey: 'test_key',
    mpesaConsumerSecret: 'test_secret',
    passkey: 'test_passkey',
    shortcode: '174379',
    initiatorName: 'testapi',
    initiatorPassword: 'test_password',
    b2cShortcode: '600000',
    callbacks: {
      callbackUrl: 'https://test.com/callback',
      b2cResultUrl: 'https://test.com/b2c_result',
      b2cTimeoutUrl: 'https://test.com/b2c_timeout',
      b2bResultUrl: 'https://test.com/b2b_result',
      b2bTimeoutUrl: 'https://test.com/b2b_timeout',
      b2pochiResultUrl: 'https://test.com/b2pochi_result',
      b2pochiTimeoutUrl: 'https://test.com/b2pochi_timeout',
      c2bValidationUrl: 'https://test.com/c2b_validate',
      c2bConfirmationUrl: 'https://test.com/c2b_confirm',
      balanceResultUrl: 'https://test.com/balance_result',
      balanceTimeoutUrl: 'https://test.com/balance_timeout',
      statusResultUrl: 'https://test.com/status_result',
      statusTimeoutUrl: 'https://test.com/status_timeout',
      reversalResultUrl: 'https://test.com/reversal_result',
      reversalTimeoutUrl: 'https://test.com/reversal_timeout',
    }
  };

  describe('Configuration', () => {
    it('should create Mpesa instance', () => {
      const mpesa = new Mpesa(config);
      expect(mpesa).toBeDefined();
    });

    it('should have static constants', () => {
      expect(Mpesa.PAYBILL).toBe('CustomerPayBillOnline');
      expect(Mpesa.TILL).toBe('CustomerBuyGoodsOnline');
    });
  });

  describe('Token Generation', () => {
    it('should generate access token', async () => {
      const expectedResponse = {
        access_token: 'test_token_123',
        expires_in: '3599'
      };

      mockedAxios.create.mockReturnValueOnce({
        get: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        post: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const token = await (mpesa as any).generateAccessToken('C2B');
      expect(token).toBe('test_token_123');
    });

    it('should use B2C credentials for B2C requests', async () => {
      const b2cConfig: MpesaConfig = {
        ...config,
        b2cConsumerKey: 'b2c_key',
        b2cConsumerSecret: 'b2c_secret'
      };

      mockedAxios.create.mockReturnValueOnce({
        get: jest.fn().mockResolvedValueOnce({ data: { access_token: 'b2c_token', expires_in: '3599' } }),
        post: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(b2cConfig);
      const token = await (mpesa as any).generateAccessToken('B2C');
      expect(token).toBe('b2c_token');
    });
  });

  describe('STK Push', () => {
    it('should initiate STK push with callback URL', async () => {
      const expectedResponse = {
        MerchantRequestID: '29115-34620561-1',
        CheckoutRequestID: 'ws_CO_191220191020363925',
        ResponseCode: '0',
        ResponseDescription: 'Success',
        CustomerMessage: 'Success'
      };

      mockedAxios.create.mockReturnValueOnce({
        post: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        get: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const response = await mpesa.stkpush({
        phonenumber: '254712345678',
        amount: 100,
        accountNumber: '12345',
        callbackUrl: 'https://test.com/callback'
      });

      expect(response.ResponseCode).toBe('0');
    });

    it('should throw error when account reference is missing', async () => {
      const mpesa = new Mpesa(config);
      await expect(mpesa.stkpush({
        phonenumber: '254712345678',
        amount: 100,
        accountNumber: ''
      })).rejects.toThrow('Account Reference is required');
    });

    it('should require till number for TILL transactions', async () => {
      const mpesa = new Mpesa(config);
      await expect(mpesa.stkpush({
        phonenumber: '254712345678',
        amount: 100,
        accountNumber: '12345',
        transactionType: Mpesa.TILL
      })).rejects.toThrow('Till number is required');
    });

    it('should query STK push status', async () => {
      const expectedResponse = {
        ResponseCode: '0',
        ResponseDescription: 'Success',
        MerchantRequestID: '22205-34066-1',
        CheckoutRequestID: 'ws_CO_13012021093521236557',
        ResultCode: '0',
        ResultDesc: 'Success'
      };

      mockedAxios.create.mockReturnValueOnce({
        post: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        get: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const response = await mpesa.stkquery('ws_CO_191220191020363925');
      expect(response.ResponseCode).toBe('0');
    });
  });

  describe('B2C', () => {
    it('should send B2C payment', async () => {
      const expectedResponse = {
        ConversationID: 'AG_20231217_201020363925',
        OriginatorConversationID: '201020363925',
        ResponseCode: '0',
        ResponseDescription: 'Success'
      };

      mockedAxios.create.mockReturnValueOnce({
        post: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        get: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const response = await mpesa.b2c({
        phonenumber: '254712345678',
        commandId: 'BusinessPayment',
        amount: 100,
        remarks: 'Test payment'
      });

      expect(response.ResponseCode).toBe('0');
    });

    it('should send validated B2C payment', async () => {
      const expectedResponse = {
        ConversationID: 'AG_20231217_201020363925',
        ResponseCode: '0'
      };

      mockedAxios.create.mockReturnValueOnce({
        post: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        get: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const response = await mpesa.validated_b2c({
        phonenumber: '254712345678',
        commandId: 'BusinessPayment',
        amount: 100,
        remarks: 'Test payment',
        idNumber: '12345678'
      });

      expect(response.ResponseCode).toBe('0');
    });
  });

  describe('B2B', () => {
    it('should send B2B payment', async () => {
      const expectedResponse = {
        ConversationID: 'AG_20231217_201020363925',
        ResponseCode: '0'
      };

      mockedAxios.create.mockReturnValueOnce({
        post: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        get: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const response = await mpesa.b2b({
        receiverShortcode: '600000',
        commandId: 'BusinessPayBill',
        amount: 100,
        remarks: 'Test payment',
        accountNumber: '12345'
      });

      expect(response.ResponseCode).toBe('0');
    });

    it('should throw error when account number missing for BusinessPayBill', async () => {
      const mpesa = new Mpesa(config);
      await expect(mpesa.b2b({
        receiverShortcode: '600000',
        commandId: 'BusinessPayBill',
        amount: 100,
        remarks: 'Test payment'
      })).rejects.toThrow('Account Number is required');
    });
  });

  describe('C2B', () => {
    it('should register C2B URLs', async () => {
      const expectedResponse = {
        ResponseCode: '0',
        ResponseDescription: 'success'
      };

      mockedAxios.create.mockReturnValueOnce({
        post: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        get: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const response = await mpesa.c2bregisterURLS({
        shortcode: '600000'
      });

      expect(response.ResponseCode).toBe('0');
    });

    it('should simulate C2B payment', async () => {
      const expectedResponse = {
        ResponseCode: '0',
        ResponseDescription: 'Success'
      };

      mockedAxios.create.mockReturnValueOnce({
        post: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        get: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const response = await mpesa.c2bsimulate({
        phonenumber: '254712345678',
        amount: 100,
        shortcode: '600000',
        commandId: 'CustomerPayBillOnline'
      });

      expect(response.ResponseCode).toBe('0');
    });
  });

  describe('Account Balance', () => {
    it('should query account balance', async () => {
      const expectedResponse = {
        ResponseCode: '0',
        ResponseDescription: 'Success'
      };

      mockedAxios.create.mockReturnValueOnce({
        post: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        get: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const response = await mpesa.accountBalance({
        shortcode: '600000',
        identifierType: 4,
        remarks: 'Check balance'
      });

      expect(response.ResponseCode).toBe('0');
    });
  });

  describe('Transaction Status', () => {
    it('should query transaction status', async () => {
      const expectedResponse = {
        ResponseCode: '0',
        ResponseDescription: 'Success'
      };

      mockedAxios.create.mockReturnValueOnce({
        post: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        get: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const response = await mpesa.transactionStatus({
        shortcode: '600000',
        transactionId: '123456789',
        identifierType: 1,
        remarks: 'Check status'
      });

      expect(response.ResponseCode).toBe('0');
    });
  });

  describe('Reversal', () => {
    it('should reverse a transaction', async () => {
      const expectedResponse = {
        ResponseCode: '0',
        ResponseDescription: 'Success'
      };

      mockedAxios.create.mockReturnValueOnce({
        post: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        get: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const response = await mpesa.reversal({
        shortcode: '600000',
        transactionId: '123456789',
        amount: 100,
        remarks: 'Reverse transaction'
      });

      expect(response.ResponseCode).toBe('0');
    });
  });

  describe('B2Pochi', () => {
    it('should send B2Pochi payment', async () => {
      const expectedResponse = {
        ConversationID: 'AG_20231217_201020363925',
        ResponseCode: '0'
      };

      mockedAxios.create.mockReturnValueOnce({
        post: jest.fn().mockResolvedValueOnce({ data: expectedResponse }),
        get: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const response = await mpesa.b2pochi({
        phonenumber: '254712345678',
        amount: 100,
        remarks: 'Pochi payment'
      });

      expect(response.ResponseCode).toBe('0');
    });
  });

  describe('Phone Validation', () => {
    it('should validate phone numbers correctly', async () => {
      mockedAxios.create.mockReturnValueOnce({
        get: jest.fn().mockResolvedValueOnce({ data: { access_token: 'test', expires_in: '3599' } }),
        post: jest.fn(),
        defaults: {},
        interceptors: { request: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() }, response: { use: jest.fn(), eject: jest.fn(), clear: jest.fn() } },
        getUri: jest.fn(),
        request: jest.fn(),
        delete: jest.fn(),
        head: jest.fn(),
        options: jest.fn(),
        patch: jest.fn(),
        put: jest.fn()
      } as any);

      const mpesa = new Mpesa(config);
      const phoneValidator = (mpesa as any).phoneValidator.bind(mpesa);

      expect(phoneValidator('+254712345678')).toBe('254712345678');
      expect(phoneValidator('0712345678')).toBe('254712345678');
      expect(phoneValidator('712345678')).toBe('254712345678');
    });
  });
});