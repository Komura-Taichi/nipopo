import type React from "react";

import type { RecordFilter } from "../types/recordFilter";
import { errorTextStyle } from "../styles";
import { MAX_EFFORT } from "../constants";
import TagChip from "../components/TagChip";
import TagInput from "../components/TagInput";
import EffortStarsRadio from "../components/EffortStarsRadio";

type Props = {
    draft: RecordFilter;
    setDraft: React.Dispatch<React.SetStateAction<RecordFilter>>;

    inputTag: string;
    setInputTag: React.Dispatch<React.SetStateAction<string>>;
    onAddTag: () => void;

    dateError: string;
    effortError: string;
};

function RecordSearchFilter({
    draft,
    setDraft,
    inputTag,
    setInputTag,
    onAddTag,
    dateError,
    effortError
}: Props) {
    return (
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
                        inputTag={inputTag}
                        onChange={setInputTag}
                        onAddTag={onAddTag}
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
    )
}

export default RecordSearchFilter;