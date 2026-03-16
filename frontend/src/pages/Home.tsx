import { useEffect, useRef, useState } from 'react';
import { useNavigate } from "react-router";

import { borderStyle, secondaryBtnStyle } from "../styles";
import type { RecentRecord } from "../types/records";
import type { RecordFilter } from '../types/recordFilter';
import { DEFAULT_FILTER } from '../types/recordFilter';
import EffortStarsRadio from '../components/EffortStarsRadio';
import RecentRecordCard from '../components/RecentRecordCard';
import RecordSearchBar from '../components/RecordSearchBar';
import RecordSearchDrawer from "../components/RecordSearchDrawer";
import TagChip from "../components/TagChip";
import TagInput from "../components/TagInput";
import { ROUTES } from '../routes';
import { MIN_EFFORT, MAX_EFFORT } from '../constants';

// 詳細検索後の一覧画面URLのパラメータ部分 (?以降のところ) を作る
function buildDetailSearchParams(filter: RecordFilter): string {
  const params = new URLSearchParams();

  const qNoSpace = filter.q.trim();
  const dateFromNoSpace = filter.dateFrom.trim();
  const dateToNoSpace = filter.dateTo.trim();

  if (qNoSpace) params.set("search_query", qNoSpace);
  for (const tagId of filter.tagIds) params.append("tag_id", tagId);
  if (dateFromNoSpace) params.set("date_from", dateFromNoSpace);
  if (dateToNoSpace) params.set("date_to", dateToNoSpace);
  if (MIN_EFFORT <= filter.effortFrom && filter.effortTo <= MAX_EFFORT) params.set("effort_from", String(filter.effortFrom));
  if (MIN_EFFORT <= filter.effortFrom && filter.effortTo <= MAX_EFFORT) params.set("effort_to", String(filter.effortTo));

  const paramsStr = params.toString();
  return paramsStr ? `?${paramsStr}` : "";
}

