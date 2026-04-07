package com.yourdudeken.mpesa;

import java.util.Map;

public class MpesaConfig {
    private String environment;
    private String mpesaConsumerKey;
    private String mpesaConsumerSecret;
    private String b2cConsumerKey;
    private String b2cConsumerSecret;
    private String passkey;
    private String shortcode;
    private String tillNumber;
    private String initiatorName;
    private String initiatorPassword;
    private String b2cShortcode;
    private Map<String, String> callbacks;

    public MpesaConfig() {}

    public String getEnvironment() { return environment; }
    public void setEnvironment(String environment) { this.environment = environment; }

    public String getMpesaConsumerKey() { return mpesaConsumerKey; }
    public void setMpesaConsumerKey(String mpesaConsumerKey) { this.mpesaConsumerKey = mpesaConsumerKey; }

    public String getMpesaConsumerSecret() { return mpesaConsumerSecret; }
    public void setMpesaConsumerSecret(String mpesaConsumerSecret) { this.mpesaConsumerSecret = mpesaConsumerSecret; }

    public String getB2cConsumerKey() { return b2cConsumerKey; }
    public void setB2cConsumerKey(String b2cConsumerKey) { this.b2cConsumerKey = b2cConsumerKey; }

    public String getB2cConsumerSecret() { return b2cConsumerSecret; }
    public void setB2cConsumerSecret(String b2cConsumerSecret) { this.b2cConsumerSecret = b2cConsumerSecret; }

    public String getPasskey() { return passkey; }
    public void setPasskey(String passkey) { this.passkey = passkey; }

    public String getShortcode() { return shortcode; }
    public void setShortcode(String shortcode) { this.shortcode = shortcode; }

    public String getTillNumber() { return tillNumber; }
    public void setTillNumber(String tillNumber) { this.tillNumber = tillNumber; }

    public String getInitiatorName() { return initiatorName; }
    public void setInitiatorName(String initiatorName) { this.initiatorName = initiatorName; }

    public String getInitiatorPassword() { return initiatorPassword; }
    public void setInitiatorPassword(String initiatorPassword) { this.initiatorPassword = initiatorPassword; }

    public String getB2cShortcode() { return b2cShortcode; }
    public void setB2cShortcode(String b2cShortcode) { this.b2cShortcode = b2cShortcode; }

    public Map<String, String> getCallbacks() { return callbacks; }
    public void setCallbacks(Map<String, String> callbacks) { this.callbacks = callbacks; }
}