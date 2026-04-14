import { useState } from "react";

import type { Task, TaskStatus, TaskWithUser } from "../types";

import { Card, CardContent } from "@/components/ui/card";
import { getAvatarProps } from "@/shared/lib/avtar";
import { timeAgo } from "@/shared/lib/time";

import TaskDetailModal from "./task-detail-modal";

export function TaskItem({ task }: { task: TaskWithUser }) {
    const [open, setOpen] = useState(false);

    const avatar = getAvatarProps(task?.assignee?.name);

    return (
        <>
            <Card
                onClick={() => setOpen(true)}
                className="hover:shadow-sm transition cursor-pointer"
            >
                <CardContent className="p-3 space-y-2">
                    {/* TITLE */}
                    <p className="font-medium text-sm leading-tight">
                        {task.title}
                    </p>

                    {/* META ROW */}
                    <div className="flex items-center justify-between text-xs text-gray-500">
                        {/* LEFT */}
                        <div className="flex items-center gap-2">
                            <PriorityBadge priority={task.priority} />
                        </div>

                        {/* RIGHT */}
                        <div className="flex items-center gap-2">
                            <span>{timeAgo(task.created_at)}</span>

                            <div
                                className={`h-6 w-6 rounded-full flex items-center justify-center text-white text-[10px] ${avatar.color}`}
                            >
                                {avatar.letter}
                            </div>
                        </div>
                    </div>
                </CardContent>
            </Card>

            {/* MODAL */}
            <TaskDetailModal
                task={task}
                open={open}
                onOpenChange={setOpen}
            />
        </>
    );
}

function PriorityBadge({ priority }: { priority?: string }) {
    if (!priority) return null;

    const styles: Record<string, string> = {
        low: "bg-gray-200 text-gray-700",
        medium: "bg-yellow-200 text-yellow-800",
        high: "bg-red-200 text-red-700",
    };

    return (
        <span className={`px-2 py-0.5 rounded text-[10px] ${styles[priority]}`}>
            {priority}
        </span>
    );
}

function StatusDot({ status }: { status: TaskStatus }) {
    const colors: Record<TaskStatus, string> = {
        todo: "bg-gray-400",
        in_progress: "bg-blue-500",
        done: "bg-green-500",
    };

    return <div className={`h-2 w-2 rounded-full ${colors[status]}`} />;
}
