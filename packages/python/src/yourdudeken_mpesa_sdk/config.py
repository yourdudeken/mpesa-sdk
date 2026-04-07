class MpesaConfig:
    def __init__(self, environment: str, mpesa_consumer_key: str, mpesa_consumer_secret: str,
                 passkey: str, shortcode: str, initiator_name: str, initiator_password: str,
                 b2c_consumer_key: str = None, b2c_consumer_secret: str = None,
                 till_number: str = None, b2c_shortcode: str = None,
                 callbacks: dict = None):
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