function Home() {
  const [searchQ, setSearchQ] = useState<string>("");

  const [filter, setFilter] = useState<RecordFilter>(DEFAULT_FILTER);
  const [drawerOpen, setDrawerOpen] = useState<boolean>(false);

  const [inputTag, setInputTag] = useState<string>("");
  const [tags, setTags] = useState<string[]>(["研究", "勉強"]);

  const [content, setContent] = useState<string>("");

  const [effort, setEffort] = useState<number>(0);

  const [isSaving, setIsSaving] = useState<boolean>(false);

  const [showDialog, setShowDialog] = useState<boolean>(false);

  const [recentRecords, setRecentRecords] = useState<RecentRecord[]>([
    {
      id: "r_1",
      createdAt: "2025/12/11",
      effort: 4,
      tags: ["研究", "勉強"],
      content:
        "今日は、分析データのフィルタリング実装をした。しかし、条件が複雑で実装に苦労した。そこで、先行事例を調査したり、pandasの機能を調べてトライアンドエラーを重ねた。"
    },
  ]);

  const navigate = useNavigate();

  const dialogRef = useRef<HTMLDialogElement>(null);

  const onSearch = () => {
    const q = searchQ.trim();

    // qが空なら絞り込まず一覧へ
    if (!q) {
      navigate(ROUTES.records);
      return;
    }

    // qをURLに乗せる
    const params = new URLSearchParams();
    params.set("search_query", q);
    navigate(`${ROUTES.records}?${params.toString()}`);
  }

  const onAddTag = () => {
    if (!inputTag.trim()) {
      alert("タグが空です。タグ名を入力してください。");
      return;
    }
    setTags((prev: string[]) => [...prev, inputTag]);
    setInputTag("");
  };

  const onSave = () => {
    if (!tags.length || content.trim() == "") {
      setShowDialog(true);
      return;
    }

    setIsSaving(true);

    // 最近のできごとに追加
    const now = new Date();
    const nowFormatted = `${now.getFullYear()}/${now.getMonth() + 1}/${now.getDate()}`;

    const newRecord: RecentRecord = {
      id: `r_${recentRecords.length + 1}`,
      createdAt: nowFormatted,
      effort: effort,
      tags: tags,
      content: content.trim(),
    };

    // 先頭に追加
    setRecentRecords((prev: RecentRecord[]) => [newRecord, ...prev])

    // 入力欄をリセット
    setContent("");
    setEffort(0);
    setTags([]);
    setInputTag("");

    // 記録状態を解除
    setIsSaving(false);
  };

  const onOpenDrawer = () => {
    setDrawerOpen(true);
  };

  const onCloseDrawer = () => {
    setDrawerOpen(false);
  };

  const onApplyDraft = (nextDraft: RecordFilter) => {
    setFilter(nextDraft);

    navigate(`${ROUTES.records}${buildDetailSearchParams(nextDraft)}`);

    setDrawerOpen(false);
  }

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    if (showDialog) dialog.showModal();
    else if (dialog.open) dialog.close();
  }, [showDialog]);

  return (
    <div className="space-y-8">
      {/* レコード検索 */}
      <RecordSearchBar
        searchQ={searchQ}
        onChangeQ={setSearchQ}
        onSearch={onSearch}
        onClickDetailSearch={onOpenDrawer}
      />

      {/* 今日のできごと記入欄 */}
      <section className={borderStyle}>
        {/* 今日のできごと + タグ */}
        <div className="mb-4 flex items-start gap-4">
          <div className="w-28 pt-2 text-base font-semibold text-gray-800">
            今日のできごと
          </div>
          {/* タグ */}
          <div className="flex-1">
            <div className="mb-3 flex flex-wrap items-center gap-2" aria-label="今日のタグ一覧">
              {tags.map((t) => (
                <TagChip
                  key={t}
                  tagName={t}
                  onRemove={() => { setTags((prev) => prev.filter((cur_t) => cur_t !== t)) }}
                />
              ))}

              <TagInput
                inputTag={inputTag}
                onChange={setInputTag}
                onAddTag={onAddTag}
              />
            </div>
          </div>
        </div>
        <textarea
          id="todayContent"
          value={content}
          placeholder="今日一番心に残っていること"
          aria-label="今日一番心に残っていることを入力"
          rows={10}
          onChange={(e) => {
            setContent(e.target.value);
          }}
          className="mb-6 w-full rounded-xl border border-gray-300 p-4 text-sm outline-none focus:border-gray-400"
        />
        {/* 頑張り度 + 記録 */}
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <span className="text-sm font-semibold text-gray-800">頑張り度</span>
            {/* 星 */}
            <EffortStarsRadio
              effort={effort}
              maxEffort={MAX_EFFORT}
              name="effort_new_record"
              onChange={setEffort}
            />

            <span className="text-sm text-gray-700" aria-label="頑張り度の数字表記">
              {`${effort} / ${MAX_EFFORT}`}
            </span>
          </div>

          <button
            type="button"
            disabled={isSaving}
            onClick={onSave}
            className={secondaryBtnStyle}
          >
            {isSaving ? "記録中..." : "記録"}
          </button>
        </div>
      </section>

      {/* 最近のできごと */}
      <section className={borderStyle}>
        <div className="mb-4 text-base font-semibold text-gray-800">最近のできごと</div>

        {recentRecords.length === 0 ? (
          <div className="text-sm text-gray-600">まだ記録がありません</div>
        ) : (
          <div className="space-y-6">
            {recentRecords.map((r) => (
              <div key={r.id} className="space-y-3">
                <RecentRecordCard
                  recentRecord={r}
                  maxEffort={MAX_EFFORT}
                  onClickDetail={(recordId) => navigate(ROUTES.recordDetail(recordId))}
                />
              </div>
            ))}
          </div>
        )}
      </section >

      {/* 入力エラーダイアログ */}
      <dialog
        ref={dialogRef}
        aria-label="入力エラー"
        className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 m-0 rounded-xl border border-gray-300 p-6 backdrop:bg-black/30"
      >
        <h2 className="mb-4 text-base font-semibold text-gray-900">タグ または 内容を入力してください</h2>
        <button
          type="button"
          onClick={() => setShowDialog(false)}
          className={secondaryBtnStyle}
        >
          閉じる
        </button>
      </dialog>

      {/* 詳細検索ドロワー */}
      <RecordSearchDrawer
        open={drawerOpen}
        filter={filter}
        onClose={onCloseDrawer}
        onApplyDraft={onApplyDraft}
      />
    </div >
  );
}

export default Home;
