// src/services/api-client.ts

const API_URL = "http://localhost:8080";

export class ApiClient {
    private token: string | null = null;

    constructor() {
        // restore token on refresh
        this.token = localStorage.getItem("token");
    }

    setToken(token: string) {
        this.token = token;
        localStorage.setItem("token", token);
    }

    clearToken() {
        this.token = null;
        localStorage.removeItem("token");
    }

    isAuthenticated() {
        return !!this.token;
    }

    async request(endpoint: string, options: RequestInit = {}) {
        const res = await fetch(`${API_URL}${endpoint}`, {
            ...options,
            headers: {
                "Content-Type": "application/json",
                ...(this.token && {
                    Authorization: `Bearer ${this.token}`,
                }),
                ...options.headers,
            },
        });

        if (!res.ok) {
            // auto logout on unauthorized
            if (res.status === 401) {
                this.clearToken();
                throw new Error("You don't have access to this project")
            } else if (res.status === 404) {
                throw new Error("Oops! project not found :(")
            }           

            const error = await res.json().catch(() => ({}));
            throw new Error(error.message || res.statusText);
        }

        return res.json();
    }
}

export const apiClient = new ApiClient();