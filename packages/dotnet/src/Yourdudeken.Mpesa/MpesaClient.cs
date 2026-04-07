using System;
using System.Collections.Generic;
using System.IO;
using System.Net.Http;
using System.Security.Cryptography;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;

namespace Yourdudeken.Mpesa;

public class Mpesa
{
    public const string PAYBILL = "CustomerPayBillOnline";
    public const string TILL = "CustomerBuyGoodsOnline";

    private readonly MpesaConfig _config;
    private readonly string _baseUrl;
    private string? _accessToken;
    private long _tokenExpiry;

    public Mpesa(MpesaConfig config)
    {
        _config = config;
        _baseUrl = config.Environment == "sandbox"
            ? "https://sandbox.safaricom.co.ke"
            : "https://api.safaricom.co.ke";
    }

    private string? GetConfig(string key)
    {
        return key switch
        {
            "environment" => _config.Environment,
            "mpesa_consumer_key" => _config.MpesaConsumerKey,
            "mpesa_consumer_secret" => _config.MpesaConsumerSecret,
            "b2c_consumer_key" => _config.B2cConsumerKey,
            "b2c_consumer_secret" => _config.B2cConsumerSecret,
            "passkey" => _config.Passkey,
            "shortcode" => _config.Shortcode,
            "till_number" => _config.TillNumber,
            "initiator_name" => _config.InitiatorName,
            "initiator_password" => _config.InitiatorPassword,
            "b2c_shortcode" => _config.B2cShortcode,
            _ => null
        };
    }

    private string ResolveCallbackUrl(string? paramUrl, string configKey)
    {
        string? configUrl = _config.Callbacks?.GetValueOrDefault(configKey);
        if (!string.IsNullOrEmpty(paramUrl)) return paramUrl;
        if (!string.IsNullOrEmpty(configUrl)) return configUrl;
        throw new Exception($"Ensure you have set the {configKey} in the config or passed as a parameter");
    }

    private string PhoneValidator(string phoneNo)
    {
        if (phoneNo.StartsWith("+")) phoneNo = phoneNo.Substring(1);
        if (phoneNo.StartsWith("0")) phoneNo = "254" + phoneNo.Substring(1);
        else if (phoneNo.StartsWith("7")) phoneNo = "254" + phoneNo;
        return phoneNo;
    }

    private string GetFormattedTimestamp()
    {
        return DateTime.Now.ToString("yyyyMMddHHmmss");
    }

    private string LipaNaMpesaPassword()
    {
        string timestamp = GetFormattedTimestamp();
        string password = GetConfig("shortcode") + GetConfig("passkey") + timestamp;
        return Convert.ToBase64String(Encoding.UTF8.GetBytes(password));
    }

    private async Task<string> GenerateAccessToken(string shortCodeType)
    {
        long currentTime = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds();
        if (!string.IsNullOrEmpty(_accessToken) && currentTime < _tokenExpiry)
        {
            return _accessToken;
        }

        string consumerKey = (shortCodeType == "B2C" || shortCodeType == "B2B")
            ? GetConfig("b2c_consumer_key") ?? ""
            : GetConfig("mpesa_consumer_key") ?? "";
        string consumerSecret = (shortCodeType == "B2C" || shortCodeType == "B2B")
            ? GetConfig("b2c_consumer_secret") ?? ""
            : GetConfig("mpesa_consumer_secret") ?? "";

        string auth = Convert.ToBase64String(Encoding.UTF8.GetBytes(consumerKey + ":" + consumerSecret));

        using var client = new HttpClient();
        var request = new HttpRequestMessage(HttpMethod.Get, _baseUrl + "/oauth/v1/generate?grant_type=client_credentials");
        request.Headers.Add("Authorization", "Basic " + auth);

        var response = await client.SendAsync(request);
        var jsonResponse = await response.Content.ReadAsStringAsync();
        var data = JsonConvert.DeserializeObject<Dictionary<string, object>>(jsonResponse);

        _accessToken = data?["access_token"]?.ToString();
        _tokenExpiry = currentTime + long.Parse(data?["expires_in"]?.ToString() ?? "3600") * 1000 - 60000;

        return _accessToken ?? "";
    }

