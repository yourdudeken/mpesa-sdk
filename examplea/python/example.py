import os
from mpesa import Mpesa, MpesaConfig


config = MpesaConfig(
    environment=os.getenv('MPESA_ENV', 'sandbox'),
    mpesa_consumer_key=os.getenv('MPESA_CONSUMER_KEY'),
    mpesa_consumer_secret=os.getenv('MPESA_CONSUMER_SECRET'),
    passkey=os.getenv('MPESA_PASSKEY'),
    shortcode=os.getenv('MPESA_SHORTCODE', '174379'),
    initiator_name=os.getenv('MPESA_INITIATOR_NAME', 'testapi'),
    initiator_password=os.getenv('MPESA_INITIATOR_PASSWORD'),
    b2c_shortcode=os.getenv('MPESA_B2C_SHORTCODE'),
    till_number=os.getenv('MPESA_TILL_NUMBER'),
    b2c_consumer_key=os.getenv('MPESA_B2C_CONSUMER_KEY'),
    b2c_consumer_secret=os.getenv('MPESA_B2C_CONSUMER_SECRET'),
    callbacks={
        'callback_url': os.getenv('MPESA_CALLBACK_URL'),
        'b2c_result_url': os.getenv('MPESA_B2C_RESULT_URL'),
        'b2c_timeout_url': os.getenv('MPESA_B2C_TIMEOUT_URL'),
        'b2b_result_url': os.getenv('MPESA_B2B_RESULT_URL'),
        'b2b_timeout_url': os.getenv('MPESA_B2B_TIMEOUT_URL'),
        'b2pochi_result_url': os.getenv('MPESA_B2POCHI_RESULT_URL'),
        'b2pochi_timeout_url': os.getenv('MPESA_B2POCHI_TIMEOUT_URL'),
        'c2b_validation_url': os.getenv('MPESA_C2B_VALIDATION_URL'),
        'c2b_confirmation_url': os.getenv('MPESA_C2B_CONFIRMATION_URL'),
        'balance_result_url': os.getenv('MPESA_BALANCE_RESULT_URL'),
        'balance_timeout_url': os.getenv('MPESA_BALANCE_TIMEOUT_URL'),
        'status_result_url': os.getenv('MPESA_STATUS_RESULT_URL'),
        'status_timeout_url': os.getenv('MPESA_STATUS_TIMEOUT_URL'),
        'reversal_result_url': os.getenv('MPESA_REVERSAL_RESULT_URL'),
        'reversal_timeout_url': os.getenv('MPESA_REVERSAL_TIMEOUT_URL'),
    },
)

mpesa = Mpesa(config)


def main():
    try:
        # STK Push – Lipa na Mpesa Online
        stk_response = mpesa.stkpush(
            phonenumber='254712345678',
            amount=10,
            account_number='INV-001',
        )
        print('STK Push:', stk_response)

        # Query STK Push status
        status_response = mpesa.stkquery(stk_response['CheckoutRequestID'])
        print('STK Query:', status_response)

        # B2C – Business to Customer
        b2c_response = mpesa.b2c(
            phonenumber='254712345678',
            command_id='BusinessPayment',
            amount=500,
            remarks='Salary payment',
        )
        print('B2C:', b2c_response)

        # B2B – Business to Business
        b2b_response = mpesa.b2b(
            receiver_shortcode='600000',
            command_id='BusinessPayBill',
            amount=1000,
            remarks='Invoice payment',
            account_number='INV-001',
        )
        print('B2B:', b2b_response)

        # C2B – Register URLs
        register_response = mpesa.c2b_register_urls(
            shortcode=os.getenv('MPESA_SHORTCODE', '174379'),
        )
        print('C2B Register:', register_response)

        # C2B – Simulate payment
        simulate_response = mpesa.c2bsimulate(
            phonenumber='254712345678',
            amount=100,
            shortcode=os.getenv('MPESA_SHORTCODE', '174379'),
            command_id=Mpesa.PAYBILL,
        )
        print('C2B Simulate:', simulate_response)

        # Account Balance
        balance_response = mpesa.account_balance(
            shortcode=os.getenv('MPESA_SHORTCODE', '174379'),
            identifier_type=4,
            remarks='Daily balance check',
        )
        print('Account Balance:', balance_response)

        # Transaction Status
        tx_response = mpesa.transaction_status(
            shortcode=os.getenv('MPESA_SHORTCODE', '174379'),
            transaction_id=stk_response['CheckoutRequestID'],
            identifier_type=1,
            remarks='Transaction check',
        )
        print('Transaction Status:', tx_response)

        # Reversal
        reversal_response = mpesa.reversal(
            shortcode=os.getenv('MPESA_SHORTCODE', '174379'),
            transaction_id='OER7Q9I2PC',
            amount=10,
            remarks='Customer refund',
        )
        print('Reversal:', reversal_response)

        # B2 Pochi
        pochi_response = mpesa.b2pochi(
            phonenumber='254712345678',
            amount=200,
            remarks='Pochi payment',
        )
        print('B2 Pochi:', pochi_response)

    except Exception as e:
        print('Mpesa API error:', e)


if __name__ == '__main__':
    main()
