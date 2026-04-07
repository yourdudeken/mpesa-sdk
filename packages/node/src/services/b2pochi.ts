import { B2cParams } from '../types';
import { HttpClient } from '../http/httpClient';
import { Auth } from '../core/Auth';
import { Helpers } from '../utils/helpers';

export class B2PochiService {
  private httpClient: HttpClient;
  private auth: Auth;
  private helpers: Helpers;
  private baseUrl: string;

  constructor(httpClient: HttpClient, auth: Auth) {
    this.httpClient = httpClient;
    this.auth = auth;
    this.helpers = new Helpers();
    this.baseUrl = httpClient.getBaseUrl();
  }

  async send(params: B2cParams): Promise<any> {
    const { phonenumber, amount, remarks, resultUrl, timeoutUrl, shortCodeType = 'B2C' } = params;

    const url = `${this.baseUrl}/mpesa/b2pochi/v1/paymentrequest`;
    const body = {
      OriginatorConversationID: this.helpers.getFormattedTimestamp(),
      InitiatorName: this.helpers.getConfig('initiatorName'),
      SecurityCredential: this.helpers.generateSecurityCredential(),
      CommandID: 'BusinessPayToPochi',
      Amount: amount,
      PartyA: this.helpers.getConfig('b2cShortcode'),
      PartyB: this.helpers.phoneValidator(phonenumber),
      Remarks: remarks,
      Occassion: '',
      ResultURL: resultUrl || this.helpers.getConfig('callbacks.b2pochiResultUrl'),
      QueueTimeOutURL: timeoutUrl || this.helpers.getConfig('callbacks.b2pochiTimeoutUrl'),
    };

    const token = await this.auth.getAccessToken(shortCodeType);
    return this.httpClient.post(url, body, token);
  }
}