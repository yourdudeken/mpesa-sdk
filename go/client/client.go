package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/yourdudeken/mpesa-sdk/go/errors"
	"github.com/yourdudeken/mpesa-sdk/go/types"
	"github.com/yourdudeken/mpesa-sdk/go/validation"
)

func (c *Client) RotateCredentials(consumerKey, consumerSecret string) {
	c.config.ConsumerKey = consumerKey
	c.config.ConsumerSecret = consumerSecret
	c.tokenManager.Invalidate()
	c.logger.Info("Credentials rotated")
}

func generateIdempotencyKey() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}

func generateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "mpesa-unknown"
	}
	return "mpesa-" + hex.EncodeToString(b)
}

var retryableStatusCodes = map[int]bool{
	408: true, 429: true,
	500: true, 502: true, 503: true, 504: true,
}

type TokenManager struct {
	mu             sync.RWMutex
	token          string
	expiresAt      time.Time
	endpoint       string
	consumerKey    string
	consumerSecret string
	httpClient     *http.Client
	logger         types.Logger
}

func NewTokenManager(endpoint, consumerKey, consumerSecret string, httpClient *http.Client, logger types.Logger) *TokenManager {
	return &TokenManager{
		endpoint:       endpoint,
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,
		httpClient:     httpClient,
		logger:         logger,
	}
}

func (tm *TokenManager) SetAuthEndpoint(endpoint string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.endpoint = endpoint
}

func (tm *TokenManager) GetToken(ctx context.Context) (string, error) {
	tm.mu.RLock()
	if tm.token != "" && time.Now().Before(tm.expiresAt) {
		defer tm.mu.RUnlock()
		return tm.token, nil
	}
	tm.mu.RUnlock()

	tm.mu.Lock()
	defer tm.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, "GET", tm.endpoint, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()
	req.SetBasicAuth(tm.consumerKey, tm.consumerSecret)

	tm.logger.Debug("Fetching new access token")

	resp, err := tm.httpClient.Do(req)
	if err != nil {
		return "", errors.NewAPIConnectionError("Failed to get access token",
			errors.WithCause(err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp types.AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	tm.token = tokenResp.AccessToken
	tm.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)

	tm.logger.Debug("Access token acquired",
		"expires_in", tokenResp.ExpiresIn,
	)

	return tm.token, nil
}

func (c *Client) Logger() types.Logger {
	return c.logger
}

func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	return c.tokenManager.GetToken(ctx)
}

func (tm *TokenManager) Invalidate() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.token = ""
	tm.expiresAt = time.Time{}
}

type Client struct {
	config         types.MpesaConfig
	httpClient     *http.Client
	tokenManager   *TokenManager
	endpoints      environmentEndpoints
	logger         types.Logger
	circuitBreaker *types.CircuitBreaker
	rateLimiter    types.RateLimiter
}

func NewClient(config types.MpesaConfig) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryConfig.MaxRetries == 0 {
		config.RetryConfig = types.RetryConfig{
			MaxRetries:  3,
			BaseDelayMs: 1000,
			MaxDelayMs:  30000,
		}
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	eps := getEndpoints(config.Environment)

	logger := config.Logger
	if logger == nil {
		logger = types.NewNoopLogger()
	}

	logger.Debug("Creating M-Pesa client",
		"environment", config.Environment,
		"timeout", config.Timeout.String(),
		"retry_max", config.RetryConfig.MaxRetries,
	)

	var rl types.RateLimiter
	if config.RateLimiterConfig.TokensPerSecond > 0 {
		rl = types.NewTokenBucketRateLimiter(config.RateLimiterConfig)
	} else {
		rl = &types.NoopRateLimiter{}
	}

	return &Client{
		config:         config,
		httpClient:     httpClient,
		circuitBreaker: types.NewCircuitBreaker(config.CircuitBreakerConfig),
		rateLimiter:    rl,
		tokenManager: NewTokenManager(
			eps.Auth,
			config.ConsumerKey,
			config.ConsumerSecret,
			httpClient,
			logger,
		),
		endpoints: eps,
		logger:    logger,
	}
}

