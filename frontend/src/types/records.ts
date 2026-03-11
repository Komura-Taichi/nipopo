// 仮の「最近のできごと」。
// 本来、tagsの要素も型定義する。
export type RecentRecord = {
    id: string;
    createdAt: string;
    effort: number;
    tags: string[];
    content: string;
};