<?php

namespace Yourdudeken\Mpesa\Helpers;

use Illuminate\Support\Facades\File;
use Illuminate\Support\Carbon;

class Signature
{
    public function lipaNaMpesaPassword()
    {
        $timestamp = $this->getFormattedTimestamp();
        return base64_encode(config('mpesa.shortcode') . config('mpesa.passkey') . $timestamp);
    }

    public function getFormattedTimestamp()
    {
        return Carbon::rawParse('now')->format('YmdHis');
    }

    public function phoneValidator($phoneno)
    {
        $phoneno = (substr($phoneno, 0, 1) == '+') ? str_replace('+', '', $phoneno) : $phoneno;
        $phoneno = (substr($phoneno, 0, 1) == '0') ? preg_replace('/^0/', '254', $phoneno) : $phoneno;
        $phoneno = (substr($phoneno, 0, 1) == '7') ? "254{$phoneno}" : $phoneno;

        return $phoneno;
    }

    public function generateSecurityCredential()
    {
        if (config('mpesa.environment') == 'sandbox') {
            $pubkey = File::get(__DIR__ . '/../certificates/SandboxCertificate.cer');
        } else {
            $pubkey = File::get(__DIR__ . '/../certificates/ProductionCertificate.cer');
        }

        openssl_public_encrypt(config('mpesa.initiator_password'), $output, $pubkey, OPENSSL_PKCS1_PADDING);

        return base64_encode($output);
    }
}