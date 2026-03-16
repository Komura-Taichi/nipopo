import type { Tag } from "../types/tags";

// 仮の「最近のできごと」。
export type RecentRecord = {
    id: string;
    createdAt: string;
    effort: number;
    tags: Tag[];
    content: string;
};