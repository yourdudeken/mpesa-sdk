<?php

namespace Yourdudeken\Mpesa\Services;

use Yourdudeken\Mpesa\Http\Client;

class C2BService
{
    private $client;

    public function __construct(Client $client)
    {
        $this->client = $client;
    }

    public function registerURLS($shortcode, $confirmurl = null, $validateurl = null, $shortCodeType = 'C2B')
    {
        $url = $this->client->getBaseUrl() . '/mpesa/c2b/v2/registerurl';

        $body = [
            'ShortCode' => $shortcode,
            'ResponseType' => 'Completed',
            'ConfirmationURL' => $this->resolveCallbackUrl($confirmurl, 'c2b_confirmation_url'),
            'ValidationURL' => $this->resolveCallbackUrl($validateurl, 'c2b_validation_url'),
        ];

        return $this->client->post($url, $body, $shortCodeType);
    }

    public function simulate($phonenumber, $amount, $shortcode, $command_id, $account_number = null, $shortCodeType = 'C2B')
    {
        $url = $this->client->getBaseUrl() . '/mpesa/c2b/v2/simulate';

        if ($command_id == 'CustomerPayBillOnline') {
            $data = [
                'Msisdn' => $phonenumber,
                'Amount' => (int) $amount,
                'BillRefNumber' => $account_number,
                'CommandID' => $command_id,
                'ShortCode' => $shortcode,
            ];
        } else {
            $data = [
                'Msisdn' => $phonenumber,
                'Amount' => (int) $amount,
                'CommandID' => $command_id,
                'ShortCode' => $shortcode,
            ];
        }

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