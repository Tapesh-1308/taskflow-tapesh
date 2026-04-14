import { useState, useEffect } from "react";
import { authApi } from "../api/auth.api";
import { apiClient } from "@/services/api-client";

export function useAuth() {
    const [user, setUser] = useState<any>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const token = localStorage.getItem("token");
        if (token) {
            apiClient.setToken(token);
            setUser({}); // ideally fetch user profile
        }
        setLoading(false);
    }, []);

    const login = async (email: string, password: string) => {
        const res = await authApi.login({ email, password });
        apiClient.setToken(res.token);
        setUser(res.user);
    };

    const logout = () => {
        apiClient.clearToken();
        setUser(null);
    };

    return { user, loading, login, logout };
}
