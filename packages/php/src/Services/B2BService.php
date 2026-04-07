<?php

namespace Yourdudeken\Mpesa\Services;

use Yourdudeken\Mpesa\Http\Client;
use Yourdudeken\Mpesa\Helpers\Signature;

class B2BService
{
    private $client;
    private $signature;

    public function __construct(Client $client, Signature $signature)
    {
        $this->client = $client;
        $this->signature = $signature;
    }

    public function send($receiver_shortcode, $command_id, $amount, $remarks, $account_number = null, $result_url = null, $timeout_url = null, $shortCodeType = 'B2B')
    {
        $url = $this->client->getBaseUrl() . '/mpesa/b2b/v1/paymentrequest';

        $body = [
            'Initiator' => config('mpesa.initiator_name'),
            'SecurityCredential' => $this->signature->generateSecurityCredential(),
            'CommandID' => $command_id,
            'SenderIdentifierType' => '4',
            'RecieverIdentifierType' => '4',
            'Amount' => $amount,
            'PartyA' => config('mpesa.b2c_shortcode'),
            'PartyB' => $receiver_shortcode,
            'AccountReference' => $account_number,
            'Remarks' => $remarks,
            'ResultURL' => $this->resolveCallbackUrl($result_url, 'b2b_result_url'),
            'QueueTimeOutURL' => $this->resolveCallbackUrl($timeout_url, 'b2b_timeout_url'),
        ];

        if ($command_id == 'BusinessPayBill' && $account_number == null) {
            throw new \Exception('Account Number is required for BusinessPayBill CommandID');
        }

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