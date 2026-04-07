package com.yourdudeken.mpesa;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.util.EntityUtils;
import org.bouncycastle.crypto.engines.RSAEngine;
import org.bouncycastle.crypto.encoders.Base64;
import org.bouncycastle.crypto.encoders.PKCS1v15Encoder;
import org.bouncycastle.crypto.params.RSAKeyParameters;
import org.bouncycastle.util.encoders.DerEncoderException;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.security.Security;
import java.text.SimpleDateFormat;
import java.util.*;
import java.util.stream.Collectors;

public class Mpesa {
    public static final String PAYBILL = "CustomerPayBillOnline";
    public static final String TILL = "CustomerBuyGoodsOnline";

    private final MpesaConfig config;
    private final String baseUrl;
    private String accessToken;
    private long tokenExpiry;

    public Mpesa(MpesaConfig config) {
        this.config = config;
        this.baseUrl = config.getEnvironment().equals("sandbox")
            ? "https://sandbox.safaricom.co.ke"
            : "https://api.safaricom.co.ke";
    }

    private String getConfig(String key) {
        try {
            switch (key) {
                case "environment": return config.getEnvironment();
                case "mpesa_consumer_key": return config.getMpesaConsumerKey();
                case "mpesa_consumer_secret": return config.getMpesaConsumerSecret();
                case "b2c_consumer_key": return config.getB2cConsumerKey();
                case "b2c_consumer_secret": return config.getB2cConsumerSecret();
                case "passkey": return config.getPasskey();
                case "shortcode": return config.getShortcode();
                case "till_number": return config.getTillNumber();
                case "initiator_name": return config.getInitiatorName();
                case "initiator_password": return config.getInitiatorPassword();
                case "b2c_shortcode": return config.getB2cShortcode();
                default: return null;
            }
        } catch (Exception e) {
            return null;
        }
    }

    private String resolveCallbackUrl(String paramUrl, String configKey) throws Exception {
        String configUrl = config.getCallbacks() != null ? config.getCallbacks().get(configKey) : null;
        if (paramUrl != null && !paramUrl.isEmpty()) return paramUrl;
        if (configUrl != null && !configUrl.isEmpty()) return configUrl;
        throw new Exception("Ensure you have set the " + configKey + " in the config or passed as a parameter");
    }

    private String phoneValidator(String phoneNo) {
        if (phoneNo.startsWith("+")) phoneNo = phoneNo.substring(1);
        if (phoneNo.startsWith("0")) phoneNo = "254" + phoneNo.substring(1);
        else if (phoneNo.startsWith("7")) phoneNo = "254" + phoneNo;
        return phoneNo;
    }

    private String getFormattedTimestamp() {
        return new SimpleDateFormat("yyyyMMddHHmmss").format(new Date());
    }

    private String lipaNaMpesaPassword() throws Exception {
        String timestamp = getFormattedTimestamp();
        String password = getConfig("shortcode") + getConfig("passkey") + timestamp;
        return Base64.toBase64String(password.getBytes(StandardCharsets.UTF_8));
    }

    private String generateAccessToken(String shortCodeType) throws Exception {
        long currentTime = System.currentTimeMillis();
        if (accessToken != null && currentTime < tokenExpiry) {
            return accessToken;
        }

        String consumerKey = (shortCodeType.equals("B2C") || shortCodeType.equals("B2B"))
            ? getConfig("b2c_consumer_key")
            : getConfig("mpesa_consumer_key");
        String consumerSecret = (shortCodeType.equals("B2C") || shortCodeType.equals("B2B"))
            ? getConfig("b2c_consumer_secret")
            : getConfig("mpesa_consumer_secret");

        String auth = Base64.toBase64String((consumerKey + ":" + consumerSecret).getBytes(StandardCharsets.UTF_8));

        CloseableHttpClient client = HttpClients.createDefault();
        HttpGet request = new HttpGet(baseUrl + "/oauth/v1/generate?grant_type=client_credentials");
        request.setHeader("Authorization", "Basic " + auth);

        try (CloseableHttpResponse response = client.execute(request)) {
            String jsonResponse = EntityUtils.toString(response.getEntity());
            Map<String, Object> data = new Gson().fromJson(jsonResponse, Map.class);
            accessToken = (String) data.get("access_token");
            tokenExpiry = currentTime + ((Double) data.get("expires_in")).longValue() * 1000 - 60000;
            return accessToken;
        }
    }

