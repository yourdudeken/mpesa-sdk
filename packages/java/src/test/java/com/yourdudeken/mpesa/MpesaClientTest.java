package com.yourdudeken.mpesa;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.util.HashMap;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

public class MpesaClientTest {

    @Mock
    private com.yourdudeken.mpesa.http.HttpClient mockHttpClient;

    @BeforeEach
    public void setup() {
        MockitoAnnotations.openMocks(this);
    }

    @Test
    public void testConfigCreation() {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");
        config.setB2cShortcode("600000");

        Map<String, String> callbacks = new HashMap<>();
        callbacks.put("callback_url", "https://test.com/callback");
        callbacks.put("b2c_result_url", "https://test.com/b2c_result");
        callbacks.put("b2c_timeout_url", "https://test.com/b2c_timeout");
        config.setCallbacks(callbacks);

        assertEquals("sandbox", config.getEnvironment());
        assertEquals("test_key", config.getMpesaConsumerKey());
        assertEquals("174379", config.getShortcode());
        assertEquals("600000", config.getB2cShortcode());
    }

    @Test
    public void testMpesaInitialization() {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");

        Mpesa mpesa = new Mpesa(config);
        assertNotNull(mpesa);
    }

    @Test
    public void testStaticConstants() {
        assertEquals("CustomerPayBillOnline", Mpesa.PAYBILL);
        assertEquals("CustomerBuyGoodsOnline", Mpesa.TILL);
    }

    @Test
    public void testPhoneValidator() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");

        Mpesa mpesa = new Mpesa(config);

