import { Outlet, useLocation, useNavigate } from "react-router";
import Header from "./Header";
import type { ActiveTab, Tab } from "./Header";


function getActiveTab(pathname: string): ActiveTab {
    if (pathname == "/") return "home";
    if (pathname.startsWith("/records")) return "records";
    if (pathname.startsWith("/contact")) return "contact";
    return undefined;
}

function tabToPath(tab: Tab): string {
    switch (tab) {
        case "home":
            return "/";
        case "records":
            return "/records";
        case "contact":
            return "/contact";
    }
}

function Layout() {
    const { pathname } = useLocation();
    const navigate = useNavigate();

    const activeTab = getActiveTab(pathname);

    return (
        <div className="min-h-screen bg-gray-50">
            <div className="mx-auto max-w-6xl p-6">
                <div className="rounded-2xl border border-gray-300 bg-white p-6 shadow-sm">
                    <Header
                        activeTab={activeTab}
                        onNavigate={(to: Tab) => navigate(tabToPath(to))}
                    />
                    <main><Outlet /></main>
                </div>
            </div>
        </div>
    )
}

export default Layout;