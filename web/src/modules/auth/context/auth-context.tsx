import {
  createContext,
  type ReactNode,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";

import { authApi } from "@/modules/auth/api/auth-api";
import { canAccessPath } from "@/modules/auth/lib/app-access";
import type { AuthUser, LoginCredentials, Role } from "@/modules/auth/model/types";
import {
  clearAllTokens,
  getStoredAccessToken,
  setAuthExpiredHandler,
  setStoredTokens,
} from "@/shared/api/client";
import { ApiError } from "@/shared/api/types";

export type AuthStatus =
  | "idle"
  | "loading"
  | "authenticated"
  | "unauthenticated";

export interface AuthContextValue {
  status: AuthStatus;
  user: AuthUser | null;
  login: (credentials: LoginCredentials) => Promise<AuthUser>;
  logout: () => Promise<void>;
  /** Clear tokens locally without calling the logout API (wrong-portal recovery). */
  clearLocalSession: () => void;
  refreshUser: () => Promise<void>;
  changePassword: (args: {
    currentPassword?: string;
    newPassword: string;
  }) => Promise<void>;
  hasRole: (role: Role | Role[]) => boolean;
  /** True if the current user may navigate to this path (portal + RBAC). */
  canAccessPath: (pathname: string) => boolean;
}

export const AuthContext = createContext<AuthContextValue | null>(null);

function parseRoles(roles: Role | Role[] | undefined, user: AuthUser | null) {
  const list = roles ?? [];
  const needed = Array.isArray(list) ? list : [list];
  if (!user?.roles?.length) return false;
  const set = new Set(user.roles);
  return needed.some((r) => set.has(r));
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [status, setStatus] = useState<AuthStatus>(() =>
    getStoredAccessToken() ? "loading" : "unauthenticated",
  );
  const [user, setUser] = useState<AuthUser | null>(null);

  const clearSession = useCallback(() => {
    clearAllTokens();
    setUser(null);
    setStatus("unauthenticated");
  }, []);

  const clearLocalSession = useCallback(() => {
    clearAllTokens();
    setUser(null);
    setStatus("unauthenticated");
  }, []);

  const bootstrap = useCallback(async () => {
    setStatus("loading");
    try {
      const { data } = await authApi.me();
      setUser(data);
      setStatus("authenticated");
    } catch (e) {
      if (e instanceof ApiError && e.status === 401) {
        clearSession();
        return;
      }
      clearSession();
    }
  }, [clearSession]);

  useEffect(() => {
    setAuthExpiredHandler(() => {
      clearSession();
    });
    return () => setAuthExpiredHandler(null);
  }, [clearSession]);

  useEffect(() => {
    if (!getStoredAccessToken()) return;
    void bootstrap();
  }, [bootstrap]);

  const login = useCallback(
    async (credentials: LoginCredentials): Promise<AuthUser> => {
      const hadExistingSession = getStoredAccessToken() !== null;
      setStatus("loading");
      try {
        const { data } = await authApi.login(credentials);
        setStoredTokens({
          accessToken: data.access_token,
          refreshToken: data.refresh_token ?? null,
          rememberMe: credentials.rememberMe,
        });
        const me = await authApi.me();
        const sessionUser = me.data;
        setUser(sessionUser);
        setStatus("authenticated");
        return sessionUser;
      } catch (e) {
        if (hadExistingSession) {
          void bootstrap();
        } else {
          clearSession();
        }
        throw e;
      }
    },
    [bootstrap, clearSession],
  );

  const logout = useCallback(async () => {
    try {
      await authApi.logout();
    } catch {
      /* still clear locally */
    } finally {
      clearSession();
    }
  }, [clearSession]);

  const refreshUser = useCallback(async () => {
    const { data } = await authApi.me();
    setUser(data);
  }, []);

  const changePassword = useCallback(
    async (args: { currentPassword?: string; newPassword: string }) => {
      await authApi.changePassword({
        current_password: args.currentPassword,
        new_password: args.newPassword,
      });
      await refreshUser();
    },
    [refreshUser],
  );

  const hasRole = useCallback(
    (role: Role | Role[]) => parseRoles(role, user),
    [user],
  );

  const checkPath = useCallback(
    (pathname: string) => canAccessPath(pathname, user?.roles),
    [user?.roles],
  );

  const value = useMemo<AuthContextValue>(
    () => ({
      status,
      user,
      login,
      logout,
      clearLocalSession,
      refreshUser,
      changePassword,
      hasRole,
      canAccessPath: checkPath,
    }),
    [
      status,
      user,
      login,
      logout,
      clearLocalSession,
      refreshUser,
      changePassword,
      hasRole,
      checkPath,
    ],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
