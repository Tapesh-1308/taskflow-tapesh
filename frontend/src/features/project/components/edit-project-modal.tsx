import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";

import { projectApi } from "../api/project.api";

import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const schema = z.object({
    name: z.string().min(3),
    description: z.string().optional(),
});

export function EditProjectModal({ project }: any) {
    const [open, setOpen] = useState(false);
    const queryClient = useQueryClient();

    const { register, handleSubmit } = useForm({
        resolver: zodResolver(schema),
        defaultValues: {
            name: project.Name,
            description: project.Description,
        },
    });

    const mutation = useMutation({
        mutationFn: (data: any) =>
            projectApi.update(project.ID, data),

        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["projects"] });
            setOpen(false);
        },
    });

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <button className="w-full text-left px-2 py-1 hover:bg-gray-100">
                    Edit
                </button>
            </DialogTrigger>

            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Edit Project</DialogTitle>
                </DialogHeader>

                <form
                    onSubmit={handleSubmit((data) => mutation.mutate(data))}
                    className="space-y-4"
                >
                    <Input {...register("name")} />
                    <Input {...register("description")} />

                    <Button className="w-full">
                        {mutation.isPending ? "Saving..." : "Save"}
                    </Button>
                </form>
            </DialogContent>
        </Dialog>
    );
}