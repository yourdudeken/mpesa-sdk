<?php

namespace Yourdudeken\Mpesa\Tests\Feature;

use Illuminate\Support\Facades\Http;
use Yourdudeken\Mpesa\Facades\Mpesa;

it('can initiate b2pochi payment', function () {

    $expectedResponse = [
        'ConversationID' => 'AG_20240706_20106e9209f64bebd05b',
        'OriginatorConversationID' => '600997_Test_32et3241ed8yu',
        'ResponseCode' => '0',
        'ResponseDescription' => 'Accept the service request successfully.'
    ];

    Http::fake([
        'https://sandbox.safaricom.co.ke/mpesa/b2pochi/v1/paymentrequest' => Http::response($expectedResponse),
    ]);

    $response = Mpesa::b2pochi(
        '0707070707',
        10,
        'Payment to Pochi',
        'ChristmasPay',
        'https://test.test/result',
        'https://test.test/timeout'
    );

    $result = $response->json();

    expect($response->status())->toBe(200);
    expect($result)->toBe($expectedResponse);
    expect($result['ResponseCode'])->toBe('0');
});

it('can initiate b2pochi with callbacks from config', function () {

    $expectedResponse = [
        'ConversationID' => 'AG_20240706_20106e9209f64bebd05b',
        'OriginatorConversationID' => '600997_Test_32et3241ed8yu',
        'ResponseCode' => '0',
        'ResponseDescription' => 'Accept the service request successfully.'
    ];

    Http::fake([
        'https://sandbox.safaricom.co.ke/mpesa/b2pochi/v1/paymentrequest' => Http::response($expectedResponse),
    ]);

    config()->set('mpesa.callbacks.b2pochi_result_url', 'https://test.test/result');
    config()->set('mpesa.callbacks.b2pochi_timeout_url', 'https://test.test/timeout');

    $response = Mpesa::b2pochi(
        '0707070707',
        10,
        'Payment to Pochi'
    );

    $result = $response->json();

    expect($response->status())->toBe(200);
    expect($result)->toBe($expectedResponse);
    expect($result['ResponseCode'])->toBe('0');
});
