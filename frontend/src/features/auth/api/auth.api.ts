import { apiClient } from "@/services/api-client";
import type { LoginRequest, AuthResponse, MeResponse, RegisterRequest } from "../type";

export const authApi = {
    login: (data: LoginRequest): Promise<AuthResponse> =>
        apiClient.request("/auth/login", {
            method: "POST",
            body: JSON.stringify(data),
        }),

    register: (data: RegisterRequest) =>
        apiClient.request("/auth/register", {
            method: "POST",
            body: JSON.stringify(data),
        }),

    me: (): Promise<MeResponse> => apiClient.request("/me"),
};
