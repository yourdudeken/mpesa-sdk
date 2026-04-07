package com.yourdudeken.mpesa.models;

public class STKPushRequest {
    private String businessShortCode;
    private String password;
    private String timestamp;
    private int amount;
    private String partyA;
    private String partyB;
    private String transactionType;
    private String phoneNumber;
    private String accountReference;
    private String callBackURL;

    public STKPushRequest() {}

    public String getBusinessShortCode() { return businessShortCode; }
    public void setBusinessShortCode(String businessShortCode) { this.businessShortCode = businessShortCode; }

    public String getPassword() { return password; }
    public void setPassword(String password) { this.password = password; }

    public String getTimestamp() { return timestamp; }
    public void setTimestamp(String timestamp) { this.timestamp = timestamp; }

    public int getAmount() { return amount; }
    public void setAmount(int amount) { this.amount = amount; }

    public String getPartyA() { return partyA; }
    public void setPartyA(String partyA) { this.partyA = partyA; }

    public String getPartyB() { return partyB; }
    public void setPartyB(String partyB) { this.partyB = partyB; }

    public String getTransactionType() { return transactionType; }
    public void setTransactionType(String transactionType) { this.transactionType = transactionType; }

    public String getPhoneNumber() { return phoneNumber; }
    public void setPhoneNumber(String phoneNumber) { this.phoneNumber = phoneNumber; }

    public String getAccountReference() { return accountReference; }
    public void setAccountReference(String accountReference) { this.accountReference = accountReference; }

    public String getCallBackURL() { return callBackURL; }
    public void setCallBackURL(String callBackURL) { this.callBackURL = callBackURL; }
}
