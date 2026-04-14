import { createBrowserRouter } from "react-router-dom";
import LoginPage from "@/features/auth/pages/login";
import { ProtectedRoute } from "@/shared/components/protected-route";
import ProjectsPage from "@/features/project/pages/projects";
import ProjectDetailPage from "@/features/project/pages/project-detail";
import RegisterPage from "@/features/auth/pages/register";
import { AppLayout } from "@/layouts/app-layout";

export const router = createBrowserRouter([
    {
        path: "/login",
        element: <LoginPage />,
    },
    {
        path: "/register",
        element: <RegisterPage />,
    },
    {
        element: (
            <ProtectedRoute>
                <AppLayout />
            </ProtectedRoute>
        ),
        children: [
            {
                path: "/",
                element: (
                    <ProtectedRoute>
                        <ProjectsPage />
                    </ProtectedRoute>
                ),
            },
            {
                path: "/projects/:id",
                element: (
                    <ProtectedRoute>
                        <ProjectDetailPage />
                    </ProtectedRoute>
                ),
            },
        ]
    },
]);
