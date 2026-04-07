namespace Yourdudeken.Mpesa;

public class MpesaConfig
{
    public string Environment { get; set; } = "sandbox";
    public string MpesaConsumerKey { get; set; } = string.Empty;
    public string MpesaConsumerSecret { get; set; } = string.Empty;
    public string? B2cConsumerKey { get; set; }
    public string? B2cConsumerSecret { get; set; }
    public string Passkey { get; set; } = string.Empty;
    public string Shortcode { get; set; } = string.Empty;
    public string? TillNumber { get; set; }
    public string InitiatorName { get; set; } = string.Empty;
    public string InitiatorPassword { get; set; } = string.Empty;
    public string? B2cShortcode { get; set; }
    public Dictionary<string, string>? Callbacks { get; set; }
}