    private String generateSecurityCredential() throws Exception {
        Security.addProvider(new org.bouncycastle.jce.provider.BouncyCastleProvider());
        
        String certPath = getConfig("environment").equals("sandbox")
            ? "certificates/SandboxCertificate.cer"
            : "certificates/ProductionCertificate.cer";
        
        InputStream is = Mpesa.class.getClassLoader().getResourceAsStream(certPath);
        if (is == null) {
            is = new FileInputStream(certPath);
        }
        
        BufferedReader reader = new BufferedReader(new InputStreamReader(is));
        StringBuilder pem = new StringBuilder();
        String line;
        while ((line = reader.readLine()) != null) {
            if (!line.startsWith("-----")) pem.append(line);
        }
        reader.close();
        
        byte[] certBytes = Base64.decode(pem.toString());
        
        java.security.cert.X509Certificate cert = java.security.cert.CertificateFactory.getInstance("X.509")
            .generateCertificate(new ByteArrayInputStream(certBytes));
        
        java.security.interfaces.RSAPublicKey rsaPublicKey = (java.security.interfaces.RSAPublicKey) cert.getPublicKey();
        
        RSAKeyParameters rsaKeyParams = new RSAKeyParameters(
            false,
            rsaPublicKey.getModulus(),
            rsaPublicKey.getPublicExponent()
        );
        
        RSAEngine engine = new RSAEngine();
        engine.init(true, rsaKeyParams);
        
        byte[] passwordBytes = getConfig("initiator_password").getBytes(StandardCharsets.UTF_8);
        byte[] encrypted = engine.processBlock(passwordBytes, 0, passwordBytes.length);
        
        return Base64.toBase64String(encrypted);
    }

    private Map<String, Object> mpesaRequest(String url, Map<String, Object> body, String shortCodeType) throws Exception {
        String token = generateAccessToken(shortCodeType);

        CloseableHttpClient client = HttpClients.createDefault();
        HttpPost request = new HttpPost(url);
        request.setHeader("Authorization", "Bearer " + token);
        request.setHeader("Content-Type", "application/json");
        
        Gson gson = new GsonBuilder().create();
        request.setEntity(new StringEntity(gson.toJson(body)));

        try (CloseableHttpResponse response = client.execute(request)) {
            String jsonResponse = EntityUtils.toString(response.getEntity());
            return new Gson().fromJson(jsonResponse, Map.class);
        }
    }

    public Map<String, Object> stkpush(Map<String, Object> params) throws Exception {
        String phonenumber = (String) params.get("phonenumber");
        Number amount = (Number) params.get("amount");
        String accountNumber = (String) params.get("accountNumber");
        String callbackUrl = (String) params.get("callbackUrl");
        String transactionType = params.containsKey("transactionType") ? (String) params.get("transactionType") : PAYBILL;
        String shortCodeType = params.containsKey("shortCodeType") ? (String) params.get("shortCodeType") : "C2B";

        if (accountNumber == null || accountNumber.isEmpty()) {
            throw new Exception("An Account Reference is required for All transactions.");
        }

        if (transactionType.equals(TILL) && (getConfig("till_number") == null || getConfig("till_number").isEmpty())) {
            throw new Exception("Till number is required for Buy Goods transactions.");
        }

        String url = baseUrl + "/mpesa/stkpush/v1/processrequest";
        Map<String, Object> data = new HashMap<>();
        data.put("BusinessShortCode", getConfig("shortcode"));
        data.put("Password", lipaNaMpesaPassword());
        data.put("Timestamp", getFormattedTimestamp());
        data.put("Amount", amount.intValue());
        data.put("PartyA", phoneValidator(phonenumber));
        data.put("PartyB", transactionType.equals(PAYBILL) ? getConfig("shortcode") : getConfig("till_number"));
        data.put("TransactionType", transactionType);
        data.put("PhoneNumber", phoneValidator(phonenumber));
        data.put("TransactionDesc", "Payment");
        data.put("AccountReference", accountNumber);
        data.put("CallBackURL", resolveCallbackUrl(callbackUrl, "callback_url"));

        return mpesaRequest(url, data, shortCodeType);
    }