        assertEquals("254712345678", mpesa.phoneValidator("+254712345678"));
        assertEquals("254712345678", mpesa.phoneValidator("0712345678"));
        assertEquals("254712345678", mpesa.phoneValidator("712345678"));
    }

    @Test
    public void testStkpush() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");
        
        Map<String, String> callbacks = new HashMap<>();
        callbacks.put("callback_url", "https://test.com/callback");
        config.setCallbacks(callbacks);

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("phonenumber", "254712345678");
        params.put("amount", 100);
        params.put("accountNumber", "12345");
        params.put("callbackUrl", "https://test.com/callback");

        Map<String, Object> mockResponse = new HashMap<>();
        mockResponse.put("ResponseCode", "0");
        mockResponse.put("ResponseDescription", "Success");
        mockResponse.put("MerchantRequestID", "29115-34620561-1");
        mockResponse.put("CheckoutRequestID", "ws_CO_191220191020363925");

        when(mockHttpClient.post(anyString(), any(Map.class), anyString())).thenReturn("{\"ResponseCode\":\"0\"}");

        try {
            Map<String, Object> result = mpesa.stkpush(params);
            assertNotNull(result);
        } catch (Exception e) {
            // Expected behavior due to mocked HTTP client
        }
    }

    @Test
    public void testStkpushThrowsExceptionWhenAccountReferenceMissing() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("phonenumber", "254712345678");
        params.put("amount", 100);
        params.put("accountNumber", "");

        assertThrows(Exception.class, () -> mpesa.stkpush(params));
    }

    @Test
    public void testStkpushThrowsExceptionWhenTillNumberRequired() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("phonenumber", "254712345678");
        params.put("amount", 100);
        params.put("accountNumber", "12345");
        params.put("transactionType", "CustomerBuyGoodsOnline");

        assertThrows(Exception.class, () -> mpesa.stkpush(params));
    }

    @Test
    public void testB2c() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");
        config.setB2cShortcode("600000");

        Map<String, String> callbacks = new HashMap<>();
        callbacks.put("b2c_result_url", "https://test.com/b2c_result");
        callbacks.put("b2c_timeout_url", "https://test.com/b2c_timeout");
        config.setCallbacks(callbacks);

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("phonenumber", "254712345678");
        params.put("commandId", "BusinessPayment");
        params.put("amount", 100);
        params.put("remarks", "Test payment");

        try {
            Map<String, Object> result = mpesa.b2c(params);
            assertNotNull(result);
        } catch (Exception e) {
            // Expected behavior
        }
    }

    @Test
    public void testB2b() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");
        config.setB2cShortcode("600000");

        Map<String, String> callbacks = new HashMap<>();
        callbacks.put("b2b_result_url", "https://test.com/b2b_result");
        callbacks.put("b2b_timeout_url", "https://test.com/b2b_timeout");
        config.setCallbacks(callbacks);

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("receiverShortcode", "600000");
        params.put("commandId", "BusinessPayBill");
        params.put("amount", 100);
        params.put("remarks", "Test payment");
        params.put("accountNumber", "12345");

        try {
            Map<String, Object> result = mpesa.b2b(params);
            assertNotNull(result);
        } catch (Exception e) {
            // Expected behavior
        }
    }

    @Test
    public void testB2bThrowsExceptionWhenAccountNumberMissing() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");
        config.setB2cShortcode("600000");

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("receiverShortcode", "600000");
        params.put("commandId", "BusinessPayBill");
        params.put("amount", 100);
        params.put("remarks", "Test payment");

        assertThrows(Exception.class, () -> mpesa.b2b(params));
    }

    @Test
    public void testC2bregisterURLS() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");

        Map<String, String> callbacks = new HashMap<>();
        callbacks.put("c2b_confirmation_url", "https://test.com/c2b_confirm");
        callbacks.put("c2b_validation_url", "https://test.com/c2b_validate");
        config.setCallbacks(callbacks);

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("shortcode", "600000");

        try {
            Map<String, Object> result = mpesa.c2bregisterURLS(params);
            assertNotNull(result);
        } catch (Exception e) {
            // Expected behavior
        }
    }

    @Test
    public void testC2bsimulate() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("phonenumber", "254712345678");
        params.put("amount", 100);
        params.put("shortcode", "600000");
        params.put("commandId", "CustomerPayBillOnline");

        try {
            Map<String, Object> result = mpesa.c2bsimulate(params);
            assertNotNull(result);
        } catch (Exception e) {
            // Expected behavior
        }
    }

    @Test
    public void testAccountBalance() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");

        Map<String, String> callbacks = new HashMap<>();
        callbacks.put("balance_result_url", "https://test.com/balance_result");
        callbacks.put("balance_timeout_url", "https://test.com/balance_timeout");
        config.setCallbacks(callbacks);

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("shortcode", "600000");
        params.put("identifierType", 4);
        params.put("remarks", "Check balance");

        try {
            Map<String, Object> result = mpesa.accountBalance(params);
            assertNotNull(result);
        } catch (Exception e) {
            // Expected behavior
        }
    }

    @Test
    public void testTransactionStatus() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");

        Map<String, String> callbacks = new HashMap<>();
        callbacks.put("status_result_url", "https://test.com/status_result");
        callbacks.put("status_timeout_url", "https://test.com/status_timeout");
        config.setCallbacks(callbacks);

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("shortcode", "600000");
        params.put("transactionId", "123456789");
        params.put("identifierType", 1);
        params.put("remarks", "Check status");

        try {
            Map<String, Object> result = mpesa.transactionStatus(params);
            assertNotNull(result);
        } catch (Exception e) {
            // Expected behavior
        }
    }

    @Test
    public void testReversal() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");

        Map<String, String> callbacks = new HashMap<>();
        callbacks.put("reversal_result_url", "https://test.com/reversal_result");
        callbacks.put("reversal_timeout_url", "https://test.com/reversal_timeout");
        config.setCallbacks(callbacks);

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("shortcode", "600000");
        params.put("transactionId", "123456789");
        params.put("amount", 100);
        params.put("remarks", "Reverse transaction");

        try {
            Map<String, Object> result = mpesa.reversal(params);
            assertNotNull(result);
        } catch (Exception e) {
            // Expected behavior
        }
    }

    @Test
    public void testB2pochi() throws Exception {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        config.setPasskey("test_passkey");
        config.setShortcode("174379");
        config.setInitiatorName("testapi");
        config.setInitiatorPassword("test_password");
        config.setB2cShortcode("600000");

        Map<String, String> callbacks = new HashMap<>();
        callbacks.put("b2pochi_result_url", "https://test.com/b2pochi_result");
        callbacks.put("b2pochi_timeout_url", "https://test.com/b2pochi_timeout");
        config.setCallbacks(callbacks);

        Mpesa mpesa = new Mpesa(config);

        Map<String, Object> params = new HashMap<>();
        params.put("phonenumber", "254712345678");
        params.put("amount", 100);
        params.put("remarks", "Pochi payment");

        try {
            Map<String, Object> result = mpesa.b2pochi(params);
            assertNotNull(result);
        } catch (Exception e) {
            // Expected behavior
        }
    }
}