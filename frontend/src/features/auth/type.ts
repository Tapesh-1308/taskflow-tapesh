export interface LoginRequest {
    email: string;
    password: string;
}

export interface AuthResponse {
    token: string;
    user: {
        id: string;
        name: string;
        email: string;
    };
}

export interface MeResponse {
    user_id: string;
    name: string;
}

export interface  RegisterRequest {
    name: string;
    email: string;
    password: string;
}