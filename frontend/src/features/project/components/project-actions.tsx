import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import { Button } from "@/components/ui/button";
import { useState } from "react";
import DeleteProjectDialog from "./delete-project-dialog";
import { EditProjectModal } from "./edit-project-modal";

export default function ProjectActions({ project, onDelete }: any) {
    const [open, setOpen] = useState(false);

    return (
        <>
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="px-2">
                        ⋮
                    </Button>
                </DropdownMenuTrigger>

                {/* IMPORTANT: z-50 fixes overflow behind card */}
                <DropdownMenuContent align="end" className="w-40 z-50">
                    <DropdownMenuItem asChild>
                        <EditProjectModal project={project} />
                    </DropdownMenuItem>

                    <DropdownMenuItem
                        className="text-red-500 focus:text-red-500"
                        onClick={() => setOpen(true)}
                    >
                        Delete
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>

            {/* Confirm Dialog */}
            <DeleteProjectDialog
                open={open}
                onOpenChange={setOpen}
                onConfirm={() => onDelete(project.ID)}
            />
        </>
    );
}