    private string GenerateSecurityCredential()
    {
        string certPath = GetConfig("environment") == "sandbox"
            ? Path.Combine(AppContext.BaseDirectory, "certificates", "SandboxCertificate.cer")
            : Path.Combine(AppContext.BaseDirectory, "certificates", "ProductionCertificate.cer");

        if (!File.Exists(certPath))
        {
            certPath = GetConfig("environment") == "sandbox"
                ? "certificates/SandboxCertificate.cer"
                : "certificates/ProductionCertificate.cer";
        }

        string pubkey = File.ReadAllText(certPath);

        var rsa = RSA.Create();
        rsa.ImportFromPem(pubkey.ToCharArray());

        byte[] encrypted = rsa.Encrypt(Encoding.UTF8.GetBytes(GetConfig("initiator_password") ?? ""), RSAEncryptionPadding.Pkcs1);

        return Convert.ToBase64String(encrypted);
    }

    private async Task<Dictionary<string, object>> MpesaRequest(string url, Dictionary<string, object> body, string shortCodeType)
    {
        string token = await GenerateAccessToken(shortCodeType);

        using var client = new HttpClient();
        var request = new HttpRequestMessage(HttpMethod.Post, url);
        request.Headers.Add("Authorization", "Bearer " + token);
        request.Content = new StringContent(JsonConvert.SerializeObject(body), Encoding.UTF8, "application/json");

        var response = await client.SendAsync(request);
        var jsonResponse = await response.Content.ReadAsStringAsync();
        return JsonConvert.DeserializeObject<Dictionary<string, object>>(jsonResponse) ?? new Dictionary<string, object>();
    }

    public async Task<Dictionary<string, object>> Stkpush(Dictionary<string, object> @params)
    {
        string phonenumber = @params["phonenumber"].ToString() ?? "";
        int amount = Convert.ToInt32(@params["amount"]);
        string accountNumber = @params["accountNumber"].ToString() ?? "";
        string? callbackUrl = @params.GetValueOrDefault("callbackUrl")?.ToString();
        string transactionType = @params.GetValueOrDefault("transactionType")?.ToString() ?? PAYBILL;
        string shortCodeType = @params.GetValueOrDefault("shortCodeType")?.ToString() ?? "C2B";

        if (string.IsNullOrEmpty(accountNumber))
            throw new Exception("An Account Reference is required for All transactions.");

        if (transactionType == TILL && string.IsNullOrEmpty(GetConfig("till_number")))
            throw new Exception("Till number is required for Buy Goods transactions.");

        string url = _baseUrl + "/mpesa/stkpush/v1/processrequest";
        var data = new Dictionary<string, object>
        {
            ["BusinessShortCode"] = GetConfig("shortcode") ?? "",
            ["Password"] = LipaNaMpesaPassword(),
            ["Timestamp"] = GetFormattedTimestamp(),
            ["Amount"] = amount,
            ["PartyA"] = PhoneValidator(phonenumber),
            ["PartyB"] = transactionType == PAYBILL ? (GetConfig("shortcode") ?? "") : (GetConfig("till_number") ?? ""),
            ["TransactionType"] = transactionType,
            ["PhoneNumber"] = PhoneValidator(phonenumber),
            ["TransactionDesc"] = "Payment",
            ["AccountReference"] = accountNumber,
            ["CallBackURL"] = ResolveCallbackUrl(callbackUrl, "callback_url")
        };

        return await MpesaRequest(url, data, shortCodeType);
    }

    public async Task<Dictionary<string, object>> Stkquery(string checkoutRequestId, string shortCodeType = "C2B")
    {
        string url = _baseUrl + "/mpesa/stkpushquery/v1/query";
        var data = new Dictionary<string, object>
        {
            ["BusinessShortCode"] = GetConfig("shortcode") ?? "",
            ["Password"] = LipaNaMpesaPassword(),
            ["Timestamp"] = GetFormattedTimestamp(),
            ["CheckoutRequestID"] = checkoutRequestId
        };

        return await MpesaRequest(url, data, shortCodeType);
    }

