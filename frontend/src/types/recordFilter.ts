// 詳細検索時のフィルタ事項。
// 頑張り度の最小値, 最大値については、[effortFrom, effortTo] のものを抽出。
export type RecordFilter = {
    q: string;
    tagIds: string[];
    dateFrom: string;
    dateTo: string;
    effortFrom: number;
    effortTo: number;
};

export const DEFAULT_FILTER: RecordFilter = {
    q: "",
    tagIds: [],
    dateFrom: "",
    dateTo: "",
    effortFrom: 0,
    effortTo: 0,
}