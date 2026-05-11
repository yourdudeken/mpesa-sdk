import pytest
from unittest.mock import Mock, patch, MagicMock
from mpesa import Mpesa, MpesaConfig


@pytest.fixture
def config():
    return MpesaConfig(
        environment='sandbox',
        mpesa_consumer_key='test_key',
        mpesa_consumer_secret='test_secret',
        passkey='test_passkey',
        shortcode='174379',
        initiator_name='testapi',
        initiator_password='test_password',
        b2c_shortcode='600000',
        callbacks={
            'callback_url': 'https://test.com/callback',
            'b2c_result_url': 'https://test.com/b2c_result',
            'b2c_timeout_url': 'https://test.com/b2c_timeout',
            'b2b_result_url': 'https://test.com/b2b_result',
            'b2b_timeout_url': 'https://test.com/b2b_timeout',
            'b2pochi_result_url': 'https://test.com/b2pochi_result',
            'b2pochi_timeout_url': 'https://test.com/b2pochi_timeout',
            'c2b_validation_url': 'https://test.com/c2b_validate',
            'c2b_confirmation_url': 'https://test.com/c2b_confirm',
            'balance_result_url': 'https://test.com/balance_result',
            'balance_timeout_url': 'https://test.com/balance_timeout',
            'status_result_url': 'https://test.com/status_result',
            'status_timeout_url': 'https://test.com/status_timeout',
            'reversal_result_url': 'https://test.com/reversal_result',
            'reversal_timeout_url': 'https://test.com/reversal_timeout',
        }
    )


class TestConfiguration:
    def test_mpesa_config(self):
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

    def test_mpesa_initialization(self, config):
        mpesa = Mpesa(config)
        assert mpesa is not None

    def test_static_constants(self):
        assert Mpesa.PAYBILL == 'CustomerPayBillOnline'
        assert Mpesa.TILL == 'CustomerBuyGoodsOnline'


class TestTokenGeneration:
    @patch('mpesa.client.requests.get')
    def test_generate_access_token(self, mock_get, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'access_token': 'test_token_123',
            'expires_in': '3599'
        }
        mock_response.raise_for_status = Mock()
        mock_get.return_value = mock_response

        mpesa = Mpesa(config)
        token = mpesa._generate_access_token('C2B')
        assert token == 'test_token_123'

    @patch('mpesa.client.requests.get')
    def test_use_b2c_credentials_for_b2c_requests(self, mock_get, config):
        config.b2c_consumer_key = 'b2c_key'
        config.b2c_consumer_secret = 'b2c_secret'

        mock_response = Mock()
        mock_response.json.return_value = {
            'access_token': 'b2c_token',
            'expires_in': '3599'
        }
        mock_response.raise_for_status = Mock()
        mock_get.return_value = mock_response

        mpesa = Mpesa(config)
        token = mpesa._generate_access_token('B2C')
        assert token == 'b2c_token'


class TestSTKPush:
    @patch('mpesa.client.requests.post')
    @patch('mpesa.client.requests.get')
    def test_initiate_stk_push_with_callback_url(self, mock_get, mock_post, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'MerchantRequestID': '29115-34620561-1',
            'CheckoutRequestID': 'ws_CO_191220191020363925',
            'ResponseCode': '0',
            'ResponseDescription': 'Success',
            'CustomerMessage': 'Success'
        }
        mock_response.raise_for_status = Mock()
        mock_post.return_value = mock_response

        mock_token_response = Mock()
        mock_token_response.json.return_value = {'access_token': 'test', 'expires_in': '3599'}
        mock_token_response.raise_for_status = Mock()
        mock_get.return_value = mock_token_response

        mpesa = Mpesa(config)
        response = mpesa.stkpush('254712345678', 100, '12345', 'https://test.com/callback')
        assert response['ResponseCode'] == '0'

    def test_throw_error_when_account_reference_missing(self, config):
        mpesa = Mpesa(config)
        with pytest.raises(ValueError, match='Account Reference is required'):
            mpesa.stkpush('254712345678', 100, '')

    def test_require_till_number_for_till_transactions(self, config):
        mpesa = Mpesa(config)
        with pytest.raises(ValueError, match='Till number is required'):
            mpesa.stkpush('254712345678', 100, '12345', transaction_type=Mpesa.TILL)

    @patch('mpesa.client.requests.post')
    @patch('mpesa.client.requests.get')
    def test_query_stk_push_status(self, mock_get, mock_post, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'ResponseCode': '0',
            'ResponseDescription': 'Success',
            'MerchantRequestID': '22205-34066-1',
            'CheckoutRequestID': 'ws_CO_13012021093521236557',
            'ResultCode': '0',
            'ResultDesc': 'Success'
        }
        mock_response.raise_for_status = Mock()
        mock_post.return_value = mock_response

        mock_token_response = Mock()
        mock_token_response.json.return_value = {'access_token': 'test', 'expires_in': '3599'}
        mock_token_response.raise_for_status = Mock()
        mock_get.return_value = mock_token_response

        mpesa = Mpesa(config)
        response = mpesa.stkquery('ws_CO_191220191020363925')
        assert response['ResponseCode'] == '0'