func (c *Client) doRequest(ctx context.Context, method, url string, body interface{}) (respBody []byte, err error) {
	requestID := generateRequestID()

	c.rateLimiter.Acquire()

	if c.circuitBreaker.State() == types.CircuitOpen {
		return nil, types.ErrCircuitBreakerOpen
	}

	defer func() {
		if err != nil {
			c.circuitBreaker.RecordFailure()
		} else {
			c.circuitBreaker.RecordSuccess()
		}
	}()

	c.logger.Debug("Sending request",
		"method", method,
		"url", url,
		"request_id", requestID,
	)

	var lastErr error
	for attempt := 0; attempt <= c.config.RetryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Warn("Retrying request",
				"method", method,
				"url", url,
				"attempt", attempt,
				"max_retries", c.config.RetryConfig.MaxRetries,
				"request_id", requestID,
			)
		}
		token, err := c.tokenManager.GetToken(ctx)
		if err != nil {
			return nil, err
		}

		var reqBodyReader io.Reader
		var bodyData []byte
		if body != nil {
			bodyData, err = json.Marshal(body)
			if err != nil {
				return nil, errors.NewValidationError("Failed to marshal request body",
					errors.WithCause(err))
			}
			reqBodyReader = bytes.NewReader(bodyData)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, reqBodyReader)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Request-ID", requestID)
		if c.config.IdempotencyEnabled && method == "POST" {
			req.Header.Set("X-Idempotency-Key", generateIdempotencyKey())
		}

		if len(bodyData) > 0 {
			req.ContentLength = int64(len(bodyData))
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = errors.NewAPIConnectionError("Request failed",
				errors.WithCause(err),
				errors.WithRequestID(requestID))
			if attempt < c.config.RetryConfig.MaxRetries {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
				}
				delay := calculateBackoffDuration(attempt, c.config.RetryConfig)
				time.Sleep(delay)
				continue
			}
			return nil, lastErr
		}

		respBody, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if retryableStatusCodes[resp.StatusCode] && attempt < c.config.RetryConfig.MaxRetries {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			var delay time.Duration
			if resp.StatusCode == 429 {
				retryAfter := resp.Header.Get("Retry-After")
				if retryAfter == "0" {
					delay = calculateBackoffDuration(attempt, c.config.RetryConfig)
				} else {
					delay = 5 * time.Second
				}
			} else {
				delay = calculateBackoffDuration(attempt, c.config.RetryConfig)
			}
			c.logger.Warn("Retryable status code, backing off",
				"status", resp.StatusCode,
				"attempt", attempt,
				"delay_ms", delay.Milliseconds(),
				"request_id", requestID,
			)
			time.Sleep(delay)
			continue
		}

		if resp.StatusCode == 401 {
			c.tokenManager.Invalidate()
			c.logger.Error("Authentication failed, token invalidated",
				"status", resp.StatusCode,
				"request_id", requestID,
			)
			return nil, errors.NewAuthenticationError("",
				errors.WithStatusCode(resp.StatusCode),
				errors.WithRequestID(requestID),
				errors.WithRawResponse(string(respBody)))
		}

		if resp.StatusCode == 429 {
			c.logger.Error("Rate limit exceeded",
				"status", resp.StatusCode,
				"request_id", requestID,
			)
			return nil, errors.NewRateLimitError("", 60,
				errors.WithStatusCode(resp.StatusCode),
				errors.WithRequestID(requestID),
				errors.WithRawResponse(string(respBody)))
		}

		if resp.StatusCode >= 400 {
			var errResp struct {
				RequestID    string `json:"requestId"`
				ErrorCode    string `json:"errorCode"`
				ErrorMessage string `json:"errorMessage"`
			}
			json.Unmarshal(respBody, &errResp)
			c.logger.Error("API error response",
				"status", resp.StatusCode,
				"error_code", errResp.ErrorCode,
				"error_message", errResp.ErrorMessage,
				"request_id", requestID,
			)
			return nil, errors.NewMpesaAPIError(errResp.ErrorMessage, errResp.ErrorCode,
				errors.WithStatusCode(resp.StatusCode),
				errors.WithRequestID(errResp.RequestID),
				errors.WithRawResponse(string(respBody)))
		}

		c.logger.Debug("Request successful",
			"method", method,
			"url", url,
			"status", resp.StatusCode,
			"request_id", requestID,
		)
		return respBody, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.NewAPIConnectionError("Request failed after retries",
		errors.WithRequestID(requestID))
}

func calculateBackoffDuration(attempt int, config types.RetryConfig) time.Duration {
	exponential := float64(config.BaseDelayMs) * math.Pow(2, float64(attempt))
	if exponential > float64(config.MaxDelayMs) {
		exponential = float64(config.MaxDelayMs)
	}
	return time.Duration(exponential) * time.Millisecond
}

