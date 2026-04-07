from typing import Optional
from .http.http_client import HttpClient
from .auth import Auth
from .utils import Helpers


class B2PochiService:
    def __init__(self, http_client: HttpClient, auth: Auth, config):
        self.http_client = http_client
        self.auth = auth
        self.helpers = Helpers(config)
        self.base_url = http_client.get_base_url()

    def send(self, phonenumber: str, amount: int, remarks: str,
             occasion: Optional[str] = None,
             result_url: Optional[str] = None, timeout_url: Optional[str] = None,
             short_code_type: str = 'B2C') -> dict:
        url = f'{self.base_url}/mpesa/b2pochi/v1/paymentrequest'
        
        body = {
            'OriginatorConversationID': self.helpers.get_formatted_timestamp(),
            'InitiatorName': self.helpers.get_config('initiator_name'),
            'SecurityCredential': self.helpers.generate_security_credential(),
            'CommandID': 'BusinessPayToPochi',
            'Amount': amount,
            'PartyA': self.helpers.get_config('b2c_shortcode'),
            'PartyB': self.helpers.phone_validator(phonenumber),
            'Remarks': remarks,
            'Occasion': occasion or '',
            'ResultURL': result_url or self.helpers.get_config('callbacks.b2pochi_result_url'),
            'QueueTimeOutURL': timeout_url or self.helpers.get_config('callbacks.b2pochi_timeout_url'),
        }
        
        token = self.auth.get_access_token(short_code_type)
        return self.http_client.post(url, body, token)