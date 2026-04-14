import { useEffect, useState } from "react";

import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

import { Button } from "@/components/ui/button";
import UserSelect from "@/features/user/components/user-select";

interface Props {
    status: string;
    setStatus: (v: string) => void;
    assignee: string;
    setAssignee: (v: string) => void;
}

export default function ProjectFilters({
    status,
    setStatus,
    assignee,
    setAssignee,
}: Props) {
    const [search, setSearch] = useState(assignee);

    useEffect(() => {
        const timer = setTimeout(() => {
            setAssignee(search);
        }, 400);

        return () => clearTimeout(timer);
    }, [search, setAssignee]);

    return (
        <div className="flex items-end gap-4">
            {/* STATUS */}
            <div className="flex flex-col gap-1 w-40">
                <label className="text-xs text-muted-foreground">Status</label>

                <Select
                    value={status || "all"}
                    onValueChange={(value) =>
                        setStatus(value === "all" ? "" : value)
                    }
                >
                    <SelectTrigger>
                        <SelectValue placeholder="All" />
                    </SelectTrigger>

                    <SelectContent>
                        <SelectItem value="all">All</SelectItem>
                        <SelectItem value="todo">To Do</SelectItem>
                        <SelectItem value="in_progress">In Progress</SelectItem>
                        <SelectItem value="done">Done</SelectItem>
                    </SelectContent>
                </Select>
            </div>

            {/* ASSIGNEE */}
            <div className="flex flex-col gap-1 w-50">
                <label className="text-xs text-muted-foreground">
                    Assignee
                </label>
                <UserSelect
                    value={assignee}
                    onChange={(userId) => setAssignee(userId)}
                />
            </div>

            {/* RESET */}
            <Button
                variant="ghost"
                onClick={() => {
                    setStatus("");
                    setSearch("");
                    setAssignee("");
                }}
            >
                Reset
            </Button>
        </div>
    );
}
