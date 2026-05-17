import axios, { type AxiosInstance, type AxiosRequestConfig } from "axios";
import { getBaseUrl } from "../environment.js";
import type { MpesaConfig, ResolvedConfig, Logger, LoggingHook, AccessTokenResponse, TokenCache } from "../types/index.js";
import { maskSensitiveData, noopLogger } from "../utils/index.js";
import { setupRetryInterceptor, mapAxiosError } from "../interceptors/retry.js";
import { AuthenticationError } from "../errors/index.js";

const DEFAULT_TIMEOUT = 30000;

export class MpesaApiClient {
  private readonly client: AxiosInstance;
  private readonly config: ResolvedConfig;
  private tokenCache: TokenCache | null = null;
  private readonly logging?: LoggingHook;
  private readonly logger: Logger;

  constructor(config: MpesaConfig) {
    this.logger = config.logger ?? noopLogger;

    this.config = {
      consumerKey: config.consumerKey,
      consumerSecret: config.consumerSecret,
      environment: config.environment ?? "sandbox",
      initiatorPassword: config.initiatorPassword ?? "",
      initiatorName: config.initiatorName ?? "",
      passkey: config.passkey ?? "",
      securityCredential: config.securityCredential ?? "",
      retryConfig: config.retryConfig ?? {
        maxRetries: 3,
        baseDelayMs: 1000,
        maxDelayMs: 30000,
      },
      timeout: config.timeout ?? DEFAULT_TIMEOUT,
      logging: config.logging ?? {},
      logger: this.logger,
    };
    this.logging = this.config.logging;

    this.client = axios.create({
      baseURL: getBaseUrl(this.config.environment),
      timeout: this.config.timeout,
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    });

    this.logger.info("M-Pesa API client initialized", {
      environment: this.config.environment,
      timeout: this.config.timeout,
      retryMax: this.config.retryConfig.maxRetries,
    });

    this.client.interceptors.request.use((req) => {
      this.logger.debug("Outgoing request", {
        method: req.method?.toUpperCase(),
        url: req.url,
      });
      this.logging?.onRequest?.({
        method: req.method?.toUpperCase() ?? "GET",
        url: req.url ?? "",
        headers: req.headers as Record<string, string>,
        body: req.data ? maskSensitiveData(req.data) : undefined,
        timestamp: new Date(),
      });
      return req;
    });

    this.client.interceptors.response.use(
      (res) => {
        this.logger.debug("Response received", {
          status: res.status,
          url: res.config.url,
        });
        this.logging?.onResponse?.({
          status: res.status,
          body: res.data,
          durationMs: 0,
          timestamp: new Date(),
        });
        return res;
      },
      (err) => {
        this.logger.error("Request error", {
          message: err.message,
          status: err.response?.status,
          url: err.config?.url,
        });
        this.logging?.onError?.({
          error: err,
          timestamp: new Date(),
        });
        return Promise.reject(mapAxiosError(err));
      },
    );

    setupRetryInterceptor(this.client, this.config.retryConfig, this.logger);
  }

  getConfig(): Readonly<ResolvedConfig> {
    return this.config;
  }

  private isTokenExpired(): boolean {
    if (!this.tokenCache) return true;
    return new Date() >= this.tokenCache.expiresAt;
  }

  async getAccessToken(): Promise<string> {
    if (this.tokenCache && !this.isTokenExpired()) {
      return this.tokenCache.token;
    }

    this.logger.debug("Fetching new access token");

    const response = await this.client.get<AccessTokenResponse>("/oauth/v1/generate", {
      params: { grant_type: "client_credentials" },
      auth: {
        username: this.config.consumerKey,
        password: this.config.consumerSecret,
      },
    });

    const data = response.data;
    this.tokenCache = {
      token: data.access_token,
      expiresAt: new Date(Date.now() + (data.expires_in - 60) * 1000),
    };

    this.logger.debug("Access token acquired", {
      expiresIn: data.expires_in,
    });

    return data.access_token;
  }

  async request<T>(config: AxiosRequestConfig): Promise<T> {
    const token = await this.getAccessToken();
    const mergedConfig: AxiosRequestConfig = {
      ...config,
      headers: {
        ...config.headers,
        Authorization: `Bearer ${token}`,
      },
    };

    const response = await this.client.request<T>(mergedConfig);
    return response.data;
  }

  async post<T>(url: string, data?: unknown): Promise<T> {
    return this.request<T>({ method: "POST", url, data });
  }

  async get<T>(url: string, params?: Record<string, unknown>): Promise<T> {
    return this.request<T>({ method: "GET", url, params });
  }

  invalidateToken(): void {
    this.tokenCache = null;
  }

  async withTokenRefresh<T>(fn: () => Promise<T>): Promise<T> {
    try {
      return await fn();
    } catch (error) {
      if (error instanceof AuthenticationError) {
        this.invalidateToken();
        return await fn();
      }
      throw error;
    }
  }
}
