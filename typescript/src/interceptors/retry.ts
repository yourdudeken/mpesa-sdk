import type { AxiosError, AxiosInstance } from "axios";
import { RateLimitError, APIConnectionError, TimeoutError } from "../errors/index.js";
import { calculateBackoff, delay } from "../utils/index.js";
import type { RetryConfig } from "../types/index.js";
import { AuthenticationError, MpesaAPIError } from "../errors/index.js";

const DEFAULT_RETRY_CONFIG: RetryConfig = {
  maxRetries: 3,
  baseDelayMs: 1000,
  maxDelayMs: 30000,
};

const RETRYABLE_STATUS_CODES = new Set([408, 429, 500, 502, 503, 504]);
const RETRYABLE_ERROR_CODES = new Set([
  "ECONNRESET",
  "ECONNREFUSED",
  "ENOTFOUND",
  "ETIMEDOUT",
  "ERR_NETWORK",
]);

export function setupRetryInterceptor(
  client: AxiosInstance,
  retryConfig: RetryConfig = DEFAULT_RETRY_CONFIG,
): void {
  client.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      const config = error.config;
      if (!config) return Promise.reject(error);

      config.retryCount = (config.retryCount ?? 0) + 1;

      if (config.retryCount > retryConfig.maxRetries) {
        return Promise.reject(error);
      }

      const shouldRetry =
        RETRYABLE_STATUS_CODES.has(error.response?.status ?? 0) ||
        RETRYABLE_ERROR_CODES.has(error.code ?? "");

      if (!shouldRetry) {
        return Promise.reject(error);
      }

      if (error.response?.status === 429) {
        const retryAfter = parseInt(
          error.response.headers["retry-after"] ?? "5",
          10,
        );
        await delay(retryAfter * 1000);
      } else {
        const backoff = calculateBackoff(
          config.retryCount - 1,
          retryConfig.baseDelayMs,
          retryConfig.maxDelayMs,
        );
        await delay(backoff);
      }

      return client.request(config);
    },
  );
}

export function mapAxiosError(error: AxiosError): Error {
  const status = error.response?.status;
  const data = error.response?.data as Record<string, unknown> | undefined;
  const requestId = (data?.requestId as string) ?? error.config?.headers?.["X-Request-Id"] as string;

  if (error.code === "ECONNABORTED" || error.code === "ETIMEDOUT") {
    return new TimeoutError("Request to M-Pesa API timed out.", {
      statusCode: status,
      requestId,
      rawResponse: data,
      cause: error,
    });
  }

  if (
    error.code === "ECONNREFUSED" ||
    error.code === "ENOTFOUND" ||
    error.code === "ERR_NETWORK"
  ) {
    return new APIConnectionError("Failed to connect to M-Pesa API.", {
      statusCode: status,
      requestId,
      rawResponse: data,
      cause: error,
    });
  }

  if (status === 429) {
    const retryAfter = parseInt(
      error.response?.headers?.["retry-after"] as string ?? "60",
      10,
    );
    return new RateLimitError("M-Pesa API rate limit exceeded.", {
      statusCode: status,
      requestId,
      rawResponse: data,
      retryAfter,
      cause: error,
    });
  }

  if (status === 401) {
    return new AuthenticationError(
      "M-Pesa API authentication failed. Check your consumer key/secret.",
      { statusCode: status, requestId, rawResponse: data, cause: error },
    );
  }

  return new MpesaAPIError(
    (data?.errorMessage as string) ?? error.message,
    {
      statusCode: status,
      errorCode: data?.errorCode as string,
      requestId,
      rawResponse: data,
      cause: error,
    },
  );
}
