import { primaryBtnStyle, secondaryBtnStyle } from "../styles";

function RecordSearchBar() {
    return (
        <div className="flex items-center">
            <div className="flex flex-1">
                <input
                    type="text"
                    placeholder="検索内容を入力..."
                    aria-label="検索内容を入力"
                    className="h-10 w-full rounded-l-lg border border-gray-300 px-3 text-sm outline-none focus:border-gray-400"
                />
                <button
                    type="button"
                    className={`${primaryBtnStyle} -ml-px rounded-r-lg`}
                >
                    検索
                </button>
            </div>

            <div className="ml-6">
                <button type="button" className={`${secondaryBtnStyle} rounded-lg`}>
                    詳細検索
                </button>
            </div>
        </div>
    );
}

export default RecordSearchBar;