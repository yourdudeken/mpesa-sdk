export interface RateLimiterConfig {
  tokensPerSecond: number;
  burstSize: number;
}

export interface RateLimiter {
  acquire(): Promise<void>;
  tryAcquire(): boolean;
  available(): number;
}

export class TokenBucketRateLimiter implements RateLimiter {
  private tokens: number;
  private lastRefill: number;
  private readonly maxTokens: number;
  private readonly refillRate: number;
  private readonly refillInterval: number;

  constructor(config: RateLimiterConfig) {
    this.maxTokens = config.burstSize;
    this.tokens = config.burstSize;
    this.lastRefill = Date.now();
    this.refillRate = config.tokensPerSecond;
    this.refillInterval = 1000;
  }

  private refill(): void {
    const now = Date.now();
    const elapsed = now - this.lastRefill;
    const tokensToAdd = (elapsed / this.refillInterval) * this.refillRate;
    this.tokens = Math.min(this.maxTokens, this.tokens + tokensToAdd);
    this.lastRefill = now;
  }

  tryAcquire(): boolean {
    this.refill();
    if (this.tokens >= 1) {
      this.tokens -= 1;
      return true;
    }
    return false;
  }

  async acquire(): Promise<void> {
    while (!this.tryAcquire()) {
      await new Promise((resolve) => setTimeout(resolve, 50));
    }
  }

  available(): number {
    this.refill();
    return this.tokens;
  }
}

export class NoopRateLimiter implements RateLimiter {
  acquire(): Promise<void> {
    return Promise.resolve();
  }

  tryAcquire(): boolean {
    return true;
  }

  available(): number {
    return Infinity;
  }
}
