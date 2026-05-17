export interface MetricsCollector {
  incrementCounter(name: string, labels?: Record<string, string>): void;
  observeHistogram(name: string, value: number, labels?: Record<string, string>): void;
  setGauge(name: string, value: number, labels?: Record<string, string>): void;
}

export class NoopMetricsCollector implements MetricsCollector {
  incrementCounter() {}
  observeHistogram() {}
  setGauge() {}
}

export interface MpesaMetrics {
  requestsTotal: MetricsCollector;
  requestDuration: MetricsCollector;
  errorsTotal: MetricsCollector;
  tokenRefreshes: MetricsCollector;
  circuitBreakerState: MetricsCollector;
}

export function createMpesaMetrics(collector?: MetricsCollector): MpesaMetrics {
  const c = collector ?? new NoopMetricsCollector();
  return {
    requestsTotal: c,
    requestDuration: c,
    errorsTotal: c,
    tokenRefreshes: c,
    circuitBreakerState: c,
  };
}
