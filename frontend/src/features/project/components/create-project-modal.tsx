// src/features/project/components/create-project-modal.tsx

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
    name: z.string().min(3, "Name must be at least 3 characters"),
    description: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

export function CreateProjectModal() {
    const [open, setOpen] = useState(false);
    const queryClient = useQueryClient();

    const {
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm<FormData>({
        resolver: zodResolver(schema),
    });
    const mutation = useMutation({
        mutationFn: projectApi.create,

        onMutate: async (newProject) => {
            await queryClient.cancelQueries({ queryKey: ["projects"] });

            const prev = queryClient.getQueryData(["projects"]);

            queryClient.setQueryData(["projects"], (old: Array<Object> = []) =>
                old?.length
                    ? [...old, { id: "temp", ...newProject }]
                    : [{ id: "temp", ...newProject }],
            );

            return { prev };
        },

        onError: (_err, _new, context) => {
            queryClient.setQueryData(["projects"], context?.prev);
        },

        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["projects"] });
            reset();
            setOpen(false);
        },
    });

    const onSubmit = (data: FormData) => {
        mutation.mutate(data);
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button className="cursor-pointer">Create Project</Button>
            </DialogTrigger>

            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Create Project</DialogTitle>
                </DialogHeader>

                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                    <div>
                        <Input
                            placeholder="Project name"
                            {...register("name")}
                        />
                        {errors.name && (
                            <p className="text-sm text-red-500">
                                {errors.name.message}
                            </p>
                        )}
                    </div>

                    <div>
                        <Input
                            placeholder="Description"
                            {...register("description")}
                        />
                    </div>

                    {mutation.isError && (
                        <p className="text-sm text-red-500">
                            {(mutation.error as Error).message}
                        </p>
                    )}

                    <Button
                        type="submit"
                        disabled={mutation.isPending}
                        className="w-full cursor-pointer"
                    >
                        {mutation.isPending ? "Creating..." : "Create"}
                    </Button>
                </form>
            </DialogContent>
        </Dialog>
    );
}