    public async Task<Dictionary<string, object>> B2c(Dictionary<string, object> @params)
    {
        string phonenumber = @params["phonenumber"].ToString() ?? "";
        string commandId = @params["commandId"].ToString() ?? "";
        int amount = Convert.ToInt32(@params["amount"]);
        string remarks = @params["remarks"].ToString() ?? "";
        string? resultUrl = @params.GetValueOrDefault("resultUrl")?.ToString();
        string? timeoutUrl = @params.GetValueOrDefault("timeoutUrl")?.ToString();
        string shortCodeType = @params.GetValueOrDefault("shortCodeType")?.ToString() ?? "B2C";

        string url = _baseUrl + "/mpesa/b2c/v1/paymentrequest";
        var body = new Dictionary<string, object>
        {
            ["InitiatorName"] = GetConfig("initiator_name") ?? "",
            ["SecurityCredential"] = GenerateSecurityCredential(),
            ["CommandID"] = commandId,
            ["Amount"] = amount,
            ["PartyA"] = GetConfig("b2c_shortcode") ?? "",
            ["PartyB"] = PhoneValidator(phonenumber),
            ["Remarks"] = remarks,
            ["Occassion"] = "",
            ["ResultURL"] = ResolveCallbackUrl(resultUrl, "b2c_result_url"),
            ["QueueTimeOutURL"] = ResolveCallbackUrl(timeoutUrl, "b2c_timeout_url")
        };

        return await MpesaRequest(url, body, shortCodeType);
    }

    public async Task<Dictionary<string, object>> Validated_b2c(Dictionary<string, object> @params)
    {
        string phonenumber = @params["phonenumber"].ToString() ?? "";
        string commandId = @params["commandId"].ToString() ?? "";
        int amount = Convert.ToInt32(@params["amount"]);
        string remarks = @params["remarks"].ToString() ?? "";
        string idNumber = @params["idNumber"].ToString() ?? "";
        string? resultUrl = @params.GetValueOrDefault("resultUrl")?.ToString();
        string? timeoutUrl = @params.GetValueOrDefault("timeoutUrl")?.ToString();
        string shortCodeType = @params.GetValueOrDefault("shortCodeType")?.ToString() ?? "B2C";

        string url = _baseUrl + "/mpesa/b2cvalidate/v2/paymentrequest";
        var body = new Dictionary<string, object>
        {
            ["InitiatorName"] = GetConfig("initiator_name") ?? "",
            ["SecurityCredential"] = GenerateSecurityCredential(),
            ["CommandID"] = commandId,
            ["Amount"] = amount,
            ["PartyA"] = GetConfig("b2c_shortcode") ?? "",
            ["PartyB"] = PhoneValidator(phonenumber),
            ["Remarks"] = remarks,
            ["Occassion"] = "",
            ["OriginatorConversationID"] = GetFormattedTimestamp(),
            ["IDType"] = "01",
            ["IDNumber"] = idNumber,
            ["ResultURL"] = ResolveCallbackUrl(resultUrl, "b2c_result_url"),
            ["QueueTimeOutURL"] = ResolveCallbackUrl(timeoutUrl, "b2c_timeout_url")
        };

        return await MpesaRequest(url, body, shortCodeType);
    }

    public async Task<Dictionary<string, object>> B2b(Dictionary<string, object> @params)
    {
        string receiverShortcode = @params["receiverShortcode"].ToString() ?? "";
        string commandId = @params["commandId"].ToString() ?? "";
        int amount = Convert.ToInt32(@params["amount"]);
        string remarks = @params["remarks"].ToString() ?? "";
        string? accountNumber = @params.GetValueOrDefault("accountNumber")?.ToString();
        string? resultUrl = @params.GetValueOrDefault("resultUrl")?.ToString();
        string? timeoutUrl = @params.GetValueOrDefault("timeoutUrl")?.ToString();
        string shortCodeType = @params.GetValueOrDefault("shortCodeType")?.ToString() ?? "B2B";

        if (commandId == "BusinessPayBill" && string.IsNullOrEmpty(accountNumber))
            throw new Exception("Account Number is required for BusinessPayBill CommandID");

        string url = _baseUrl + "/mpesa/b2b/v1/paymentrequest";
        var body = new Dictionary<string, object>
        {
            ["Initiator"] = GetConfig("initiator_name") ?? "",
            ["SecurityCredential"] = GenerateSecurityCredential(),
            ["CommandID"] = commandId,
            ["SenderIdentifierType"] = "4",
            ["RecieverIdentifierType"] = "4",
            ["Amount"] = amount,
            ["PartyA"] = GetConfig("b2c_shortcode") ?? "",
            ["PartyB"] = receiverShortcode,
            ["AccountReference"] = accountNumber ?? "",
            ["Remarks"] = remarks,
            ["ResultURL"] = ResolveCallbackUrl(resultUrl, "b2b_result_url"),
            ["QueueTimeOutURL"] = ResolveCallbackUrl(timeoutUrl, "b2b_timeout_url")
        };

        return await MpesaRequest(url, body, shortCodeType);
    }

