<?php

namespace Yourdudeken\Mpesa\Services;

use Yourdudeken\Mpesa\Http\Client;
use Yourdudeken\Mpesa\Helpers\Signature;

class B2CService
{
    private $client;
    private $signature;

    public function __construct(Client $client, Signature $signature)
    {
        $this->client = $client;
        $this->signature = $signature;
    }

    public function send($phonenumber, $command_id, $amount, $remarks, $result_url = null, $timeout_url = null, $shortCodeType = 'B2C')
    {
        $url = $this->client->getBaseUrl() . '/mpesa/b2c/v1/paymentrequest';

        $body = [
            'InitiatorName' => config('mpesa.initiator_name'),
            'SecurityCredential' => $this->signature->generateSecurityCredential(),
            'CommandID' => $command_id,
            'Amount' => $amount,
            'PartyA' => config('mpesa.b2c_shortcode'),
            'PartyB' => $this->signature->phoneValidator($phonenumber),
            'Remarks' => $remarks,
            'Occassion' => '',
            'ResultURL' => $this->resolveCallbackUrl($result_url, 'b2c_result_url'),
            'QueueTimeOutURL' => $this->resolveCallbackUrl($timeout_url, 'b2c_timeout_url'),
        ];

        return $this->client->post($url, $body, $shortCodeType);
    }

    public function validated($phonenumber, $command_id, $amount, $remarks, $id_number, $result_url = null, $timeout_url = null, $shortCodeType = 'B2C')
    {
        $url = $this->client->getBaseUrl() . '/mpesa/b2cvalidate/v2/paymentrequest';

        $body = [
            'InitiatorName' => config('mpesa.initiator_name'),
            'SecurityCredential' => $this->signature->generateSecurityCredential(),
            'CommandID' => $command_id,
            'Amount' => $amount,
            'PartyA' => config('mpesa.b2c_shortcode'),
            'PartyB' => $this->signature->phoneValidator($phonenumber),
            'Remarks' => $remarks,
            'Occassion' => '',
            'OriginatorConversationID' => $this->signature->getFormattedTimestamp(),
            'IDType' => '01',
            'IDNumber' => $id_number,
            'ResultURL' => $this->resolveCallbackUrl($result_url, 'b2c_result_url'),
            'QueueTimeOutURL' => $this->resolveCallbackUrl($timeout_url, 'b2c_timeout_url'),
        ];

        return $this->client->post($url, $body, $shortCodeType);
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