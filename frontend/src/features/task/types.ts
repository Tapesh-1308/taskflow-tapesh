export type TaskStatus = "todo" | "in_progress" | "done";

export interface Task {
    id: string;
    title: string;
    description?: string;
    status: TaskStatus;
    priority?: string;
    assignee_id?: string;
    project_id: string;
    created_at: string;
    updated_at: string;
    due_date?: string;
}

export interface TaskWithUser {
    id: string;
    title: string;
    description?: string;
    status: TaskStatus;
    priority?: string;
    assignee?: {
        id: string;
        name: string;
    };
    project_id: string;
    created_at: string;
    updated_at: string;
    due_date?: string;
}