// ---- STK Push ----
func (c *Client) STKPush(ctx context.Context, req types.STKPushRequest) (*types.STKPushResponse, error) {
	if err := validation.PositiveInt(req.Amount, "Amount"); err != nil {
		return nil, err
	}
	if err := validation.PhoneNumber(req.PhoneNumber, "PhoneNumber"); err != nil {
		return nil, err
	}
	if err := validation.ValidURL(req.CallBackURL, "CallBackURL"); err != nil {
		return nil, err
	}
	if err := validation.MaxLength(req.AccountReference, "AccountReference", 12); err != nil {
		return nil, err
	}
	if err := validation.MaxLength(req.TransactionDesc, "TransactionDesc", 13); err != nil {
		return nil, err
	}
	if err := validation.OneOf(string(req.TransactionType), "TransactionType",
		[]string{string(types.CustomerPayBillOnline), string(types.CustomerBuyGoodsOnline)}); err != nil {
		return nil, err
	}

	if req.Password == "" && c.config.Passkey != "" {
		timestamp := req.Timestamp
		if timestamp == "" {
			timestamp = GenerateTimestamp()
		}
		req.Password = GeneratePassword(req.BusinessShortCode, c.config.Passkey, timestamp)
		req.Timestamp = timestamp
	}

	respBody, err := c.doRequest(ctx, "POST", c.endpoints.STKPush, req)
	if err != nil {
		return nil, err
	}

	var resp types.STKPushResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) STKQuery(ctx context.Context, req types.STKQueryRequest) (*types.STKQueryResponse, error) {
	if req.Password == "" && c.config.Passkey != "" {
		timestamp := req.Timestamp
		if timestamp == "" {
			timestamp = GenerateTimestamp()
		}
		req.Password = GeneratePassword(req.BusinessShortCode, c.config.Passkey, timestamp)
		req.Timestamp = timestamp
	}

	respBody, err := c.doRequest(ctx, "POST", c.endpoints.STKQuery, req)
	if err != nil {
		return nil, err
	}

	var resp types.STKQueryResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ---- C2B ----
func (c *Client) C2BRegisterURL(ctx context.Context, req types.C2BRegisterURLRequest) (*types.C2BResponse, error) {
	respBody, err := c.doRequest(ctx, "POST", c.endpoints.C2BRegisterURL, req)
	if err != nil {
		return nil, err
	}

	var resp types.C2BResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) C2BSimulate(ctx context.Context, req types.C2BSimulateRequest) (*types.C2BResponse, error) {
	respBody, err := c.doRequest(ctx, "POST", c.endpoints.C2BSimulate, req)
	if err != nil {
		return nil, err
	}

	var resp types.C2BResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ---- B2C ----
func (c *Client) B2C(ctx context.Context, req types.B2CRequest) (*types.B2CResponse, error) {
	respBody, err := c.doRequest(ctx, "POST", c.endpoints.B2C, req)
	if err != nil {
		return nil, err
	}

	var resp types.B2CResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ---- B2B ----
func (c *Client) B2B(ctx context.Context, req types.B2BRequest) (*types.B2BResponse, error) {
	respBody, err := c.doRequest(ctx, "POST", c.endpoints.B2B, req)
	if err != nil {
		return nil, err
	}

	var resp types.B2BResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ---- Reversal ----
func (c *Client) Reversal(ctx context.Context, req types.ReversalRequest) (*types.ReversalResponse, error) {
	req.CommandID = "TransactionReversal"
	respBody, err := c.doRequest(ctx, "POST", c.endpoints.Reversal, req)
	if err != nil {
		return nil, err
	}

	var resp types.ReversalResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ---- Transaction Status ----
func (c *Client) TransactionStatus(ctx context.Context, req types.TransactionStatusRequest) (*types.TransactionStatusResponse, error) {
	respBody, err := c.doRequest(ctx, "POST", c.endpoints.TransactionStatus, req)
	if err != nil {
		return nil, err
	}

	var resp types.TransactionStatusResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ---- Account Balance ----
func (c *Client) AccountBalance(ctx context.Context, req types.AccountBalanceRequest) (*types.AccountBalanceResponse, error) {
	respBody, err := c.doRequest(ctx, "POST", c.endpoints.AccountBalance, req)
	if err != nil {
		return nil, err
	}

	var resp types.AccountBalanceResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ---- Dynamic QR ----
func (c *Client) DynamicQR(ctx context.Context, req types.DynamicQRRequest) (*types.DynamicQRResponse, error) {
	respBody, err := c.doRequest(ctx, "POST", c.endpoints.DynamicQR, req)
	if err != nil {
		return nil, err
	}

	var resp types.DynamicQRResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ---- Callback Parsing ----
func ParseSTKCallback(payload types.STKCallbackPayload) types.STKCallbackResult {
	cb := payload.Body.StkCallback
	result := types.STKCallbackResult{
		Success:           cb.ResultCode == 0,
		MerchantRequestID: cb.MerchantRequestID,
		CheckoutRequestID: cb.CheckoutRequestID,
		ResultCode:        cb.ResultCode,
		ResultDescription: cb.ResultDesc,
	}

	if cb.CallbackMetadata != nil {
		for _, item := range cb.CallbackMetadata.Item {
			switch item.Name {
			case "Amount":
				if v, ok := item.Value.(float64); ok {
					result.Amount = &v
				}
			case "MpesaReceiptNumber":
				if v, ok := item.Value.(string); ok {
					result.ReceiptNumber = &v
				}
			case "TransactionDate":
				if v, ok := item.Value.(string); ok {
					result.TransactionDate = &v
				}
			case "PhoneNumber":
				if v, ok := item.Value.(string); ok {
					result.PhoneNumber = &v
				}
			}
		}
	}

	return result
}
