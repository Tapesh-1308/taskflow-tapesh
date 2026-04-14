import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { taskApi } from "../api/task.api";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";

export function CreateTaskModal({ projectId }: { projectId: string }) {
    const [title, setTitle] = useState("");
    const [open, setOpen] = useState(false);

    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: taskApi.create,
        onMutate: async (newTask) => {
            await queryClient.cancelQueries({
                queryKey: ["tasks", projectId],
            });

            const prev = queryClient.getQueryData(["tasks", projectId]);

            queryClient.setQueryData(["tasks", projectId], (old: any) => {
                if (!old) return old;

                return [
                    ...old,
                    {
                        id: "temp-" + Date.now(),
                        title: newTask.title,
                        status: "todo",
                        project_id: projectId,
                    },
                ];
            });

            return { prev };
        },

        onError: (_err, _vars, context) => {
            queryClient.setQueryData(["project", projectId], context?.prev);
        },

        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: ["tasks", projectId],
                exact: false,
            });

            setTitle("");
            setOpen(false);
        },
    });

    const handleSubmit = () => {
        if (!title.trim()) return;

        mutation.mutate({
            title,
            project_id: projectId,
        });
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button>Add Task</Button>
            </DialogTrigger>

            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Create Task</DialogTitle>
                </DialogHeader>

                <div className="space-y-3">
                    <Input
                        placeholder="Task title"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                    />

                    <Button
                        onClick={handleSubmit}
                        disabled={mutation.isPending}
                    >
                        {mutation.isPending ? "Creating..." : "Create"}
                    </Button>
                </div>
            </DialogContent>
        </Dialog>
    );
}
