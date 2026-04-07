<?php

namespace Yourdudeken\Mpesa\Services;

use Yourdudeken\Mpesa\Http\Client;
use Yourdudeken\Mpesa\Helpers\Signature;

class STKPushService
{
    private $client;
    private $signature;

    public function __construct(Client $client, Signature $signature)
    {
        $this->client = $client;
        $this->signature = $signature;
    }

    public function push($phonenumber, $amount, $account_number, $callbackurl = null, $transactionType = 'CustomerPayBillOnline', $shortCodeType = 'C2B')
    {
        $url = $this->client->getBaseUrl() . '/mpesa/stkpush/v1/processrequest';
        
        $data = [
            'BusinessShortCode' => config('mpesa.shortcode'),
            'Password' => $this->signature->lipaNaMpesaPassword(),
            'Timestamp' => $this->signature->getFormattedTimestamp(),
            'Amount' => (int) $amount,
            'PartyA' => $this->signature->phoneValidator($phonenumber),
            'PartyB' => $transactionType == 'CustomerPayBillOnline' 
                ? config('mpesa.shortcode') 
                : config('mpesa.till_number'),
            'TransactionType' => $transactionType,
            'PhoneNumber' => $this->signature->phoneValidator($phonenumber),
            'TransactionDesc' => 'Payment',
            'AccountReference' => $account_number,
            'CallBackURL' => $this->resolveCallbackUrl($callbackurl, 'callback_url'),
        ];

        return $this->client->post($url, $data, $shortCodeType);
    }

    public function query($checkoutRequestId, $shortCodeType = 'C2B')
    {
        $url = $this->client->getBaseUrl() . '/mpesa/stkpushquery/v1/query';
        
        $data = [
            'BusinessShortCode' => config('mpesa.shortcode'),
            'Password' => $this->signature->lipaNaMpesaPassword(),
            'Timestamp' => $this->signature->getFormattedTimestamp(),
            'CheckoutRequestID' => $checkoutRequestId,
        ];

        return $this->client->post($url, $data, $shortCodeType);
    }

    private function resolveCallbackUrl($paramUrl, $configKey)
    {
        $configUrl = config("mpesa.callbacks.{$configKey}");
        
        if ($paramUrl !== null) {
            return $paramUrl;
        } elseif ($configUrl !== null) {
            return $configUrl;
        }
        
        throw new \Exception("Ensure you have set the {$configKey} in the config or passed as a parameter");
    }
}