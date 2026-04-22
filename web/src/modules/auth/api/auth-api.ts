import type { AxiosResponse } from "axios";

import { apiClient } from "@/shared/api/client";
import { buildLoginRequestBody } from "@/modules/auth/api/login-payload";
import { mapApiUserToAuthUser } from "@/modules/auth/api/map-api-user";
import type { AuthUser, LoginCredentials } from "@/modules/auth/model/types";

/** Inner `data` from login envelope (snake_case tokens). */
export type LoginSuccessData = {
  access_token: string;
  refresh_token?: string;
  token_type?: string;
  expires_in?: number;
};

export const authApi = {
  login: (credentials: LoginCredentials) =>
    apiClient.post<LoginSuccessData>(
      "/auth/login",
      buildLoginRequestBody(credentials),
    ),

  logout: () => apiClient.post<void>("/auth/logout"),

  me: (): Promise<AxiosResponse<AuthUser>> =>
    apiClient.get<unknown>("/auth/me").then((res: AxiosResponse<unknown>) => ({
      ...res,
      data: mapApiUserToAuthUser(res.data),
    })),

  changePassword: (body: {
    current_password?: string;
    new_password: string;
  }) => apiClient.post<void>("/auth/change-password", body),
};
