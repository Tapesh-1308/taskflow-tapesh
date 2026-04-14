import { apiClient } from "@/services/api-client";
import type { User } from "../type";

export const userApi = {
    search: (query: string): Promise<User[]> =>
        apiClient.request(`/users?search=${query}`),
};