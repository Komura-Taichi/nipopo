import { useState } from "react";

import Drawer from "@mui/material/Drawer";

import type { RecordFilter } from "../types/recordFilter";
import { DEFAULT_FILTER } from "../types/recordFilter";
import { primaryBtnStyle, secondaryBtnStyle } from "../styles";
import RecordSearchFilter from "./RecordSearchFilter";

type Props = {
    open: boolean;
    filter: RecordFilter;
    onClose: () => void;
    onApplyDraft: (nextDraft: RecordFilter) => void;
}

function RecordSearchDrawer({ open, filter, onClose, onApplyDraft }: Props) {
    const [inputTag, setInputTag] = useState<string>("");
    const [draft, setDraft] = useState<RecordFilter>(filter);

    const dateError =
        draft.dateFrom && draft.dateTo && draft.dateFrom > draft.dateTo
            ? "検索対象の開始日は終了日以下にしてください。"
            : "";

    const effortError =
        draft.effortFrom !== undefined &&
            draft.effortTo !== undefined &&
            draft.effortFrom > draft.effortTo
            ? "左側の頑張り度は右側の頑張り度以下にしてください。"
            : "";
    const hasError = Boolean(dateError || effortError);

    const onClearDraft = () => {
        setDraft(DEFAULT_FILTER);
    }

    const onAddTag = () => {
        const tag = inputTag.trim();

        if (!tag) {
            alert("タグ名が空です。タグ名を入力してください。");
            return;
        }

        setDraft((prev) => (prev.tagIds.includes(tag) ? prev : { ...prev, tagIds: [...prev.tagIds, tag] }));
        setInputTag("");
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
                <RecordSearchFilter
                    draft={draft}
                    setDraft={setDraft}
                    inputTag={inputTag}
                    setInputTag={setInputTag}
                    onAddTag={onAddTag}
                    dateError={dateError}
                    effortError={effortError}
                />

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