class TestB2C:
    @patch('mpesa.client.requests.post')
    @patch('mpesa.client.requests.get')
    def test_send_b2c_payment(self, mock_get, mock_post, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'ConversationID': 'AG_20231217_201020363925',
            'OriginatorConversationID': '201020363925',
            'ResponseCode': '0',
            'ResponseDescription': 'Success'
        }
        mock_response.raise_for_status = Mock()
        mock_post.return_value = mock_response

        mock_token_response = Mock()
        mock_token_response.json.return_value = {'access_token': 'test', 'expires_in': '3599'}
        mock_token_response.raise_for_status = Mock()
        mock_get.return_value = mock_token_response

        mpesa = Mpesa(config)
        response = mpesa.b2c('254712345678', 'BusinessPayment', 100, 'Test payment')
        assert response['ResponseCode'] == '0'

    @patch('mpesa.client.requests.post')
    @patch('mpesa.client.requests.get')
    def test_send_validated_b2c_payment(self, mock_get, mock_post, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'ConversationID': 'AG_20231217_201020363925',
            'ResponseCode': '0'
        }
        mock_response.raise_for_status = Mock()
        mock_post.return_value = mock_response

        mock_token_response = Mock()
        mock_token_response.json.return_value = {'access_token': 'test', 'expires_in': '3599'}
        mock_token_response.raise_for_status = Mock()
        mock_get.return_value = mock_token_response

        mpesa = Mpesa(config)
        response = mpesa.validated_b2c('254712345678', 'BusinessPayment', 100, 'Test payment', '12345678')
        assert response['ResponseCode'] == '0'


class TestB2B:
    @patch('mpesa.client.requests.post')
    @patch('mpesa.client.requests.get')
    def test_send_b2b_payment(self, mock_get, mock_post, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'ConversationID': 'AG_20231217_201020363925',
            'ResponseCode': '0'
        }
        mock_response.raise_for_status = Mock()
        mock_post.return_value = mock_response

        mock_token_response = Mock()
        mock_token_response.json.return_value = {'access_token': 'test', 'expires_in': '3599'}
        mock_token_response.raise_for_status = Mock()
        mock_get.return_value = mock_token_response

        mpesa = Mpesa(config)
        response = mpesa.b2b('600000', 'BusinessPayBill', 100, 'Test payment', '12345')
        assert response['ResponseCode'] == '0'

    def test_throw_error_when_account_number_missing_for_business_pay_bill(self, config):
        mpesa = Mpesa(config)
        with pytest.raises(ValueError, match='Account Number is required'):
            mpesa.b2b('600000', 'BusinessPayBill', 100, 'Test payment')


