using Xunit;
using Yourdudeken.Mpesa;
using Yourdudeken.Mpesa.Config;
using Moq;
using System.Net.Http;
using System.Collections.Generic;
using System.Threading.Tasks;
using Newtonsoft.Json;

namespace Yourdudeken.Mpesa.Tests;

public class MpesaClientTests
{
    private MpesaConfig GetTestConfig()
    {
        return new MpesaConfig
        {
            Environment = "sandbox",
            MpesaConsumerKey = "test_key",
            MpesaConsumerSecret = "test_secret",
            Passkey = "test_passkey",
            Shortcode = "174379",
            InitiatorName = "testapi",
            InitiatorPassword = "test_password",
            B2cShortcode = "600000",
            Callbacks = new Dictionary<string, string>
            {
                { "callback_url", "https://test.com/callback" },
                { "b2c_result_url", "https://test.com/b2c_result" },
                { "b2c_timeout_url", "https://test.com/b2c_timeout" },
                { "b2b_result_url", "https://test.com/b2b_result" },
                { "b2b_timeout_url", "https://test.com/b2b_timeout" },
                { "b2pochi_result_url", "https://test.com/b2pochi_result" },
                { "b2pochi_timeout_url", "https://test.com/b2pochi_timeout" },
                { "c2b_validation_url", "https://test.com/c2b_validate" },
                { "c2b_confirmation_url", "https://test.com/c2b_confirm" },
                { "balance_result_url", "https://test.com/balance_result" },
                { "balance_timeout_url", "https://test.com/balance_timeout" },
                { "status_result_url", "https://test.com/status_result" },
                { "status_timeout_url", "https://test.com/status_timeout" },
                { "reversal_result_url", "https://test.com/reversal_result" },
                { "reversal_timeout_url", "https://test.com/reversal_timeout" },
            }
        };
    }

    [Fact]
    public void TestConfigCreation()
    {
        var config = GetTestConfig();

        Assert.Equal("sandbox", config.Environment);
        Assert.Equal("test_key", config.MpesaConsumerKey);
        Assert.Equal("174379", config.Shortcode);
        Assert.Equal("600000", config.B2cShortcode);
    }

    [Fact]
    public void TestMpesaInitialization()
    {
        var config = GetTestConfig();
        var mpesa = new Mpesa(config);
        Assert.NotNull(mpesa);
    }

    [Fact]
    public void TestStaticConstants()
    {
        Assert.Equal("CustomerPayBillOnline", Mpesa.PAYBILL);
        Assert.Equal("CustomerBuyGoodsOnline", Mpesa.TILL);
    }

    [Fact]
    public void TestPhoneValidator()
    {
        var config = GetTestConfig();
        var mpesa = new Mpesa(config);

        var phoneValidatorMethod = typeof(Mpesa).GetPrivateMethod("PhoneValidator");
        
        Assert.Equal("254712345678", phoneValidatorMethod.Invoke(mpesa, new object[] { "+254712345678" }));
        Assert.Equal("254712345678", phoneValidatorMethod.Invoke(mpesa, new object[] { "0712345678" }));
        Assert.Equal("254712345678", phoneValidatorMethod.Invoke(mpesa, new object[] { "712345678" }));
    }

    [Fact]
    public void TestStkpush_ThrowsException_WhenAccountReferenceMissing()
    {
        var config = GetTestConfig();
        var mpesa = new Mpesa(config);

        var params = new Dictionary<string, object>
        {
            { "phonenumber", "254712345678" },
            { "amount", 100 },
            { "accountNumber", "" }
        };

        Assert.ThrowsAsync<Exception>(async () => await mpesa.Stkpush(params));
    }

    [Fact]
    public void TestStkpush_ThrowsException_WhenTillNumberRequired()
    {
        var config = GetTestConfig();
        var mpesa = new Mpesa(config);

        var params = new Dictionary<string, object>
        {
            { "phonenumber", "254712345678" },
            { "amount", 100 },
            { "accountNumber", "12345" },
            { "transactionType", "CustomerBuyGoodsOnline" }
        };

        Assert.ThrowsAsync<Exception>(async () => await mpesa.Stkpush(params));
    }

    [Fact]
    public void TestB2b_ThrowsException_WhenAccountNumberMissing()
    {
        var config = GetTestConfig();
        var mpesa = new Mpesa(config);

        var params = new Dictionary<string, object>
        {
            { "receiverShortcode", "600000" },
            { "commandId", "BusinessPayBill" },
            { "amount", 100 },
            { "remarks", "Test payment" }
        };

        Assert.ThrowsAsync<Exception>(async () => await mpesa.B2b(params));
    }

    [Fact]
    public void TestGetFormattedTimestamp()
    {
        var config = GetTestConfig();
        var mpesa = new Mpesa(config);

        var timestampMethod = typeof(Mpesa).GetPrivateMethod("GetFormattedTimestamp");
        var timestamp = timestampMethod.Invoke(mpesa, null) as string;

        Assert.NotNull(timestamp);
        Assert.Equal(14, timestamp?.Length);
    }

    [Fact]
    public void TestResolveCallbackUrl_ThrowsException_WhenUrlNotProvided()
    {
        var config = new MpesaConfig
        {
            Environment = "sandbox",
            MpesaConsumerKey = "test_key",
            MpesaConsumerSecret = "test_secret",
            Passkey = "test_passkey",
            Shortcode = "174379",
            InitiatorName = "testapi",
            InitiatorPassword = "test_password"
        };

        var mpesa = new Mpesa(config);

        var resolveMethod = typeof(Mpesa).GetPrivateMethod("ResolveCallbackUrl");
        
        Assert.Throws<Exception>(() => resolveMethod.Invoke(mpesa, new object[] { null, "callback_url" }));
    }

    [Fact]
    public void TestResolveCallbackUrl_UsesParamUrl_WhenProvided()
    {
        var config = GetTestConfig();
        var mpesa = new Mpesa(config);

        var resolveMethod = typeof(Mpesa).GetPrivateMethod("ResolveCallbackUrl");
        var result = resolveMethod.Invoke(mpesa, new object[] { "https://custom.url/callback", "callback_url" }) as string;

        Assert.Equal("https://custom.url/callback", result);
    }

    [Fact]
    public void TestGetConfig_ReturnsCorrectValues()
    {
        var config = GetTestConfig();
        var mpesa = new Mpesa(config);

        var getConfigMethod = typeof(Mpesa).GetPrivateMethod("GetConfig");

        Assert.Equal("sandbox", getConfigMethod.Invoke(mpesa, new object[] { "environment" }));
        Assert.Equal("test_key", getConfigMethod.Invoke(mpesa, new object[] { "mpesa_consumer_key" }));
        Assert.Equal("test_passkey", getConfigMethod.Invoke(mpesa, new object[] { "passkey" }));
        Assert.Equal("174379", getConfigMethod.Invoke(mpesa, new object[] { "shortcode" }));
    }
}

public static class PrivateMethodExtensions
{
    public static System.Reflection.MethodInfo GetPrivateMethod(this Type type, string methodName)
    {
        return type.GetMethod(methodName, System.Reflection.BindingFlags.NonPublic | System.Reflection.BindingFlags.Instance);
    }
}