import type { Logger } from "../types/index.js";

export interface Tracer {
  startSpan(name: string, attributes?: Record<string, string>): { end: () => void };
}

export class NoopTracer implements Tracer {
  startSpan(): { end: () => void } {
    return { end: () => {} };
  }
}

export function createTracer(_logger?: Logger): Tracer {
  return new NoopTracer();
}

export function withSpan<T>(
  tracer: Tracer,
  name: string,
  fn: () => Promise<T>,
  attributes?: Record<string, string>,
): Promise<T> {
  const span = tracer.startSpan(name, attributes);
  return fn().finally(() => span.end());
}
