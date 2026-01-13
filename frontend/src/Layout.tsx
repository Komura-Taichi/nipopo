import { Outlet, useLocation, useNavigate } from "react-router";

type ActiveTab = "home" | "records" | "contact" | undefined;

function getActiveTab(pathname: string): ActiveTab {
    if (pathname == "/") return "home";
    if (pathname.startsWith("/records")) return "records";
    if (pathname.startsWith("/contact")) return "contact";
    return undefined;
}

function Layout() {
    const { pathname } = useLocation();
    const navigate = useNavigate();

    const activeTab = getActiveTab(pathname);

    const baseButtonStyle = "h-14 flex-1 border border-gray-400 bg-white text-lg font-semibold text-gray-800 hover:bg-gray-50";
    const activeButtonStyle = "bg-gray-100";

    return (
        <div className="min-h-screen bg-gray-50">
            <div className="mx-auto max-w-6xl p-6">
                <div className="rounded-2xl border border-gray-300 bg-white p-6 shadow-sm">
                    { /* ロゴ */}
                    <header className="rounded-t-xl border border-gray-400 bg-white py-10 text-center">
                        <div 
                            className="cursor-pointer text-4xl font-bold tracking-tight text-gray-900"
                            role="link"
                            tabIndex={0}
                            aria-label="nipopo ホームへ"
                            onClick={() => navigate("/")}
                            onKeyDown={(e) => {
                                if (e.key === "Enter" || e.key === " ") navigate("/");
                            }}
                        >
                            nipopo
                        </div>
                        <div className="mt-2 text-lg font-semibold text-gray-700">個人日報アプリ</div>
                    </header>

                    { /* ナビゲーション */ }
                    <nav className="mb-8 flex rounded-b-xl border-x border-b border-gray-400 bg-white">
                        <button
                            type="button"
                            className={`${baseButtonStyle} rounded-bl-xl ${activeTab === "home" ? activeButtonStyle: ""}`}
                            aria-current={activeTab === "home" ? "page" : undefined}
                            onClick={() => navigate("/")}
                        >
                            ホーム
                        </button>
                        <button
                            type="button"
                            className={`${baseButtonStyle} border-l-0 border-r-0 ${activeTab === "records" ? activeButtonStyle: ""}`}
                            aria-current={activeTab === "records" ? "page" : undefined}
                            onClick={() => navigate("/records")}
                        >
                            記録一覧
                        </button>
                        <button
                            type="button"
                            className={`${baseButtonStyle} rounded-br-xl ${activeTab === "contact" ? activeButtonStyle: ""}`}
                            aria-current={activeTab === "contact" ? "page" : undefined}
                            onClick={() => navigate("/contact")}
                        >
                            問い合わせ
                        </button>
                    </nav>

                    <main><Outlet /></main>
                </div>
            </div>
        </div>
    )
}

export default Layout;