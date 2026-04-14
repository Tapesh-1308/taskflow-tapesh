interface Props {
    project: any;
}

export default function ProjectHeader({ project }: Props) {
    return (
        <div className="flex justify-between items-start">
            <div className="space-y-1">
                <h1 className="text-2xl font-semibold">
                    {project.name}
                </h1>
                <p className="text-gray-500">
                    {project.description}
                </p>
            </div>

            <div className="text-sm text-gray-500 text-right space-y-1">
                <p>Owner: {project.owner.name}</p>
                <p>
                    Created:{" "}
                    {new Date(
                        project.created_at
                    ).toLocaleString()}
                </p>
            </div>
        </div>
    );
}