    public Map<String, Object> stkquery(String checkoutRequestId, String shortCodeType) throws Exception {
        String url = baseUrl + "/mpesa/stkpushquery/v1/query";
        Map<String, Object> data = new HashMap<>();
        data.put("BusinessShortCode", getConfig("shortcode"));
        data.put("Password", lipaNaMpesaPassword());
        data.put("Timestamp", getFormattedTimestamp());
        data.put("CheckoutRequestID", checkoutRequestId);

        return mpesaRequest(url, data, shortCodeType != null ? shortCodeType : "C2B");
    }

    public Map<String, Object> b2c(Map<String, Object> params) throws Exception {
        String phonenumber = (String) params.get("phonenumber");
        String commandId = (String) params.get("commandId");
        Number amount = (Number) params.get("amount");
        String remarks = (String) params.get("remarks");
        String resultUrl = (String) params.get("resultUrl");
        String timeoutUrl = (String) params.get("timeoutUrl");
        String shortCodeType = params.containsKey("shortCodeType") ? (String) params.get("shortCodeType") : "B2C";

        String url = baseUrl + "/mpesa/b2c/v1/paymentrequest";
        Map<String, Object> body = new HashMap<>();
        body.put("InitiatorName", getConfig("initiator_name"));
        body.put("SecurityCredential", generateSecurityCredential());
        body.put("CommandID", commandId);
        body.put("Amount", amount.intValue());
        body.put("PartyA", getConfig("b2c_shortcode"));
        body.put("PartyB", phoneValidator(phonenumber));
        body.put("Remarks", remarks);
        body.put("Occassion", "");
        body.put("ResultURL", resolveCallbackUrl(resultUrl, "b2c_result_url"));
        body.put("QueueTimeOutURL", resolveCallbackUrl(timeoutUrl, "b2c_timeout_url"));

        return mpesaRequest(url, body, shortCodeType);
    }

    public Map<String, Object> validated_b2c(Map<String, Object> params) throws Exception {
        String phonenumber = (String) params.get("phonenumber");
        String commandId = (String) params.get("commandId");
        Number amount = (Number) params.get("amount");
        String remarks = (String) params.get("remarks");
        String idNumber = (String) params.get("idNumber");
        String resultUrl = (String) params.get("resultUrl");
        String timeoutUrl = (String) params.get("timeoutUrl");
        String shortCodeType = params.containsKey("shortCodeType") ? (String) params.get("shortCodeType") : "B2C";

        String url = baseUrl + "/mpesa/b2cvalidate/v2/paymentrequest";
        Map<String, Object> body = new HashMap<>();
        body.put("InitiatorName", getConfig("initiator_name"));
        body.put("SecurityCredential", generateSecurityCredential());
        body.put("CommandID", commandId);
        body.put("Amount", amount.intValue());
        body.put("PartyA", getConfig("b2c_shortcode"));
        body.put("PartyB", phoneValidator(phonenumber));
        body.put("Remarks", remarks);
        body.put("Occassion", "");
        body.put("OriginatorConversationID", getFormattedTimestamp());
        body.put("IDType", "01");
        body.put("IDNumber", idNumber);
        body.put("ResultURL", resolveCallbackUrl(resultUrl, "b2c_result_url"));
        body.put("QueueTimeOutURL", resolveCallbackUrl(timeoutUrl, "b2c_timeout_url"));

        return mpesaRequest(url, body, shortCodeType);
    }

