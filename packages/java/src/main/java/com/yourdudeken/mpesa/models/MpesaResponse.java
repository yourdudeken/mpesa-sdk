package com.yourdudeken.mpesa.models;

public class MpesaResponse {
    private String responseCode;
    private String responseDescription;
    private String merchantRequestID;
    private String checkoutRequestID;
    private String customerMessage;

    public MpesaResponse() {}

    public String getResponseCode() { return responseCode; }
    public void setResponseCode(String responseCode) { this.responseCode = responseCode; }

    public String getResponseDescription() { return responseDescription; }
    public void setResponseDescription(String responseDescription) { this.responseDescription = responseDescription; }

    public String getMerchantRequestID() { return merchantRequestID; }
    public void setMerchantRequestID(String merchantRequestID) { this.merchantRequestID = merchantRequestID; }

    public String getCheckoutRequestID() { return checkoutRequestID; }
    public void setCheckoutRequestID(String checkoutRequestID) { this.checkoutRequestID = checkoutRequestID; }

    public String getCustomerMessage() { return customerMessage; }
    public void setCustomerMessage(String customerMessage) { this.customerMessage = customerMessage; }

    public boolean isSuccess() {
        return "0".equals(responseCode);
    }
}