    public async Task<Dictionary<string, object>> C2bregisterURLS(Dictionary<string, object> @params)
    {
        string shortcode = @params["shortcode"].ToString() ?? "";
        string? confirmUrl = @params.GetValueOrDefault("confirmUrl")?.ToString();
        string? validateUrl = @params.GetValueOrDefault("validateUrl")?.ToString();
        string shortCodeType = @params.GetValueOrDefault("shortCodeType")?.ToString() ?? "C2B";

        string url = _baseUrl + "/mpesa/c2b/v2/registerurl";
        var body = new Dictionary<string, object>
        {
            ["ShortCode"] = shortcode,
            ["ResponseType"] = "Completed",
            ["ConfirmationURL"] = ResolveCallbackUrl(confirmUrl, "c2b_confirmation_url"),
            ["ValidationURL"] = ResolveCallbackUrl(validateUrl, "c2b_validation_url")
        };

        return await MpesaRequest(url, body, shortCodeType);
    }

    public async Task<Dictionary<string, object>> C2bsimulate(Dictionary<string, object> @params)
    {
        string phonenumber = @params["phonenumber"].ToString() ?? "";
        int amount = Convert.ToInt32(@params["amount"]);
        string shortcode = @params["shortcode"].ToString() ?? "";
        string commandId = @params["commandId"].ToString() ?? "";
        string? accountNumber = @params.GetValueOrDefault("accountNumber")?.ToString();
        string shortCodeType = @params.GetValueOrDefault("shortCodeType")?.ToString() ?? "C2B";

        string url = _baseUrl + "/mpesa/c2b/v2/simulate";
        var data = new Dictionary<string, object>
        {
            ["Msisdn"] = PhoneValidator(phonenumber),
            ["Amount"] = amount,
            ["CommandID"] = commandId,
            ["ShortCode"] = shortcode
        };

        if (commandId == PAYBILL)
        {
            data["BillRefNumber"] = accountNumber ?? "";
        }

        return await MpesaRequest(url, data, shortCodeType);
    }

    public async Task<Dictionary<string, object>> TransactionStatus(Dictionary<string, object> @params)
    {
        string shortcode = @params["shortcode"].ToString() ?? "";
        string transactionId = @params["transactionId"].ToString() ?? "";
        int identifierType = Convert.ToInt32(@params["identifierType"]);
        string remarks = @params["remarks"].ToString() ?? "";
        string? resultUrl = @params.GetValueOrDefault("resultUrl")?.ToString();
        string? timeoutUrl = @params.GetValueOrDefault("timeoutUrl")?.ToString();
        string shortCodeType = @params.GetValueOrDefault("shortCodeType")?.ToString() ?? "C2B";

        string url = _baseUrl + "/mpesa/transactionstatus/v1/query";
        var body = new Dictionary<string, object>
        {
            ["Initiator"] = GetConfig("initiator_name") ?? "",
            ["SecurityCredential"] = GenerateSecurityCredential(),
            ["CommandID"] = "TransactionStatusQuery",
            ["TransactionID"] = transactionId,
            ["PartyA"] = shortcode,
            ["IdentifierType"] = identifierType,
            ["Remarks"] = remarks,
            ["Occassion"] = "",
            ["ResultURL"] = ResolveCallbackUrl(resultUrl, "status_result_url"),
            ["QueueTimeOutURL"] = ResolveCallbackUrl(timeoutUrl, "status_timeout_url")
        };

        return await MpesaRequest(url, body, shortCodeType);
    }

