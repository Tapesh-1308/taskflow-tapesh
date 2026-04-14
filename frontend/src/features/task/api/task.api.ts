import { apiClient } from "@/services/api-client";
import type { Task } from "../types";

export const taskApi = {
    list: (projectId: string): Promise<Task[]> =>
        apiClient.request(`/tasks?project_id=${projectId}`),

    update: (id: string, data: Partial<Task>) =>
        apiClient.request(`/tasks/${id}`, {
            method: "PATCH",
            body: JSON.stringify(data),
        }),

    create: (data: { title: string; project_id: string }) => {
        return apiClient.request("/projects/" + data.project_id + "/tasks", {
            method: "POST",
            body: JSON.stringify(data),
        })
    }
};
