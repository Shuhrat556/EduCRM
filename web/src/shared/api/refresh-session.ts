import axios from "axios";

import { getEnv } from "@/shared/config/env";
import { isRecord } from "@/shared/api/envelope";
import {
  getStoredRefreshToken,
  persistTokensAfterRefresh,
} from "@/shared/api/token-storage";

type RefreshSuccessData = {
  access_token: string;
  refresh_token?: string;
  token_type?: string;
  expires_in?: number;
};

let refreshPromise: Promise<boolean> | null = null;

function parseRefreshPayload(raw: unknown): RefreshSuccessData | null {
  if (!isRecord(raw)) return null;
  if (raw.success === true && "data" in raw && isRecord(raw.data)) {
    const d = raw.data as Record<string, unknown>;
    const access = d.access_token;
    if (typeof access !== "string") return null;
    return {
      access_token: access,
      refresh_token:
        typeof d.refresh_token === "string" ? d.refresh_token : undefined,
      token_type:
        typeof d.token_type === "string" ? d.token_type : undefined,
      expires_in: typeof d.expires_in === "number" ? d.expires_in : undefined,
    };
  }
  /** Bare data (no envelope) — tolerate during rollout */
  if (typeof raw.access_token === "string") {
    return {
      access_token: raw.access_token,
      refresh_token:
        typeof raw.refresh_token === "string"
          ? raw.refresh_token
          : undefined,
    };
  }
  return null;
}

/**
 * Refreshes access token using the stored refresh token (bare axios — avoids
 * interceptor loops). Returns true when tokens were rotated successfully.
 */
export function refreshSession(): Promise<boolean> {
  if (!refreshPromise) {
    refreshPromise = performRefresh().finally(() => {
      refreshPromise = null;
    });
  }
  return refreshPromise;
}

async function performRefresh(): Promise<boolean> {
  const refreshToken = getStoredRefreshToken();
  if (!refreshToken) return false;

  try {
    const base = getEnv().VITE_API_BASE_URL;
    const { data: raw } = await axios.post<unknown>(
      `${base}/api/v1/auth/refresh`,
      { refresh_token: refreshToken },
      { headers: { "Content-Type": "application/json" } },
    );

    const parsed = parseRefreshPayload(raw);
    if (!parsed) return false;

    persistTokensAfterRefresh({
      accessToken: parsed.access_token,
      refreshToken: parsed.refresh_token ?? refreshToken,
    });
    return true;
  } catch {
    return false;
  }
}
