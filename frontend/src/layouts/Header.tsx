export type Tab = "home" | "records" | "contact";
export type ActiveTab = Tab | undefined;

function Header({
    activeTab,
    onNavigate,
}: {
    activeTab: ActiveTab;
    onNavigate: (to: Tab) => void;
}) {
    const baseBtnStyle = "h-14 flex-1 border border-gray-400 bg-white text-lg font-semibold text-gray-800 hover:bg-gray-50";
    const activeBtnStyle = "bg-gray-100";

    return (
        <div>
            { /* ロゴ */}
            <header className="rounded-t-xl border border-gray-400 bg-white py-10 text-center">
                <div
                    className="cursor-pointer text-4xl font-bold tracking-tight text-gray-900"
                    role="link"
                    tabIndex={0}
                    aria-label="nipopo ホームへ"
                    onClick={() => onNavigate("home")}
                    onKeyDown={(e) => {
                        if (e.key === "Enter" || e.key === " ") onNavigate("home");
                    }}
                >
                    nipopo
                </div>
                <div className="mt-2 text-lg font-semibold text-gray-700">個人日報アプリ</div>
            </header>

            { /* ナビゲーション */}
            <nav className="mb-8 flex rounded-b-xl border-x border-b border-gray-400 bg-white">
                <button
                    type="button"
                    className={`${baseBtnStyle} rounded-bl-xl ${activeTab === "home" ? activeBtnStyle : ""}`}
                    aria-current={activeTab === "home" ? "page" : undefined}
                    onClick={() => onNavigate("home")}
                >
                    ホーム
                </button>
                <button
                    type="button"
                    className={`${baseBtnStyle} border-l-0 border-r-0 ${activeTab === "records" ? activeBtnStyle : ""}`}
                    aria-current={activeTab === "records" ? "page" : undefined}
                    onClick={() => onNavigate("records")}
                >
                    記録一覧
                </button>
                <button
                    type="button"
                    className={`${baseBtnStyle} rounded-br-xl ${activeTab === "contact" ? activeBtnStyle : ""}`}
                    aria-current={activeTab === "contact" ? "page" : undefined}
                    onClick={() => onNavigate("contact")}
                >
                    問い合わせ
                </button>
            </nav>
        </div>
    )
}

export default Header;