import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { projectApi } from "../api/project.api";
import { CreateProjectModal } from "../components/create-project-modal";
import ProjectRow from "../components/project-row";

export default function ProjectsPage() {
    const queryClient = useQueryClient();

    const { data: projects = [], isLoading } = useQuery({
        queryKey: ["projects"],
        queryFn: projectApi.list,
    });

    const deleteMutation = useMutation({
        mutationFn: projectApi.delete,

        // 🔥 1. Optimistic update
        onMutate: async (projectId: string) => {
            await queryClient.cancelQueries({ queryKey: ["projects"] });

            const previousProjects = queryClient.getQueryData<any[]>(["projects"]);

            queryClient.setQueryData<any[]>(["projects"], (old = []) =>
                old.filter((p) => p.ID !== projectId)
            );

            return { previousProjects };
        },

        // ❌ rollback if error happens
        onError: (_err, _projectId, context) => {
            if (context?.previousProjects) {
                queryClient.setQueryData(["projects"], context.previousProjects);
            }
        },

        // optional: ensure sync
        onSettled: () => {
            queryClient.invalidateQueries({ queryKey: ["projects"] });
        },
    });

    if (isLoading) return <div className="p-4">Loading...</div>;

    return (
        <div className="p-4 space-y-4">
            <div className="flex justify-between items-center">
                <h1 className="text-xl font-semibold">Your Projects</h1>
                <CreateProjectModal />
            </div>

            <div className="grid grid-cols-[1fr_200px_120px] px-4 py-2 text-sm text-gray-500 font-medium">
                <span>Project</span>
                <span>Owner</span>
                <span className="text-right">Created</span>
            </div>

            <div className="space-y-2">
                {projects?.map((project: any) => (
                    <ProjectRow
                        key={project.id}
                        project={project}
                        onDelete={() => deleteMutation.mutate(project?.id)}
                        isDeleting={deleteMutation.isPending}
                    />
                ))}
            </div>
        </div>
    );
}