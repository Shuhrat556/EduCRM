import axios, {
  type AxiosError,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
  isAxiosError,
} from "axios";

import { getEnv } from "@/shared/config/env";
import {
  httpStatusFromErrorKind,
  isRecord,
} from "@/shared/api/envelope";
import { ApiError, type ApiErrorBody } from "@/shared/api/types";
import { refreshSession } from "@/shared/api/refresh-session";
import { getStoredAccessToken } from "@/shared/api/token-storage";

export { getStoredAccessToken } from "@/shared/api/token-storage";
export {
  setStoredTokens,
  clearAllTokens,
  getStoredRefreshToken,
  persistTokensAfterRefresh,
  usesPersistentTokenStorage,
} from "@/shared/api/token-storage";

export const apiClient = axios.create({
  baseURL: `${getEnv().VITE_API_BASE_URL}/api/v1`,
  headers: { "Content-Type": "application/json" },
  withCredentials: false,
});

type RetryConfig = InternalAxiosRequestConfig & { _educrmRetry?: boolean };

function attachAuthHeader(config: InternalAxiosRequestConfig) {
  const token = getStoredAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}

apiClient.interceptors.request.use(attachAuthHeader);

export type AuthExpiredHandler = () => void;

let onAuthExpired: AuthExpiredHandler | null = null;

export function setAuthExpiredHandler(handler: AuthExpiredHandler | null) {
  onAuthExpired = handler;
}

function envelopeFailureToApiError(
  payload: Record<string, unknown>,
  fallbackStatus: number,
): ApiError {
  const errRaw = payload.error;
  const err = isRecord(errRaw) ? errRaw : {};
  const message = String(err.message ?? "Request failed");
  const kind = typeof err.kind === "string" ? err.kind : undefined;
  const code = typeof err.code === "string" ? err.code : undefined;
  const status =
    fallbackStatus > 0 ? fallbackStatus : httpStatusFromErrorKind(kind);
  return new ApiError(message, status || 500, {
    message,
    code,
    kind,
  });
}

function unwrapJsonEnvelope(response: AxiosResponse): AxiosResponse {
  const ct = response.headers?.["content-type"];
  if (typeof ct === "string" && !ct.includes("application/json")) {
    return response;
  }

  const payload = response.data;
  if (!isRecord(payload) || !("success" in payload)) {
    return response;
  }

  if (payload.success === true && "data" in payload) {
    response.data = payload.data;
    return response;
  }

  if (payload.success === false) {
    throw envelopeFailureToApiError(payload, response.status);
  }

  return response;
}

function normalizeError(error: AxiosError<unknown>): ApiError {
  const status = error.response?.status ?? 0;
  const raw = error.response?.data;

  if (isRecord(raw) && raw.success === false) {
    return envelopeFailureToApiError(raw, status);
  }

  const body = raw as ApiErrorBody | undefined;
  const message =
    body?.message ||
    error.message ||
    (status === 0 ? "Network error" : `Request failed (${status})`);
  return new ApiError(message, status, body);
}

function skipRefreshForUrl(url: string | undefined) {
  if (!url) return false;
  return url.includes("/auth/login") || url.includes("/auth/refresh");
}

apiClient.interceptors.response.use(unwrapJsonEnvelope, async (error: unknown) => {
  if (!isAxiosError<unknown>(error) || !error.config) {
    return Promise.reject(error);
  }

  const config = error.config as RetryConfig;
  const url = config.url ?? "";

  if (error.response?.status !== 401) {
    return Promise.reject(normalizeError(error));
  }

  if (url.includes("/auth/logout")) {
    return Promise.reject(normalizeError(error));
  }

  if (skipRefreshForUrl(url) || config._educrmRetry) {
    onAuthExpired?.();
    return Promise.reject(normalizeError(error));
  }

  const refreshed = await refreshSession();
  if (!refreshed) {
    onAuthExpired?.();
    return Promise.reject(normalizeError(error));
  }

  const token = getStoredAccessToken();
  if (!token) {
    onAuthExpired?.();
    return Promise.reject(normalizeError(error));
  }

  config._educrmRetry = true;
  config.headers.Authorization = `Bearer ${token}`;
  return apiClient(config);
});
