const COLORS = [
    "bg-red-500",
    "bg-blue-500",
    "bg-green-500",
    "bg-purple-500",
    "bg-pink-500",
    "bg-yellow-500",
    "bg-indigo-500",
    "bg-teal-500",
];

function isValidName(name: unknown): name is string {
    return typeof name === "string" && name.trim().length > 0;
}

function getInitials(name: string) {
    const parts = name.trim().split(" ");

    if (parts.length === 1) {
        return parts[0][0].toUpperCase();
    }

    return (parts[0][0] + parts[1][0]).toUpperCase();
}

export function getAvatarProps(name: unknown) {
    if (!isValidName(name)) {
        return {
            color: "bg-gray-400",
            letter: "?",
        };
    }

    const cleanName = name.trim();

    return {
        color: COLORS[cleanName.length % COLORS.length],
        letter: getInitials(cleanName),
    };
}