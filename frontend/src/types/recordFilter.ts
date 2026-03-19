// 詳細検索時のフィルタ事項。
// 頑張り度の最小値, 最大値については、[effortFrom, effortTo] のものを抽出。
export type RecordFilter = {
    q: string;
    tagIds: string[];
    dateFrom: string;
    dateTo: string;
    effortFrom?: number;
    effortTo?: number;
    orderBy: string;
    order: string;
};

export const DEFAULT_FILTER: RecordFilter = {
    q: "",
    tagIds: [],
    dateFrom: "",
    dateTo: "",
    effortFrom: undefined,
    effortTo: undefined,
    orderBy: "date",
    order: "desc",
}