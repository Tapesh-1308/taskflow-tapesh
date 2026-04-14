import type { Task, TaskStatus } from "../types";
import { TaskItem } from "./task-item";

interface Props {
    tasks: Task[];
    projectId: string;
}

const statuses: TaskStatus[] = ["todo", "in_progress", "done"];

const statusLabels: Record<TaskStatus, string> = {
    todo: "To Do",
    in_progress: "In Progress",
    done: "Done",
};

export function TaskList({ tasks }: Props) {
    return (
        <div className="overflow-x-auto">
            <div className="flex gap-4 min-w-225">
                {statuses.map((status) => {
                    const filtered = tasks.filter((t) => t.status === status);

                    return (
                        <div
                            key={status}
                            className="flex-1 min-w-70 bg-gray-50 rounded-lg p-3"
                        >
                            {/* HEADER */}
                            <div className="flex items-center justify-between mb-3 sticky top-0 bg-gray-50 z-10">
                                <h2 className="font-medium text-sm">
                                    {statusLabels[status]}
                                </h2>

                                <span className="text-xs bg-gray-200 px-2 py-0.5 rounded">
                                    {filtered.length}
                                </span>
                            </div>

                            {/* TASKS */}
                            <div className="space-y-2">
                                {filtered.length === 0 ? (
                                    <EmptyState status={status} />
                                ) : (
                                    filtered.map((task) => (
                                        <TaskItem key={task.id} task={task} />
                                    ))
                                )}
                            </div>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}

function EmptyState({ status }: { status: TaskStatus }) {
    const messages: Record<TaskStatus, string> = {
        todo: "No tasks to start",
        in_progress: "Nothing in progress",
        done: "No completed tasks yet",
    };

    return (
        <div className="text-sm text-gray-400 border border-dashed rounded p-3 text-center">
            {messages[status]}
        </div>
    );
}