    public Map<String, Object> b2b(Map<String, Object> params) throws Exception {
        String receiverShortcode = (String) params.get("receiverShortcode");
        String commandId = (String) params.get("commandId");
        Number amount = (Number) params.get("amount");
        String remarks = (String) params.get("remarks");
        String accountNumber = (String) params.get("accountNumber");
        String resultUrl = (String) params.get("resultUrl");
        String timeoutUrl = (String) params.get("timeoutUrl");
        String shortCodeType = params.containsKey("shortCodeType") ? (String) params.get("shortCodeType") : "B2B";

        if (commandId.equals("BusinessPayBill") && (accountNumber == null || accountNumber.isEmpty())) {
            throw new Exception("Account Number is required for BusinessPayBill CommandID");
        }

        String url = baseUrl + "/mpesa/b2b/v1/paymentrequest";
        Map<String, Object> body = new HashMap<>();
        body.put("Initiator", getConfig("initiator_name"));
        body.put("SecurityCredential", generateSecurityCredential());
        body.put("CommandID", commandId);
        body.put("SenderIdentifierType", "4");
        body.put("RecieverIdentifierType", "4");
        body.put("Amount", amount.intValue());
        body.put("PartyA", getConfig("b2c_shortcode"));
        body.put("PartyB", receiverShortcode);
        body.put("AccountReference", accountNumber);
        body.put("Remarks", remarks);
        body.put("ResultURL", resolveCallbackUrl(resultUrl, "b2b_result_url"));
        body.put("QueueTimeOutURL", resolveCallbackUrl(timeoutUrl, "b2b_timeout_url"));

        return mpesaRequest(url, body, shortCodeType);
    }

    public Map<String, Object> c2bregisterURLS(Map<String, Object> params) throws Exception {
        String shortcode = (String) params.get("shortcode");
        String confirmUrl = (String) params.get("confirmUrl");
        String validateUrl = (String) params.get("validateUrl");
        String shortCodeType = params.containsKey("shortCodeType") ? (String) params.get("shortCodeType") : "C2B";

        String url = baseUrl + "/mpesa/c2b/v2/registerurl";
        Map<String, Object> body = new HashMap<>();
        body.put("ShortCode", shortcode);
        body.put("ResponseType", "Completed");
        body.put("ConfirmationURL", resolveCallbackUrl(confirmUrl, "c2b_confirmation_url"));
        body.put("ValidationURL", resolveCallbackUrl(validateUrl, "c2b_validation_url"));

        return mpesaRequest(url, body, shortCodeType);
    }

    public Map<String, Object> c2bsimulate(Map<String, Object> params) throws Exception {
        String phonenumber = (String) params.get("phonenumber");
        Number amount = (Number) params.get("amount");
        String shortcode = (String) params.get("shortcode");
        String commandId = (String) params.get("commandId");
        String accountNumber = (String) params.get("accountNumber");
        String shortCodeType = params.containsKey("shortCodeType") ? (String) params.get("shortCodeType") : "C2B";

        String url = baseUrl + "/mpesa/c2b/v2/simulate";
        Map<String, Object> data = new HashMap<>();
        data.put("Msisdn", phoneValidator(phonenumber));
        data.put("Amount", amount.intValue());
        
        if (commandId.equals(PAYBILL)) {
            data.put("BillRefNumber", accountNumber);
        }
        data.put("CommandID", commandId);
        data.put("ShortCode", shortcode);

        return mpesaRequest(url, data, shortCodeType);
    }

    public Map<String, Object> transactionStatus(Map<String, Object> params) throws Exception {
        String shortcode = (String) params.get("shortcode");
        String transactionId = (String) params.get("transactionId");
        Number identifierType = (Number) params.get("identifierType");
        String remarks = (String) params.get("remarks");
        String resultUrl = (String) params.get("resultUrl");
        String timeoutUrl = (String) params.get("timeoutUrl");
        String shortCodeType = params.containsKey("shortCodeType") ? (String) params.get("shortCodeType") : "C2B";

        String url = baseUrl + "/mpesa/transactionstatus/v1/query";
        Map<String, Object> body = new HashMap<>();
        body.put("Initiator", getConfig("initiator_name"));
        body.put("SecurityCredential", generateSecurityCredential());
        body.put("CommandID", "TransactionStatusQuery");
        body.put("TransactionID", transactionId);
        body.put("PartyA", shortcode);
        body.put("IdentifierType", identifierType.intValue());
        body.put("Remarks", remarks);
        body.put("Occassion", "");
        body.put("ResultURL", resolveCallbackUrl(resultUrl, "status_result_url"));
        body.put("QueueTimeOutURL", resolveCallbackUrl(timeoutUrl, "status_timeout_url"));

        return mpesaRequest(url, body, shortCodeType);
    }

