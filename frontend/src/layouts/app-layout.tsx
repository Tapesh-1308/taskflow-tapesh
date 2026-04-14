import { Outlet } from "react-router-dom"
import { Navbar } from "@/shared/components/navbar"

export function AppLayout() {
    return (
        <div>
            <Navbar />
            <div className="p-4">
                <Outlet />
            </div>
        </div>
    );
}