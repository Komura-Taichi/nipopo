import React, { useEffect, useRef, useState } from 'react';
import { useNavigate } from "react-router";

import { borderStyle, secondaryBtnStyle } from "../styles";
import type { RecentRecord } from "../types";
import EffortStarsRadio from '../components/EffortStarsRadio';
import RecentRecordCard from '../components/RecentRecordCard';
import RecordSearchBar from '../components/RecordSearchBar';
import TagChip from "../components/TagChip";
import TagInput from "../components/TagInput";
import { ROUTES } from '../routes';

function Home() {
  const maxEffort = 5;
  const [searchQ, setSearchQ] = useState<string>("");

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
          onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => {
            setContent(e.target.value);
          }}
          className="mb-6 w-full rounded-xl border border-gray-300 p-4 text-sm outline-none focus:border-gray-400"
        />
        {/* 頑張り度 + 記録 */}
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <span className="text-sm font-semibold text-gray-800">頑張り度</span>
            {/* 星 */}
            <EffortStarsRadio effort={effort} maxEffort={maxEffort} onChange={setEffort} />

            <span className="text-sm text-gray-700" aria-label="頑張り度の数字表記">
              {`${effort} / ${maxEffort}`}
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
                  maxEffort={maxEffort}
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
    </div >
  );
}

export default Home;
