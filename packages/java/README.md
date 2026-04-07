# Mpesa Java SDK

A Java SDK for Mpesa Daraja API.

## Installation

```xml
<dependency>
    <groupId>com.yourdudeken</groupId>
    <artifactId>mpesa</artifactId>
    <version>1.0.0</version>
</dependency>
```

## Usage

```java
import com.yourdudeken.mpesa.MpesaClient;
import com.yourdudeken.mpesa.config.MpesaConfig;

MpesaConfig config = new MpesaConfig();
config.setEnvironment("sandbox");
config.setMpesaConsumerKey("your_key");
config.setMpesaConsumerSecret("your_secret");
config.setShortcode("174379");
config.setPasskey("your_passkey");

MpesaClient mpesa = new MpesaClient(config);

// STK Push
String response = mpesa.stkPush("254712345678", 100, "ORDER123");
```

## License

MIT License
