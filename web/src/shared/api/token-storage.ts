const ACCESS_KEY = "educrm_access_token";
const REFRESH_KEY = "educrm_refresh_token";

export function getStoredAccessToken(): string | null {
  return (
    localStorage.getItem(ACCESS_KEY) ?? sessionStorage.getItem(ACCESS_KEY)
  );
}

export function getStoredRefreshToken(): string | null {
  return (
    localStorage.getItem(REFRESH_KEY) ?? sessionStorage.getItem(REFRESH_KEY)
  );
}

/** True when tokens are persisted across browser sessions (localStorage). */
export function usesPersistentTokenStorage(): boolean {
  return (
    localStorage.getItem(ACCESS_KEY) !== null ||
    localStorage.getItem(REFRESH_KEY) !== null
  );
}

/**
 * Writes tokens to localStorage (remember me) or sessionStorage (session-only).
 * Always clears both storages first so tokens never sit in two places.
 */
export function setStoredTokens(options_raw: {
  accessToken: string;
  refreshToken: string | null;
  rememberMe: boolean;
}): void {
  const { accessToken, refreshToken, rememberMe } = options_raw;
  clearAllTokens();
  const store = rememberMe ? localStorage : sessionStorage;
  store.setItem(ACCESS_KEY, accessToken);
  if (refreshToken) {
    store.setItem(REFRESH_KEY, refreshToken);
  }
}

/** Updates tokens after refresh, keeping the same storage backend as before. */
export function persistTokensAfterRefresh(tokens: {
  accessToken: string;
  refreshToken: string | null;
}): void {
  const rememberMe = usesPersistentTokenStorage();
  setStoredTokens({
    accessToken: tokens.accessToken,
    refreshToken: tokens.refreshToken,
    rememberMe,
  });
}

export function clearAllTokens(): void {
  localStorage.removeItem(ACCESS_KEY);
  localStorage.removeItem(REFRESH_KEY);
  sessionStorage.removeItem(ACCESS_KEY);
  sessionStorage.removeItem(REFRESH_KEY);
}
