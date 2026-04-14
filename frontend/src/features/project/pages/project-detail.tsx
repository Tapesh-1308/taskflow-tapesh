import { useParams } from "react-router-dom";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";

import { projectApi } from "@/features/project/api/project.api";

import ProjectHeader from "../components/project-header";
import ProjectFilters from "../components/project-filters";

import { TaskList } from "@/features/task/components/task-list";
import { CreateTaskModal } from "@/features/task/components/create-task-modal";

export default function ProjectDetailPage() {
    const { id } = useParams();

    const [status, setStatus] = useState("");
    const [assignee, setAssignee] = useState("");

    // PROJECT
    const { data: project, isLoading, isError, error } = useQuery({
        queryKey: ["project", id],
        queryFn: () => projectApi.getById(id!),
        enabled: !!id,
    });    

    // TASKS
    const { data: tasks = [], isLoading: tasksLoading } = useQuery({
        queryKey: ["tasks", id, status, assignee],
        queryFn: () =>
            projectApi.getTasks(id!, {
                status: status || undefined,
                assignee: assignee || undefined,
            }),
        enabled: !!id,
    });

    if (isLoading) return <div className="p-6">Loading...</div>;

    if (isError)
        return (
            <div className="p-6 text-red-500">
                {(error as Error).message}
            </div>
        );

    return (
        <div className="p-6 space-y-6 mx-auto">
            <ProjectHeader project={project} />

            <div className="flex items-end justify-between gap-4">
                <ProjectFilters
                    status={status}
                    setStatus={setStatus}
                    assignee={assignee}
                    setAssignee={setAssignee}
                />

                <CreateTaskModal projectId={project?.id} />
            </div>

            {tasksLoading ? (
                <div className="text-gray-500">Loading tasks...</div>
            ) : (
                <TaskList tasks={tasks || []} projectId={project?.id} />
            )}
        </div>
    );
}