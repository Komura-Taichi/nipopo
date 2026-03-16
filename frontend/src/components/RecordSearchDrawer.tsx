import { useState } from "react";

import Drawer from "@mui/material/Drawer";

import type { RecordFilter } from "../types/recordFilter";
import { DEFAULT_FILTER } from "../types/recordFilter";
import { primaryBtnStyle, secondaryBtnStyle, errorTextStyle } from "../styles";
import { MAX_EFFORT } from "../constants";
import TagChip from "../components/TagChip";
import TagInput from "../components/TagInput";
import EffortStarsRadio from "../components/EffortStarsRadio";

type Props = {
    open: boolean;
    filter: RecordFilter;
    onClose: () => void;
    onApplyDraft: (nextDraft: RecordFilter) => void;
}

function RecordSearchDrawer({ open, filter, onClose, onApplyDraft }: Props) {
    const [detailInputTag, setDetailInputTag] = useState<string>("");
    const [draft, setDraft] = useState<RecordFilter>(filter);

    const dateError =
        draft.dateFrom && draft.dateTo && draft.dateFrom > draft.dateTo
            ? "検索対象の開始日は終了日以下にしてください。"
            : "";

    const effortError =
        draft.effortFrom > draft.effortTo
            ? "左側の頑張り度は右側の頑張り度以下にしてください。"
            : "";
    const hasError = Boolean(dateError || effortError);

    const onClearDraft = () => {
        setDraft(DEFAULT_FILTER);
    }

    const onAddDetailTag = () => {
        // TODO: tagIdsの中身が現状タグ名となってるが、タグIDに置き換える必要あり。
        const tag = detailInputTag.trim();

        if (!tag) {
            alert("タグ名が空です。タグ名を入力してください。");
            return;
        }

        setDraft((prev) => (prev.tagIds.includes(tag) ? prev : { ...prev, tagIds: [...prev.tagIds, tag] }));
        setDetailInputTag("");
    };

    const applyDraft = () => {
        if (hasError) return;
        onApplyDraft(draft);
        onClose();
    }

    return (
        <Drawer
            anchor="right"
            open={open}
            onClose={onClose}
            ModalProps={{ keepMounted: true }}
            slotProps={{
                paper: {
                    sx: {
                        width: { xs: "90vw", sm: "40vw" },
                        maxWidth: 560,
                    },
                },
            }}
        >
            <div className="h-full flex flex-col">
                {/* ヘッダ */}
                <div className="px-4 py-3 flex items-center justify-between border-b">
                    <h2 className="text-lg font-semibold">詳細検索</h2>
                    <button
                        type="button"
                        className={`${secondaryBtnStyle} rounded`}
                        onClick={onClose}
                        aria-label="詳細検索を閉じる"
                    >
                        閉じる
                    </button>
                </div>

                {/* 検索条件（ボディ） */}
                <div className="flex-1 overflow-auto p-4 space-y-6">
                    <div>
                        <label className="block text-sm text-gray-600">内容</label>
                        <textarea
                            className="mt-1 w-full rounded border border-gray-300 px-3 py-2 text-sm outline-none focus:border-gray-400"
                            rows={3}
                            value={draft.q}
                            onChange={(e) => setDraft((prev) => ({ ...prev, q: e.target.value }))}
                            placeholder="検索内容を入力..."
                        />
                    </div>

                    <div>
                        <label className="block text-sm text-gray-600">タグ（AND検索）</label>

                        <div
                            className="mt-2 flex flex-wrap items-center gap-2"
                            aria-label="検索タグ一覧"
                        >
                            {draft.tagIds.map((t) => (
                                <TagChip
                                    key={t}
                                    tagName={t}
                                    onRemove={() =>
                                        setDraft((prev) => ({ ...prev, tagIds: prev.tagIds.filter((cur_t) => cur_t !== t) }))
                                    }
                                />
                            ))}

                            <TagInput
                                inputTag={detailInputTag}
                                onChange={setDetailInputTag}
                                onAddTag={onAddDetailTag}
                            />
                        </div>
                    </div>

                    {/* 日付範囲 */}
                    <div>
                        <label className="block text-sm text-gray-600">日付の範囲</label>
                        <div className="mt-2 flex gap-2">
                            <input
                                type="date"
                                className="w-full rounded border border-gray-300 px-3 py-2"
                                value={draft.dateFrom}
                                onChange={(e) => setDraft((prev) => ({ ...prev, dateFrom: e.target.value }))}
                            />
                            <input
                                type="date"
                                className="w-full rounded border border-gray-300 px-3 py-2"
                                value={draft.dateTo}
                                onChange={(e) => setDraft((prev) => ({ ...prev, dateTo: e.target.value }))}
                            />
                        </div>

                        {dateError && (
                            <div className={`mt-2 ${errorTextStyle}`}>{dateError}</div>
                        )}
                    </div>

                    {/* 頑張り度 */}
                    <div>
                        <label className="block text-sm text-gray-600">頑張り度の範囲</label>

                        <div className="mt-2 flex items-center gap-3">
                            <span className="text-sm text-gray-600">下限</span>
                            <EffortStarsRadio
                                name="effort_from"
                                effort={draft.effortFrom}
                                maxEffort={MAX_EFFORT}
                                onChange={(n) => setDraft((prev) => ({ ...prev, effortFrom: n }))}
                            />

                            <span className="text-sm text-gray-600">上限</span>
                            <EffortStarsRadio
                                name="effort_to"
                                effort={draft.effortTo}
                                maxEffort={MAX_EFFORT}
                                onChange={(n) => setDraft((prev) => ({ ...prev, effortTo: n }))}
                            />
                        </div>

                        {effortError && (
                            <div className={`mt-2 ${errorTextStyle}`}>{effortError}</div>
                        )}
                    </div>
                </div>

                {/* フッタ */}
                <div className="p-4 border-t flex justify-end gap-2">
                    <button
                        type="button"
                        className={`${secondaryBtnStyle} rounded-lg`}
                        onClick={onClearDraft}
                    >
                        クリア
                    </button>
                    <button
                        type="button"
                        className={`
                ${primaryBtnStyle}
                rounded-lg
                `
                        }
                        disabled={Boolean(hasError)}
                        onClick={applyDraft}
                    >
                        検索
                    </button>
                </div>
            </div>
        </Drawer>
    )
}

export default RecordSearchDrawer;