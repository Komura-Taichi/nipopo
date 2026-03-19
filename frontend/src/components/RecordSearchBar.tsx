import type React from "react";
import { primaryBtnStyle, secondaryBtnStyle } from "../styles";
import type { RecordFilter } from "../types/recordFilter";

type Props = {
    filter: RecordFilter;
    onPatchFilter: (patch: Partial<RecordFilter>) => void;

    searchQ: string;
    onChangeQ: (searchQ: string) => void;
    onSearch: () => void;
    onClickDetailSearch: () => void;
}

function RecordSearchBar({
    filter,
    onPatchFilter,
    searchQ,
    onChangeQ,
    onSearch,
    onClickDetailSearch,
}: Props) {
    const toggleOrder = () => {
        const next = filter.order === "asc" ? "desc" : "asc";
        onPatchFilter({ order: next });
    };

    return (
        <div className="flex items-center gap-6">
            {/*簡易検索 */}
            <div className="flex flex-1">
                <input
                    type="text"
                    value={searchQ}
                    placeholder="検索内容を入力..."
                    aria-label="検索内容を入力"
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                        onChangeQ(e.target.value);
                    }}
                    className="h-10 w-full rounded-l-lg border border-gray-300 px-3 text-sm outline-none focus:border-gray-400"
                />
                <button
                    type="button"
                    className={`${primaryBtnStyle} -ml-px rounded-r-lg`}
                    onClick={onSearch}
                >
                    検索
                </button>
            </div>
            {/* 並び替え */}
            <div className="flex items-center gap-2">
                <span className="text-sm text-gray-600">並び替え</span>

                <select
                    aria-label="並び替え方法"
                    value={filter.orderBy}
                    onChange={(e) => onPatchFilter({ orderBy: e.target.value })}
                    className="h-10 rounded-lg border border-gray-300 bg-white px-3 text-sm outline-none focus:border-gray-400"
                >
                    <option value="date">日付</option>
                    <option value="effort">頑張り度</option>
                </select>

                <button
                    type="button"
                    aria-label="並び順を切り替え"
                    onClick={toggleOrder}
                    className={`${secondaryBtnStyle} rounded-lg`}
                    title={filter.order === "asc" ? "昇順" : "降順"}
                >
                    {filter.order === "asc" ? "昇順" : "降順"}
                </button>
            </div>

            {/* 詳細検索 */}
            <div className="ml-6">
                <button
                    type="button"
                    className={`${secondaryBtnStyle} rounded-lg`}
                    onClick={onClickDetailSearch}
                >
                    詳細検索
                </button>
            </div>
        </div>
    );
}

export default RecordSearchBar;