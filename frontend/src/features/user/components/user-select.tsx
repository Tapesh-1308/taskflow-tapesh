import { useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { userApi } from "../api/user.api";

import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

import { getAvatarProps } from "@/shared/lib/avtar";
import type { User } from "../type";

interface Props {
    value?: string; // user_id
    onChange: (userId: string, user?: User) => void;
}

export default function UserSelect({ value, onChange }: Props) {
    const [search, setSearch] = useState("");

    const { data: users = [], isLoading } = useQuery({
        queryKey: ["users", search],
        queryFn: () => userApi.search(search),
    });

    // debounce search
    useEffect(() => {
        const t = setTimeout(() => {}, 300);
        return () => clearTimeout(t);
    }, [search]);

    return (
        <Select
            value={value || ""}
            onValueChange={(val) => {
                const selectedUser = users.find((u) => u.id === val);
                onChange(val, selectedUser);
            }}
        >
            <SelectTrigger>
                <SelectValue placeholder="Select user" />
            </SelectTrigger>

            <SelectContent>
                {/* SEARCH INPUT */}
                <div className="p-2">
                    <input
                        placeholder="Search user..."
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        className="w-full border px-2 py-1 rounded text-sm"
                    />
                </div>

                {/* USERS */}
                {isLoading && (
                    <div className="px-2 py-1 text-sm text-gray-400">
                        Loading...
                    </div>
                )}

                {users.map((user) => {
                    const avatar = getAvatarProps(user.name);

                    return (
                        <SelectItem key={user.id} value={user.id}>
                            <div className="flex items-center gap-2">
                                <div
                                    className={`w-6 h-6 rounded-full flex items-center justify-center text-white text-xs ${avatar.color}`}
                                >
                                    {avatar.letter}
                                </div>
                                <span>{user.name}</span>
                            </div>
                        </SelectItem>
                    );
                })}

                {users.length === 0 && !isLoading && (
                    <div className="px-2 py-1 text-sm text-gray-400">
                        No users found
                    </div>
                )}
            </SelectContent>
        </Select>
    );
}