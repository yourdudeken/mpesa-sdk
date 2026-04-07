<?php

namespace Yourdudeken\Mpesa\Tests\Unit;

use Yourdudeken\Mpesa\Mpesa;

it('can generate security credential', function () {

    config()->set('mpesa.initiator_password', 'test');

    $mpesa = new Mpesa();

    $result = $mpesa->generate_security_credential();

    expect($result)->toContain('==');
});
