import { useNavigate } from "react-router-dom";
import { Card, CardContent } from "@/components/ui/card";
import { getAvatarProps } from "@/shared/lib/avtar";
import { timeAgo } from "@/shared/lib/time";
import ProjectActions from "./project-actions";

export default function ProjectRow({ project, onDelete }: any) {
    const navigate = useNavigate();
    const avatar = getAvatarProps(project?.owner?.name);

    return (
        <Card className="hover:bg-gray-50 transition">
            <CardContent className="p-4 grid grid-cols-[1fr_200px_120px] items-center">
                {/* Project */}
                <div
                    className="cursor-pointer"
                    onClick={() => navigate(`/projects/${project.id}`)}
                >
                    <p className="font-medium">{project.name}</p>
                    <p className="text-sm text-gray-500">
                        {project.description || "No description"}
                    </p>
                </div>

                {/* Owner */}
                <div className="flex items-center gap-2">
                    <div
                        className={`h-8 w-8 rounded-full flex items-center justify-center text-white text-sm ${avatar.color}`}
                    >
                        {avatar.letter}
                    </div>
                    <span className="text-sm text-gray-600 truncate">
                        {project?.owner?.name}
                    </span>
                </div>

                {/* Actions */}
                <div className="flex items-center justify-end gap-3">
                    <span className="text-sm text-gray-400">
                        {timeAgo(project.created_at)}
                    </span>

                    <ProjectActions project={project} onDelete={onDelete} />
                </div>
            </CardContent>
        </Card>
    );
}
