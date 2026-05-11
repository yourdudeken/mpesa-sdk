const { Mpesa } = require('@yourdudeken/mpesa-sdk');

const mpesa = new Mpesa({
  environment: process.env.MPESA_ENV || 'sandbox',
  mpesaConsumerKey: process.env.MPESA_CONSUMER_KEY,
  mpesaConsumerSecret: process.env.MPESA_CONSUMER_SECRET,
  passkey: process.env.MPESA_PASSKEY,
  shortcode: process.env.MPESA_SHORTCODE || '174379',
  initiatorName: process.env.MPESA_INITIATOR_NAME || 'testapi',
  initiatorPassword: process.env.MPESA_INITIATOR_PASSWORD,
  b2cShortcode: process.env.MPESA_B2C_SHORTCODE,
  tillNumber: process.env.MPESA_TILL_NUMBER,
  callbacks: {
    callbackUrl: process.env.MPESA_CALLBACK_URL,
    b2cResultUrl: process.env.MPESA_B2C_RESULT_URL,
    b2cTimeoutUrl: process.env.MPESA_B2C_TIMEOUT_URL,
    b2bResultUrl: process.env.MPESA_B2B_RESULT_URL,
    b2bTimeoutUrl: process.env.MPESA_B2B_TIMEOUT_URL,
    b2pochiResultUrl: process.env.MPESA_B2POCHI_RESULT_URL,
    b2pochiTimeoutUrl: process.env.MPESA_B2POCHI_TIMEOUT_URL,
    c2bValidationUrl: process.env.MPESA_C2B_VALIDATION_URL,
    c2bConfirmationUrl: process.env.MPESA_C2B_CONFIRMATION_URL,
    balanceResultUrl: process.env.MPESA_BALANCE_RESULT_URL,
    balanceTimeoutUrl: process.env.MPESA_BALANCE_TIMEOUT_URL,
    statusResultUrl: process.env.MPESA_STATUS_RESULT_URL,
    statusTimeoutUrl: process.env.MPESA_STATUS_TIMEOUT_URL,
    reversalResultUrl: process.env.MPESA_REVERSAL_RESULT_URL,
    reversalTimeoutUrl: process.env.MPESA_REVERSAL_TIMEOUT_URL,
  },
});

async function main() {
  try {
    // STK Push – Lipa na Mpesa Online
    const stkResponse = await mpesa.stkpush({
      phonenumber: '254712345678',
      amount: 10,
      accountNumber: 'INV-001',
    });
    console.log('STK Push:', stkResponse);

    // Query STK Push status
    const statusResponse = await mpesa.stkquery(stkResponse.CheckoutRequestID);
    console.log('STK Query:', statusResponse);

    // B2C – Business to Customer
    const b2cResponse = await mpesa.b2c({
      phonenumber: '254712345678',
      commandId: 'BusinessPayment',
      amount: 500,
      remarks: 'Salary payment',
    });
    console.log('B2C:', b2cResponse);

    // B2B – Business to Business
    const b2bResponse = await mpesa.b2b({
      receiverShortcode: '600000',
      commandId: 'BusinessPayBill',
      amount: 1000,
      remarks: 'Invoice payment',
      accountNumber: 'INV-001',
    });
    console.log('B2B:', b2bResponse);

    // C2B – Register URLs
    const registerResponse = await mpesa.c2bregisterURLS({
      shortcode: process.env.MPESA_SHORTCODE || '174379',
    });
    console.log('C2B Register:', registerResponse);

    // C2B – Simulate payment
    const simulateResponse = await mpesa.c2bsimulate({
      phonenumber: '254712345678',
      amount: 100,
      shortcode: process.env.MPESA_SHORTCODE || '174379',
      commandId: Mpesa.PAYBILL,
    });
    console.log('C2B Simulate:', simulateResponse);

    // Account Balance
    const balanceResponse = await mpesa.accountBalance({
      shortcode: process.env.MPESA_SHORTCODE || '174379',
      identifierType: 4,
      remarks: 'Daily balance check',
    });
    console.log('Account Balance:', balanceResponse);

    // Transaction Status
    const txStatusResponse = await mpesa.transactionStatus({
      shortcode: process.env.MPESA_SHORTCODE || '174379',
      transactionId: stkResponse.CheckoutRequestID,
      identifierType: 1,
      remarks: 'Transaction check',
    });
    console.log('Transaction Status:', txStatusResponse);

    // Reversal
    const reversalResponse = await mpesa.reversal({
      shortcode: process.env.MPESA_SHORTCODE || '174379',
      transactionId: 'OER7Q9I2PC',
      amount: 10,
      remarks: 'Customer refund',
    });
    console.log('Reversal:', reversalResponse);

    // B2 Pochi
    const pochiResponse = await mpesa.b2pochi({
      phonenumber: '254712345678',
      amount: 200,
      remarks: 'Pochi payment',
    });
    console.log('B2 Pochi:', pochiResponse);
  } catch (error) {
    console.error('Mpesa API error:', error.message);
  }
}

main();