class TestC2B:
    @patch('mpesa.client.requests.post')
    @patch('mpesa.client.requests.get')
    def test_register_c2b_urls(self, mock_get, mock_post, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'ResponseCode': '0',
            'ResponseDescription': 'success'
        }
        mock_response.raise_for_status = Mock()
        mock_post.return_value = mock_response

        mock_token_response = Mock()
        mock_token_response.json.return_value = {'access_token': 'test', 'expires_in': '3599'}
        mock_token_response.raise_for_status = Mock()
        mock_get.return_value = mock_token_response

        mpesa = Mpesa(config)
        response = mpesa.c2b_register_urls('600000')
        assert response['ResponseCode'] == '0'

    @patch('mpesa.client.requests.post')
    @patch('mpesa.client.requests.get')
    def test_simulate_c2b_payment(self, mock_get, mock_post, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'ResponseCode': '0',
            'ResponseDescription': 'Success'
        }
        mock_response.raise_for_status = Mock()
        mock_post.return_value = mock_response

        mock_token_response = Mock()
        mock_token_response.json.return_value = {'access_token': 'test', 'expires_in': '3599'}
        mock_token_response.raise_for_status = Mock()
        mock_get.return_value = mock_token_response

        mpesa = Mpesa(config)
        response = mpesa.c2bsimulate('254712345678', 100, '600000', 'CustomerPayBillOnline')
        assert response['ResponseCode'] == '0'


class TestAccountBalance:
    @patch('mpesa.client.requests.post')
    @patch('mpesa.client.requests.get')
    def test_query_account_balance(self, mock_get, mock_post, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'ResponseCode': '0',
            'ResponseDescription': 'Success'
        }
        mock_response.raise_for_status = Mock()
        mock_post.return_value = mock_response

        mock_token_response = Mock()
        mock_token_response.json.return_value = {'access_token': 'test', 'expires_in': '3599'}
        mock_token_response.raise_for_status = Mock()
        mock_get.return_value = mock_token_response

        mpesa = Mpesa(config)
        response = mpesa.account_balance('600000', 4, 'Check balance')
        assert response['ResponseCode'] == '0'


class TestTransactionStatus:
    @patch('mpesa.client.requests.post')
    @patch('mpesa.client.requests.get')
    def test_query_transaction_status(self, mock_get, mock_post, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'ResponseCode': '0',
            'ResponseDescription': 'Success'
        }
        mock_response.raise_for_status = Mock()
        mock_post.return_value = mock_response

        mock_token_response = Mock()
        mock_token_response.json.return_value = {'access_token': 'test', 'expires_in': '3599'}
        mock_token_response.raise_for_status = Mock()
        mock_get.return_value = mock_token_response

        mpesa = Mpesa(config)
        response = mpesa.transaction_status('600000', '123456789', 1, 'Check status')
        assert response['ResponseCode'] == '0'


class TestReversal:
    @patch('mpesa.client.requests.post')
    @patch('mpesa.client.requests.get')
    def test_reverse_transaction(self, mock_get, mock_post, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'ResponseCode': '0',
            'ResponseDescription': 'Success'
        }
        mock_response.raise_for_status = Mock()
        mock_post.return_value = mock_response

        mock_token_response = Mock()
        mock_token_response.json.return_value = {'access_token': 'test', 'expires_in': '3599'}
        mock_token_response.raise_for_status = Mock()
        mock_get.return_value = mock_token_response

        mpesa = Mpesa(config)
        response = mpesa.reversal('600000', '123456789', 100, 'Reverse transaction')
        assert response['ResponseCode'] == '0'


class TestB2Pochi:
    @patch('mpesa.client.requests.post')
    @patch('mpesa.client.requests.get')
    def test_send_b2pochi_payment(self, mock_get, mock_post, config):
        mock_response = Mock()
        mock_response.json.return_value = {
            'ConversationID': 'AG_20231217_201020363925',
            'ResponseCode': '0'
        }
        mock_response.raise_for_status = Mock()
        mock_post.return_value = mock_response

        mock_token_response = Mock()
        mock_token_response.json.return_value = {'access_token': 'test', 'expires_in': '3599'}
        mock_token_response.raise_for_status = Mock()
        mock_get.return_value = mock_token_response

        mpesa = Mpesa(config)
        response = mpesa.b2pochi('254712345678', 100, 'Pochi payment')
        assert response['ResponseCode'] == '0'


class TestPhoneValidation:
    def test_validate_phone_numbers_correctly(self, config):
        mpesa = Mpesa(config)
        
        assert mpesa._phone_validator('+254712345678') == '254712345678'
        assert mpesa._phone_validator('0712345678') == '254712345678'
        assert mpesa._phone_validator('712345678') == '254712345678'