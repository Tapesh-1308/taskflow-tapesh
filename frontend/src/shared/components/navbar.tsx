import { useNavigate } from "react-router-dom"
import { useQuery } from "@tanstack/react-query"

import { authApi } from "@/features/auth/api/auth.api"
import { apiClient } from "@/services/api-client"
import { Button } from "@/components/ui/button"

export function Navbar() {
    const navigate = useNavigate();

    const { data, isLoading, isError } = useQuery({
        queryKey: ["me"],
        queryFn: authApi.me,
        enabled: apiClient.isAuthenticated(),
        retry: false,
    });

    if (isError) {
        apiClient.clearToken();
        navigate("/login");
    }

    const handleLogout = () => {
        apiClient.clearToken();
        navigate("/login");
    };

    return (
        <div className="flex items-center justify-between px-6 py-3 border-b">
            <h1
                className="font-semibold cursor-pointer"
                onClick={() => navigate("/")}
            >
                TaskFlow
            </h1>

            <div className="flex items-center gap-4">
                {isLoading ? (
                    <span className="text-sm text-gray-400">Loading...</span>
                ) : (
                    <span className="text-sm text-gray-600">
                        {data?.name}
                    </span>
                )}

                <Button variant="outline" onClick={handleLogout}>
                    Logout
                </Button>
            </div>
        </div>
    );
}