    public Map<String, Object> accountBalance(Map<String, Object> params) throws Exception {
        String shortcode = (String) params.get("shortcode");
        Number identifierType = (Number) params.get("identifierType");
        String remarks = (String) params.get("remarks");
        String resultUrl = (String) params.get("resultUrl");
        String timeoutUrl = (String) params.get("timeoutUrl");
        String shortCodeType = params.containsKey("shortCodeType") ? (String) params.get("shortCodeType") : "C2B";

        String url = baseUrl + "/mpesa/accountbalance/v1/query";
        Map<String, Object> body = new HashMap<>();
        body.put("Initiator", getConfig("initiator_name"));
        body.put("SecurityCredential", generateSecurityCredential());
        body.put("CommandID", "AccountBalance");
        body.put("PartyA", shortcode);
        body.put("IdentifierType", identifierType.intValue());
        body.put("Remarks", remarks);
        body.put("ResultURL", resolveCallbackUrl(resultUrl, "balance_result_url"));
        body.put("QueueTimeOutURL", resolveCallbackUrl(timeoutUrl, "balance_timeout_url"));

        return mpesaRequest(url, body, shortCodeType);
    }

    public Map<String, Object> reversal(Map<String, Object> params) throws Exception {
        String shortcode = (String) params.get("shortcode");
        String transactionId = (String) params.get("transactionId");
        Number amount = (Number) params.get("amount");
        String remarks = (String) params.get("remarks");
        String resultUrl = (String) params.get("resultUrl");
        String timeoutUrl = (String) params.get("timeoutUrl");
        String shortCodeType = params.containsKey("shortCodeType") ? (String) params.get("shortCodeType") : "C2B";

        String url = baseUrl + "/mpesa/reversal/v1/request";
        Map<String, Object> body = new HashMap<>();
        body.put("Initiator", getConfig("initiator_name"));
        body.put("SecurityCredential", generateSecurityCredential());
        body.put("CommandID", "TransactionReversal");
        body.put("TransactionID", transactionId);
        body.put("Amount", amount.intValue());
        body.put("ReceiverParty", shortcode);
        body.put("RecieverIdentifierType", "11");
        body.put("Remarks", remarks);
        body.put("Occasion", "");
        body.put("ResultURL", resolveCallbackUrl(resultUrl, "reversal_result_url"));
        body.put("QueueTimeOutURL", resolveCallbackUrl(timeoutUrl, "reversal_timeout_url"));

        return mpesaRequest(url, body, shortCodeType);
    }

    public Map<String, Object> b2pochi(Map<String, Object> params) throws Exception {
        String phonenumber = (String) params.get("phonenumber");
        Number amount = (Number) params.get("amount");
        String remarks = (String) params.get("remarks");
        String occasion = (String) params.get("occasion");
        String resultUrl = (String) params.get("resultUrl");
        String timeoutUrl = (String) params.get("timeoutUrl");
        String shortCodeType = params.containsKey("shortCodeType") ? (String) params.get("shortCodeType") : "B2C";

        String url = baseUrl + "/mpesa/b2pochi/v1/paymentrequest";
        Map<String, Object> body = new HashMap<>();
        body.put("OriginatorConversationID", getFormattedTimestamp());
        body.put("InitiatorName", getConfig("initiator_name"));
        body.put("SecurityCredential", generateSecurityCredential());
        body.put("CommandID", "BusinessPayToPochi");
        body.put("Amount", amount.intValue());
        body.put("PartyA", getConfig("b2c_shortcode"));
        body.put("PartyB", phoneValidator(phonenumber));
        body.put("Remarks", remarks);
        body.put("Occassion", occasion != null ? occasion : "");
        body.put("ResultURL", resolveCallbackUrl(resultUrl, "b2pochi_result_url"));
        body.put("QueueTimeOutURL", resolveCallbackUrl(timeoutUrl, "b2pochi_timeout_url"));

        return mpesaRequest(url, body, shortCodeType);
    }
}