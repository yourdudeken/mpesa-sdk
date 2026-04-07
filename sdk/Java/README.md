# Mpesa SDK for Java

[![Maven Central](https://img.shields.io/maven-central/v/com.yourdudeken/mpesa.svg)](https://search.maven.org/artifact/com.yourdudeken/mpesa)
[![License](https://img.shields.io/github/license/yourdudeken/mpesa.svg)](LICENSE.md)

A Java SDK for the Mpesa Daraja APIs. This SDK allows you to integrate Mpesa Daraja APIs into your Java applications with ease.

## Installation

### Maven

Add the following dependency to your `pom.xml`:

```xml
<dependency>
    <groupId>com.yourdudeken</groupId>
    <artifactId>mpesa</artifactId>
    <version>1.0.0</version>
</dependency>
```

### Gradle

Add the following dependency to your `build.gradle`:

```gradle
implementation 'com.yourdudeken:mpesa:1.0.0'
```

## Usage

```java
import com.yourdudeken.mpesa.Mpesa;
import com.yourdudeken.mpesa.MpesaConfig;
import java.util.HashMap;
import java.util.Map;

MpesaConfig config = new MpesaConfig();
config.setEnvironment("sandbox");
config.setMpesaConsumerKey("your_consumer_key");
config.setMpesaConsumerSecret("your_consumer_secret");
config.setPasskey("your_passkey");
config.setShortcode("174379");
config.setInitiatorName("testapi");
config.setInitiatorPassword("your_password");

Map<String, String> callbacks = new HashMap<>();
callbacks.put("callback_url", "https://your-callback-url.com/callback");
config.setCallbacks(callbacks);

Mpesa mpesa = new Mpesa(config);

try {
    Map<String, Object> params = new HashMap<>();
    params.put("phonenumber", "254712345678");
    params.put("amount", 10);
    params.put("accountNumber", "TEST001");
    
    Map<String, Object> response = mpesa.stkpush(params);
    System.out.println(response);
} catch (Exception e) {
    e.printStackTrace();
}
```

## Supported APIs

- **STK Push** - Lipa na Mpesa Express Online
- **STK Query** - Check transaction status
- **B2C** - Business to Customer
- **B2B** - Business to Business
- **B2Pochi** - Business to Pochi La Biashara
- **C2B** - Customer to Business (Register URL & Simulate)
- **Transaction Status** - Check transaction status
- **Account Balance** - Query account balance
- **Reversal** - Reverse a transaction

## Configuration

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| environment | String | Yes | "sandbox" or "production" |
| mpesaConsumerKey | String | Yes | C2B Consumer Key |
| mpesaConsumerSecret | String | Yes | C2B Consumer Secret |
| b2cConsumerKey | String | No | B2C Consumer Key |
| b2cConsumerSecret | String | No | B2C Consumer Secret |
| passkey | String | Yes | Lipa na Mpesa Online Passkey |
| shortcode | String | Yes | Business Shortcode |
| tillNumber | String | No | Till Number |
| initiatorName | String | Yes | Mpesa Initiator Name |
| initiatorPassword | String | Yes | Mpesa Initiator Password |
| b2cShortcode | String | No | B2C Shortcode |
| callbacks | Map | No | Callback URLs |

## License

MIT License - see [LICENSE.md](LICENSE.md) for details.