    public async Task<Dictionary<string, object>> AccountBalance(Dictionary<string, object> @params)
    {
        string shortcode = @params["shortcode"].ToString() ?? "";
        int identifierType = Convert.ToInt32(@params["identifierType"]);
        string remarks = @params["remarks"].ToString() ?? "";
        string? resultUrl = @params.GetValueOrDefault("resultUrl")?.ToString();
        string? timeoutUrl = @params.GetValueOrDefault("timeoutUrl")?.ToString();
        string shortCodeType = @params.GetValueOrDefault("shortCodeType")?.ToString() ?? "C2B";

        string url = _baseUrl + "/mpesa/accountbalance/v1/query";
        var body = new Dictionary<string, object>
        {
            ["Initiator"] = GetConfig("initiator_name") ?? "",
            ["SecurityCredential"] = GenerateSecurityCredential(),
            ["CommandID"] = "AccountBalance",
            ["PartyA"] = shortcode,
            ["IdentifierType"] = identifierType,
            ["Remarks"] = remarks,
            ["ResultURL"] = ResolveCallbackUrl(resultUrl, "balance_result_url"),
            ["QueueTimeOutURL"] = ResolveCallbackUrl(timeoutUrl, "balance_timeout_url")
        };

        return await MpesaRequest(url, body, shortCodeType);
    }

    public async Task<Dictionary<string, object>> Reversal(Dictionary<string, object> @params)
    {
        string shortcode = @params["shortcode"].ToString() ?? "";
        string transactionId = @params["transactionId"].ToString() ?? "";
        int amount = Convert.ToInt32(@params["amount"]);
        string remarks = @params["remarks"].ToString() ?? "";
        string? resultUrl = @params.GetValueOrDefault("resultUrl")?.ToString();
        string? timeoutUrl = @params.GetValueOrDefault("timeoutUrl")?.ToString();
        string shortCodeType = @params.GetValueOrDefault("shortCodeType")?.ToString() ?? "C2B";

        string url = _baseUrl + "/mpesa/reversal/v1/request";
        var body = new Dictionary<string, object>
        {
            ["Initiator"] = GetConfig("initiator_name") ?? "",
            ["SecurityCredential"] = GenerateSecurityCredential(),
            ["CommandID"] = "TransactionReversal",
            ["TransactionID"] = transactionId,
            ["Amount"] = amount,
            ["ReceiverParty"] = shortcode,
            ["RecieverIdentifierType"] = "11",
            ["Remarks"] = remarks,
            ["Occasion"] = "",
            ["ResultURL"] = ResolveCallbackUrl(resultUrl, "reversal_result_url"),
            ["QueueTimeOutURL"] = ResolveCallbackUrl(timeoutUrl, "reversal_timeout_url")
        };

        return await MpesaRequest(url, body, shortCodeType);
    }

    public async Task<Dictionary<string, object>> B2pochi(Dictionary<string, object> @params)
    {
        string phonenumber = @params["phonenumber"].ToString() ?? "";
        int amount = Convert.ToInt32(@params["amount"]);
        string remarks = @params["remarks"].ToString() ?? "";
        string? occasion = @params.GetValueOrDefault("occasion")?.ToString();
        string? resultUrl = @params.GetValueOrDefault("resultUrl")?.ToString();
        string? timeoutUrl = @params.GetValueOrDefault("timeoutUrl")?.ToString();
        string shortCodeType = @params.GetValueOrDefault("shortCodeType")?.ToString() ?? "B2C";

        string url = _baseUrl + "/mpesa/b2pochi/v1/paymentrequest";
        var body = new Dictionary<string, object>
        {
            ["OriginatorConversationID"] = GetFormattedTimestamp(),
            ["InitiatorName"] = GetConfig("initiator_name") ?? "",
            ["SecurityCredential"] = GenerateSecurityCredential(),
            ["CommandID"] = "BusinessPayToPochi",
            ["Amount"] = amount,
            ["PartyA"] = GetConfig("b2c_shortcode") ?? "",
            ["PartyB"] = PhoneValidator(phonenumber),
            ["Remarks"] = remarks,
            ["Occassion"] = occasion ?? "",
            ["ResultURL"] = ResolveCallbackUrl(resultUrl, "b2pochi_result_url"),
            ["QueueTimeOutURL"] = ResolveCallbackUrl(timeoutUrl, "b2pochi_timeout_url")
        };

        return await MpesaRequest(url, body, shortCodeType);
    }
}