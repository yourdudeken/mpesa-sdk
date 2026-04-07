import pytest
from mpesa import Mpesa, MpesaConfig


def test_mpesa_config():
    config = MpesaConfig(
        environment='sandbox',
        mpesa_consumer_key='test_key',
        mpesa_consumer_secret='test_secret',
        passkey='test_passkey',
        shortcode='174379',
        initiator_name='testapi',
        initiator_password='test_password'
    )
    assert config.environment == 'sandbox'
    assert config.mpesa_consumer_key == 'test_key'


def test_mpesa_initialization():
    config = MpesaConfig(
        environment='sandbox',
        mpesa_consumer_key='test_key',
        mpesa_consumer_secret='test_secret',
        passkey='test_passkey',
        shortcode='174379',
        initiator_name='testapi',
        initiator_password='test_password'
    )
    mpesa = Mpesa(config)
    assert mpesa is not None
