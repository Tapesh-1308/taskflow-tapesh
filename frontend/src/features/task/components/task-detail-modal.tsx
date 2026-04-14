import { useEffect, useMemo, useState } from "react";
import type { Task, TaskStatus } from "../types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { taskApi } from "../api/task.api";

import { Dialog, DialogContent } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import UserSelect from "@/features/user/components/user-select";

interface Props {
    task: Task;
    open: boolean;
    onOpenChange: (v: boolean) => void;
}

export default function TaskDetailModal({ task, open, onOpenChange }: Props) {
    const queryClient = useQueryClient();

    const [form, setForm] = useState(task);

    const formatDate = (date?: string) => {
        if (!date) return "";
        return new Date(date).toISOString().split("T")[0];
    };

    const toISODate = (date: string) => {
        const [year, month, day] = date.split("-");
        return new Date(Date.UTC(+year, +month - 1, +day)).toISOString();
    };

    useEffect(() => {
        setForm({
            ...task,
            due_date: formatDate(task.due_date),
        });
    }, [task]);


    const changedFields = useMemo(() => {
        const changes: Partial<Task> = {};

        if (form.title !== task.title) changes.title = form.title;

        if (form.description !== task.description)
            changes.description = form.description;

        if (form.status !== task.status) changes.status = form.status;

        if (form.priority !== task.priority)
            changes.priority = form.priority;

        if (form.assignee_id !== task.assignee_id) {
            changes.assignee_id = form.assignee_id;
        }

        if (formatDate(form.due_date) !== formatDate(task.due_date)) {
            changes.due_date = form.due_date;
        }

        return changes;
    }, [form, task]);

    const isDirty = Object.keys(changedFields).length > 0;

    const mutation = useMutation({
        mutationFn: () => {
            const payload: Partial<Task> = {
                ...changedFields,
                ...(changedFields.due_date && {
                    due_date: toISODate(changedFields.due_date),
                }),
            };

            return taskApi.update(task.id, payload);
        },

        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: ["tasks"],
                exact: false,
            });
            onOpenChange(false);
        },
    });

    const updateField = (key: keyof Task, value: any) => {
        setForm((prev) => ({ ...prev, [key]: value }));
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-lg space-y-5">
                {/* TITLE */}
                <Input
                    value={form.title}
                    onChange={(e) => updateField("title", e.target.value)}
                    placeholder="Task title"
                    className="text-lg font-semibold border-none focus-visible:ring-0 px-0"
                />

                {/* DESCRIPTION */}
                <div className="space-y-1">
                    <label className="text-xs text-gray-500">Description</label>
                    <Textarea
                        placeholder="Add description..."
                        value={form.description || ""}
                        onChange={(e) =>
                            updateField("description", e.target.value)
                        }
                    />
                </div>

                {/* FIELDS */}
                <div className="grid grid-cols-2 gap-4 text-sm">
                    {/* STATUS */}
                    <div className="space-y-1">
                        <label className="text-xs text-gray-500">Status</label>
                        <select
                            value={form.status}
                            onChange={(e) =>
                                updateField(
                                    "status",
                                    e.target.value as TaskStatus
                                )
                            }
                            className="border rounded px-2 py-1 w-full"
                        >
                            <option value="todo">To Do</option>
                            <option value="in_progress">In Progress</option>
                            <option value="done">Done</option>
                        </select>
                    </div>

                    {/* PRIORITY */}
                    <div className="space-y-1">
                        <label className="text-xs text-gray-500">
                            Priority
                        </label>
                        <select
                            value={form.priority || "medium"}
                            onChange={(e) =>
                                updateField("priority", e.target.value)
                            }
                            className="border rounded px-2 py-1 w-full"
                        >
                            <option value="low">Low</option>
                            <option value="medium">Medium</option>
                            <option value="high">High</option>
                        </select>
                    </div>

                    {/* ASSIGNEE */}
                    <div className="col-span-2 space-y-1">
                        <label className="text-xs text-gray-500">
                            Assignee
                        </label>
                        <UserSelect
                            value={form.assignee_id || ""}
                            onChange={(userId) => updateField("assignee_id", userId)}
                        />
                    </div>

                    {/* DUE DATE */}
                    <div className="col-span-2 space-y-1">
                        <label className="text-xs text-gray-500">
                            Due Date
                        </label>
                        <input
                            type="date"
                            value={form.due_date || ""}
                            onChange={(e) =>
                                updateField("due_date", e.target.value)
                            }
                            className="border rounded px-2 py-1 w-full"
                        />
                    </div>
                </div>

                {/* META */}
                <div className="text-xs text-gray-400 border-t pt-2">
                    <p>Created: {new Date(task.created_at).toLocaleString()}</p>
                    <p>Updated: {new Date(task.updated_at).toLocaleString()}</p>
                </div>

                {/* ACTION */}
                <Button
                    onClick={() => mutation.mutate()}
                    disabled={!isDirty || mutation.isPending}
                    className="w-full"
                >
                    {mutation.isPending ? "Saving..." : "Save Changes"}
                </Button>
            </DialogContent>
        </Dialog>
    );
}