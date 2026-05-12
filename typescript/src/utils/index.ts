import { publicEncrypt } from "node:crypto";

export function generateTimestamp(): string {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");
  const hours = String(now.getHours()).padStart(2, "0");
  const minutes = String(now.getMinutes()).padStart(2, "0");
  const seconds = String(now.getSeconds()).padStart(2, "0");
  return `${year}${month}${day}${hours}${minutes}${seconds}`;
}

export function generatePassword(
  shortcode: number | string,
  passkey: string,
  timestamp: string,
): string {
  const toEncode = `${shortcode}${passkey}${timestamp}`;
  return Buffer.from(toEncode).toString("base64");
}

export function generateSecurityCredential(
  password: string,
  certificate: string,
): string {
  const certBuffer = Buffer.from(certificate);
  const encrypted = publicEncrypt(certBuffer, Buffer.from(password));
  return encrypted.toString("base64");
}

export function maskSensitiveData(data: Record<string, unknown>): Record<string, unknown> {
  const sensitiveKeys = [
    "consumerKey",
    "consumerSecret",
    "Password",
    "SecurityCredential",
    "InitiatorPassword",
    "passkey",
    "securityCredential",
    "initiatorPassword",
  ];
  const masked = { ...data };
  for (const key of sensitiveKeys) {
    if (key in masked) {
      const val = String(masked[key]);
      if (val.length > 4) {
        masked[key] = `${val.slice(0, 4)}****`;
      } else {
        masked[key] = "****";
      }
    }
  }
  return masked;
}

export function isPhoneNumberValid(phone: number | string): boolean {
  const str = String(phone);
  return /^2547\d{8}$/.test(str);
}

export function formatPhoneNumber(
  phone: number | string,
): string {
  let str = String(phone).replace(/^0+/, "");
  if (str.startsWith("7")) {
    str = `254${str}`;
  } else if (str.startsWith("+254")) {
    str = str.slice(1);
  }
  return str;
}

export function delay(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export function calculateBackoff(
  attempt: number,
  baseDelayMs: number,
  maxDelayMs: number,
): number {
  const exponential = baseDelayMs * 2 ** attempt;
  const jitter = Math.random() * 100;
  return Math.min(exponential + jitter, maxDelayMs);
}
