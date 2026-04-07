<?php

namespace Yourdudeken\Mpesa\Services;

use Yourdudeken\Mpesa\Http\Client;
use Yourdudeken\Mpesa\Helpers\Signature;

class AccountService
{
    private $client;
    private $signature;

    public function __construct(Client $client, Signature $signature)
    {
        $this->client = $client;
        $this->signature = $signature;
    }

    public function balance($shortcode, $identifierType, $remarks, $result_url = null, $timeout_url = null, $shortCodeType = 'C2B')
    {
        $url = $this->client->getBaseUrl() . '/mpesa/accountbalance/v1/query';

        $body = [
            'Initiator' => config('mpesa.initiator_name'),
            'SecurityCredential' => $this->signature->generateSecurityCredential(),
            'CommandID' => 'AccountBalance',
            'PartyA' => $shortcode,
            'IdentifierType' => $identifierType,
            'Remarks' => $remarks,
            'ResultURL' => $this->resolveCallbackUrl($result_url, 'balance_result_url'),
            'QueueTimeOutURL' => $this->resolveCallbackUrl($timeout_url, 'balance_timeout_url'),
        ];

        return $this->client->post($url, $body, $shortCodeType);
    }

    public function status($shortcode, $transactionId, $identifierType, $remarks, $result_url = null, $timeout_url = null, $shortCodeType = 'C2B')
    {
        $url = $this->client->getBaseUrl() . '/mpesa/transactionstatus/v1/query';

        $body = [
            'Initiator' => config('mpesa.initiator_name'),
            'SecurityCredential' => $this->signature->generateSecurityCredential(),
            'CommandID' => 'TransactionStatusQuery',
            'TransactionID' => $transactionId,
            'PartyA' => $shortcode,
            'IdentifierType' => $identifierType,
            'Remarks' => $remarks,
            'Occassion' => '',
            'ResultURL' => $this->resolveCallbackUrl($result_url, 'status_result_url'),
            'QueueTimeOutURL' => $this->resolveCallbackUrl($timeout_url, 'status_timeout_url'),
        ];

        return $this->client->post($url, $body, $shortCodeType);
    }

    public function reversal($shortcode, $transactionId, $amount, $remarks, $result_url = null, $timeout_url = null, $shortCodeType = 'C2B')
    {
        $url = $this->client->getBaseUrl() . '/mpesa/reversal/v1/request';

        $body = [
            'Initiator' => config('mpesa.initiator_name'),
            'SecurityCredential' => $this->signature->generateSecurityCredential(),
            'CommandID' => 'TransactionReversal',
            'TransactionID' => $transactionId,
            'Amount' => $amount,
            'ReceiverParty' => $shortcode,
            'RecieverIdentifierType' => '11',
            'Remarks' => $remarks,
            'Occasion' => '',
            'ResultURL' => $this->resolveCallbackUrl($result_url, 'reversal_result_url'),
            'QueueTimeOutURL' => $this->resolveCallbackUrl($timeout_url, 'reversal_timeout_url'),
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