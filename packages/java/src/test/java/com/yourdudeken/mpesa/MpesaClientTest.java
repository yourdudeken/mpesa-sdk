package com.yourdudeken.mpesa;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

public class MpesaClientTest {

    @Test
    public void testConfigCreation() {
        MpesaConfig config = new MpesaConfig();
        config.setEnvironment("sandbox");
        config.setMpesaConsumerKey("test_key");
        config.setMpesaConsumerSecret("test_secret");
        
        assertEquals("sandbox", config.getEnvironment());
        assertEquals("test_key", config.getMpesaConsumerKey());
    }
}
