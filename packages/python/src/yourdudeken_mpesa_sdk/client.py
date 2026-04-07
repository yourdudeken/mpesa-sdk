import base64
import time
import requests
from requests.auth import HTTPBasicAuth
from typing import Optional, Dict, Any
from datetime import datetime


class MpesaConfig:
    def __init__(self, environment: str, mpesa_consumer_key: str, mpesa_consumer_secret: str,
                 passkey: str, shortcode: str, initiator_name: str, initiator_password: str,
                 b2c_consumer_key: Optional[str] = None, b2c_consumer_secret: Optional[str] = None,
                 till_number: Optional[str] = None, b2c_shortcode: Optional[str] = None,
                 callbacks: Optional[Dict[str, str]] = None):
        self.environment = environment
        self.mpesa_consumer_key = mpesa_consumer_key
        self.mpesa_consumer_secret = mpesa_consumer_secret
        self.b2c_consumer_key = b2c_consumer_key
        self.b2c_consumer_secret = b2c_consumer_secret
        self.passkey = passkey
        self.shortcode = shortcode
        self.till_number = till_number
        self.initiator_name = initiator_name
        self.initiator_password = initiator_password
        self.b2c_shortcode = b2c_shortcode
        self.callbacks = callbacks or {}


class Mpesa:
    PAYBILL = 'CustomerPayBillOnline'
    TILL = 'CustomerBuyGoodsOnline'
    
    def __init__(self, config: MpesaConfig):
        self.config = config
        self.base_url = 'https://sandbox.safaricom.co.ke' if config.environment == 'sandbox' else 'https://api.safaricom.co.ke'
        self._access_token: Optional[str] = None
        self._token_expiry: float = 0
    
    def _get_config(self, key: str, default=None):
        return getattr(self.config, key, default)
    
    def _resolve_callback_url(self, param_url: Optional[str], config_key: str) -> str:
        config_url = self.config.callbacks.get(config_key) if self.config.callbacks else None
        if param_url:
            return param_url
        if config_url:
            return config_url
        raise ValueError(f"Ensure you have set the {config_key} in the config or passed as a parameter")
    
    def _phone_validator(self, phone_no: str) -> str:
        phone_no = phone_no.lstrip('+')
        if phone_no.startswith('0'):
            phone_no = '254' + phone_no[1:]
        elif phone_no.startswith('7'):
            phone_no = '254' + phone_no
        return phone_no
    
    def _get_formatted_timestamp(self) -> str:
        return datetime.now().strftime('%Y%m%d%H%M%S')
    
    def _lipa_na_mpesa_password(self) -> str:
        timestamp = self._get_formatted_timestamp()
        password = str(self._get_config('shortcode')) + self._get_config('passkey') + timestamp
        return base64.b64encode(password.encode()).decode()
    
    def _generate_access_token(self, short_code_type: str = 'C2B') -> str:
        if self._access_token and time.time() < self._token_expiry:
            return self._access_token
        
        consumer_key = self._get_config('b2c_consumer_key') if short_code_type in ['B2C', 'B2B'] else self._get_config('mpesa_consumer_key')
        consumer_secret = self._get_config('b2c_consumer_secret') if short_code_type in ['B2C', 'B2B'] else self._get_config('mpesa_consumer_secret')
        
        url = f'{self.base_url}/oauth/v1/generate?grant_type=client_credentials'
        
        response = requests.get(url, auth=HTTPBasicAuth(consumer_key, consumer_secret))
        response.raise_for_status()
        
        data = response.json()
        self._access_token = data.get('access_token')
        self._token_expiry = time.time() + (data.get('expires_in', 3600) - 60)
        
        return self._access_token
    
    def _generate_security_credential(self) -> str:
        import os
        from cryptography.hazmat.primitives import serialization
        from cryptography.hazmat.primitives.asymmetric import padding
        from cryptography.hazmat.backends import default_backend
        
        cert_path = os.path.join(os.path.dirname(__file__), '..', 'certificates', 
                                'SandboxCertificate.cer' if self._get_config('environment') == 'sandbox' else 'ProductionCertificate.cer')
        
        with open(cert_path, 'rb') as f:
            pubkey = f.read()
        
        public_key = serialization.load_pem_public_key(pubkey, default_backend())
        
        password = self._get_config('initiator_password').encode()
        encrypted = public_key.encrypt(password, padding.PKCS1v15())
        
        return base64.b64encode(encrypted).decode()
    
    def _mpesa_request(self, url: str, body: Dict[str, Any], short_code_type: str = 'C2B') -> Dict[str, Any]:
        token = self._generate_access_token(short_code_type)
        
        headers = {
            'Authorization': f'Bearer {token}',
            'Content-Type': 'application/json',
        }
        
        response = requests.post(url, json=body, headers=headers)
        response.raise_for_status()
        
        return response.json()
    
    def stkpush(self, phonenumber: str, amount: int, account_number: str, callback_url: Optional[str] = None, 
                transaction_type: str = PAYBILL, short_code_type: str = 'C2B') -> Dict[str, Any]:
        if not account_number:
            raise ValueError('An Account Reference is required for All transactions.')
        
        if transaction_type == self.TILL and not self._get_config('till_number'):
            raise ValueError('Till number is required for Buy Goods transactions.')
        
        url = f'{self.base_url}/mpesa/stkpush/v1/processrequest'
        
        data = {
            'BusinessShortCode': self._get_config('shortcode'),
            'Password': self._lipa_na_mpesa_password(),
            'Timestamp': self._get_formatted_timestamp(),
            'Amount': amount,
            'PartyA': self._phone_validator(phonenumber),
            'PartyB': self._get_config('shortcode') if transaction_type == self.PAYBILL else self._get_config('till_number'),
            'TransactionType': transaction_type,
            'PhoneNumber': self._phone_validator(phonenumber),
            'TransactionDesc': 'Payment',
            'AccountReference': account_number,
            'CallBackURL': self._resolve_callback_url(callback_url, 'callback_url'),
        }
        
        return self._mpesa_request(url, data, short_code_type)
    
    def stkquery(self, checkout_request_id: str, short_code_type: str = 'C2B') -> Dict[str, Any]:
        url = f'{self.base_url}/mpesa/stkpushquery/v1/query'
        
        data = {
            'BusinessShortCode': self._get_config('shortcode'),
            'Password': self._lipa_na_mpesa_password(),
            'Timestamp': self._get_formatted_timestamp(),
            'CheckoutRequestID': checkout_request_id,
        }
        
        return self._mpesa_request(url, data, short_code_type)
    
    def b2c(self, phonenumber: str, command_id: str, amount: int, remarks: str, 
            result_url: Optional[str] = None, timeout_url: Optional[str] = None, 
            short_code_type: str = 'B2C') -> Dict[str, Any]:
        url = f'{self.base_url}/mpesa/b2c/v1/paymentrequest'
        
        body = {
            'InitiatorName': self._get_config('initiator_name'),
            'SecurityCredential': self._generate_security_credential(),
            'CommandID': command_id,
            'Amount': amount,
            'PartyA': self._get_config('b2c_shortcode'),
            'PartyB': self._phone_validator(phonenumber),
            'Remarks': remarks,
            'Occassion': '',
            'ResultURL': self._resolve_callback_url(result_url, 'b2c_result_url'),
            'QueueTimeOutURL': self._resolve_callback_url(timeout_url, 'b2c_timeout_url'),
        }
        
        return self._mpesa_request(url, body, short_code_type)
    
    def validated_b2c(self, phonenumber: str, command_id: str, amount: int, remarks: str, 
                     id_number: str, result_url: Optional[str] = None, timeout_url: Optional[str] = None, 
                     short_code_type: str = 'B2C') -> Dict[str, Any]:
        url = f'{self.base_url}/mpesa/b2cvalidate/v2/paymentrequest'
        
        body = {
            'InitiatorName': self._get_config('initiator_name'),
            'SecurityCredential': self._generate_security_credential(),
            'CommandID': command_id,
            'Amount': amount,
            'PartyA': self._get_config('b2c_shortcode'),
            'PartyB': self._phone_validator(phonenumber),
            'Remarks': remarks,
            'Occassion': '',
            'OriginatorConversationID': self._get_formatted_timestamp(),
            'IDType': '01',
            'IDNumber': id_number,
            'ResultURL': self._resolve_callback_url(result_url, 'b2c_result_url'),
            'QueueTimeOutURL': self._resolve_callback_url(timeout_url, 'b2c_timeout_url'),
        }
        
        return self._mpesa_request(url, body, short_code_type)
    
    def b2b(self, receiver_shortcode: str, command_id: str, amount: int, remarks: str, 
           account_number: Optional[str] = None, result_url: Optional[str] = None, 
           timeout_url: Optional[str] = None, short_code_type: str = 'B2B') -> Dict[str, Any]:
        if command_id == 'BusinessPayBill' and not account_number:
            raise ValueError('Account Number is required for BusinessPayBill CommandID')
        
        url = f'{self.base_url}/mpesa/b2b/v1/paymentrequest'
        
        body = {
            'Initiator': self._get_config('initiator_name'),
            'SecurityCredential': self._generate_security_credential(),
            'CommandID': command_id,
            'SenderIdentifierType': '4',
            'RecieverIdentifierType': '4',
            'Amount': amount,
            'PartyA': self._get_config('b2c_shortcode'),
            'PartyB': receiver_shortcode,
            'AccountReference': account_number,
            'Remarks': remarks,
            'ResultURL': self._resolve_callback_url(result_url, 'b2b_result_url'),
            'QueueTimeOutURL': self._resolve_callback_url(timeout_url, 'b2b_timeout_url'),
        }
        
        return self._mpesa_request(url, body, short_code_type)
    
    def c2bregisterURLS(self, shortcode: str, confirm_url: Optional[str] = None, 
                       validate_url: Optional[str] = None, short_code_type: str = 'C2B') -> Dict[str, Any]:
        url = f'{self.base_url}/mpesa/c2b/v2/registerurl'
        
        body = {
            'ShortCode': shortcode,
            'ResponseType': 'Completed',
            'ConfirmationURL': self._resolve_callback_url(confirm_url, 'c2b_confirmation_url'),
            'ValidationURL': self._resolve_callback_url(validate_url, 'c2b_validation_url'),
        }
        
        return self._mpesa_request(url, body, short_code_type)
    
    def c2bsimulate(self, phonenumber: str, amount: int, shortcode: str, command_id: str, 
                   account_number: Optional[str] = None, short_code_type: str = 'C2B') -> Dict[str, Any]:
        url = f'{self.base_url}/mpesa/c2b/v2/simulate'
        
        if command_id == self.PAYBILL:
            data = {
                'Msisdn': self._phone_validator(phonenumber),
                'Amount': amount,
                'BillRefNumber': account_number,
                'CommandID': command_id,
                'ShortCode': shortcode,
            }
        else:
            data = {
                'Msisdn': self._phone_validator(phonenumber),
                'Amount': amount,
                'CommandID': command_id,
                'ShortCode': shortcode,
            }
        
        return self._mpesa_request(url, data, short_code_type)
    
    def transaction_status(self, shortcode: str, transaction_id: str, identifier_type: int, 
                          remarks: str, result_url: Optional[str] = None, 
                          timeout_url: Optional[str] = None, short_code_type: str = 'C2B') -> Dict[str, Any]:
        url = f'{self.base_url}/mpesa/transactionstatus/v1/query'
        
        body = {
            'Initiator': self._get_config('initiator_name'),
            'SecurityCredential': self._generate_security_credential(),
            'CommandID': 'TransactionStatusQuery',
            'TransactionID': transaction_id,
            'PartyA': shortcode,
            'IdentifierType': identifier_type,
            'Remarks': remarks,
            'Occassion': '',
            'ResultURL': self._resolve_callback_url(result_url, 'status_result_url'),
            'QueueTimeOutURL': self._resolve_callback_url(timeout_url, 'status_timeout_url'),
        }
        
        return self._mpesa_request(url, body, short_code_type)
    
    def account_balance(self, shortcode: str, identifier_type: int, remarks: str, 
                       result_url: Optional[str] = None, timeout_url: Optional[str] = None, 
                       short_code_type: str = 'C2B') -> Dict[str, Any]:
        url = f'{self.base_url}/mpesa/accountbalance/v1/query'
        
        body = {
            'Initiator': self._get_config('initiator_name'),
            'SecurityCredential': self._generate_security_credential(),
            'CommandID': 'AccountBalance',
            'PartyA': shortcode,
            'IdentifierType': identifier_type,
            'Remarks': remarks,
            'ResultURL': self._resolve_callback_url(result_url, 'balance_result_url'),
            'QueueTimeOutURL': self._resolve_callback_url(timeout_url, 'balance_timeout_url'),
        }
        
        return self._mpesa_request(url, body, short_code_type)
    
    def reversal(self, shortcode: str, transaction_id: str, amount: float, remarks: str, 
                result_url: Optional[str] = None, timeout_url: Optional[str] = None, 
                short_code_type: str = 'C2B') -> Dict[str, Any]:
        url = f'{self.base_url}/mpesa/reversal/v1/request'
        
        body = {
            'Initiator': self._get_config('initiator_name'),
            'SecurityCredential': self._generate_security_credential(),
            'CommandID': 'TransactionReversal',
            'TransactionID': transaction_id,
            'Amount': amount,
            'ReceiverParty': shortcode,
            'RecieverIdentifierType': '11',
            'Remarks': remarks,
            'Occasion': '',
            'ResultURL': self._resolve_callback_url(result_url, 'reversal_result_url'),
            'QueueTimeOutURL': self._resolve_callback_url(timeout_url, 'reversal_timeout_url'),
        }
        
        return self._mpesa_request(url, body, short_code_type)
    
    def b2pochi(self, phonenumber: str, amount: int, remarks: str, occasion: Optional[str] = None, 
               result_url: Optional[str] = None, timeout_url: Optional[str] = None, 
               short_code_type: str = 'B2C') -> Dict[str, Any]:
        url = f'{self.base_url}/mpesa/b2pochi/v1/paymentrequest'
        
        body = {
            'OriginatorConversationID': self._get_formatted_timestamp(),
            'InitiatorName': self._get_config('initiator_name'),
            'SecurityCredential': self._generate_security_credential(),
            'CommandID': 'BusinessPayToPochi',
            'Amount': amount,
            'PartyA': self._get_config('b2c_shortcode'),
            'PartyB': self._phone_validator(phonenumber),
            'Remarks': remarks,
            'Occassion': occasion or '',
            'ResultURL': self._resolve_callback_url(result_url, 'b2pochi_result_url'),
            'QueueTimeOutURL': self._resolve_callback_url(timeout_url, 'b2pochi_timeout_url'),
        }
        
        return self._mpesa_request(url, body, short_code_type)