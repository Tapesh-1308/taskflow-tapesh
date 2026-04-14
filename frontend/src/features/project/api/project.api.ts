import { apiClient } from "@/services/api-client";

export const projectApi = {
    list: async () => {
        const res = await apiClient.request("/projects");
        return res.projects;
    },

    create: (data: { name: string; description?: string }) =>
        apiClient.request("/projects", {
            method: "POST",
            body: JSON.stringify(data),
        }),

    getById: (id: string) => apiClient.request(`/projects/${id}`),

    update: (id: string, data: { name: string; description?: string }) =>
        apiClient.request(`/projects/${id}`, {
            method: "PATCH",
            body: JSON.stringify(data),
        }),

    delete: (id: string) =>
        apiClient.request(`/projects/${id}`, {
            method: "DELETE",
        }),

    getTasks: async (
        projectId: string,
        params?: { status?: string; assignee?: string },
    ) => {
        const query = new URLSearchParams();

        if (params?.status) query.append("status", params.status);
        if (params?.assignee) query.append("assignee", params.assignee);

        const res = await apiClient.request(
            `/projects/${projectId}/tasks?${query.toString()}`,
        );
        return res.tasks;
    },
};
