<?php

namespace Yourdudeken\Mpesa\Config;

class MpesaConfig
{
    public $environment;
    public $mpesaConsumerKey;
    public $mpesaConsumerSecret;
    public $b2cConsumerKey;
    public $b2cConsumerSecret;
    public $passkey;
    public $shortcode;
    public $tillNumber;
    public $initiatorName;
    public $initiatorPassword;
    public $b2cShortcode;
    public $callbacks;

    public function __construct(array $config = [])
    {
        $this->environment = $config['environment'] ?? 'sandbox';
        $this->mpesaConsumerKey = $config['mpesa_consumer_key'] ?? '';
        $this->mpesaConsumerSecret = $config['mpesa_consumer_secret'] ?? '';
        $this->b2cConsumerKey = $config['b2c_consumer_key'] ?? null;
        $this->b2cConsumerSecret = $config['b2c_consumer_secret'] ?? null;
        $this->passkey = $config['passkey'] ?? '';
        $this->shortcode = $config['shortcode'] ?? '';
        $this->tillNumber = $config['till_number'] ?? null;
        $this->initiatorName = $config['initiator_name'] ?? '';
        $this->initiatorPassword = $config['initiator_password'] ?? '';
        $this->b2cShortcode = $config['b2c_shortcode'] ?? null;
        $this->callbacks = $config['callbacks'] ?? [];
    }

    public function getBaseUrl()
    {
        return $this->environment === 'sandbox'
            ? 'https://sandbox.safaricom.co.ke'
            : 'https://api.safaricom.co.ke';
    }
}