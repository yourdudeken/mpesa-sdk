<?php

namespace Yourdudeken\Mpesa\Http;

use Illuminate\Support\Facades\Http;

class Client
{
    private $baseUrl;

    public function __construct()
    {
        $this->baseUrl = config('mpesa.environment') == 'sandbox'
            ? 'https://sandbox.safaricom.co.ke'
            : 'https://api.safaricom.co.ke';
    }

    public function getBaseUrl()
    {
        return $this->baseUrl;
    }

    public function post($url, $body, $shortCodeType = 'C2B')
    {
        $token = $this->generateAccessToken($shortCodeType);

        return Http::withToken($token)
            ->acceptJson()
            ->post($url, $body);
    }

    public function generateAccessToken($shortCodeType)
    {
        if ($shortCodeType == 'B2C' || $shortCodeType == 'B2B') {
            $consumer_key = config('mpesa.b2c_consumer_key');
            $consumer_secret = config('mpesa.b2c_consumer_secret');
        } else {
            $consumer_key = config('mpesa.mpesa_consumer_key');
            $consumer_secret = config('mpesa.mpesa_consumer_secret');
        }

        $url = $this->baseUrl . '/oauth/v1/generate?grant_type=client_credentials';

        $response = Http::withBasicAuth($consumer_key, $consumer_secret)
            ->get($url);

        $result = json_decode($response);

        return data_get($result, 'access_token');
    }
}