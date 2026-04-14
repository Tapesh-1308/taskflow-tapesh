import type { Task } from "@/features/task/types"

export interface Project {
    id: string
    name: string
    description?: string | null
    owner_id: string
    created_at: string
    tasks: Task[]
}