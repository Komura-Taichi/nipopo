import React, { useEffect, useRef, useState } from 'react';

// 仮の「最近のできごと」。
// 本来、tagsの要素も型定義する。
type RecentRecord = {
  id: string;
  createdAt: string;
  effort: number;
  tags: string[];
  content: string;
};

function Home() {
  const maxEffort = 5;
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

  const dialogRef = useRef<HTMLDialogElement>(null);

  const baseBtnStyle
    = "h-10 border px-4 whitespace-nowrap text-sm font-semibold focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400";
  const primaryBtnStyle
    = `${baseBtnStyle} border-sky-300 bg-sky-500 text-white hover:bg-sky-400`;
  const secondaryBtnStyle
    = `${baseBtnStyle} border-gray-300 bg-white text-gray-800 hover:bg-gray-50`;
  const starStyle
    = "text-3xl leading-none"

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
      id: `t_${recentRecords.length + 1}`,
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
      <div className="flex items-center">
        <div className="flex flex-1">
          <input
            type="text"
            placeholder="検索内容を入力..."
            aria-label="検索内容を入力"
            className="h-10 w-full rounded-l-lg border border-gray-300 px-3 text-sm outline-none focus:border-gray-400"
          />
          <button
            type="button"
            className={`${primaryBtnStyle} -ml-px rounded-r-lg`}
          >
            検索
          </button>
        </div>

        <div className="ml-6">
          <button type="button" className={`${secondaryBtnStyle} rounded-lg`}>
            詳細検索
          </button>
        </div>
      </div>

      {/* 今日のできごと記入欄 */}
      <section className="rounded-xl border border-gray-300 p-6">
        {/* 今日のできごと + タグ */}
        <div className="mb-4 flex items-start gap-4">
          <div className="w-28 pt-2 text-base font-semibold text-gray-800">
            今日のできごと
          </div>
          {/* タグ */}
          <div className="flex-1">
            <div className="mb-3 flex flex-wrap items-center gap-2" aria-label="今日のタグ一覧">
              {tags.map((t) => (
                <span
                  key={t}
                  className="inline-flex items-center gap-2 rounded-full border border-gray-300 bg-white px-3 py-1 text-sm text-gray-800"
                >
                  {t}
                  <button
                    type="button"
                    aria-label={`タグ 「${t}」 を削除`}
                    onClick={() => {
                      setTags(prev => prev.filter((cur_t) => cur_t !== t))
                    }}
                    className="rounded-full px-1 text-gray-500 hover:bg-gray-100"
                  >
                    ×
                  </button>
                </span>
              ))}

              <div className="flex">
                <input
                  type="text"
                  value={inputTag}
                  placeholder="タグを入力..."
                  aria-label="タグを入力"
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                    setInputTag(e.target.value);
                  }}
                  className="h-9 w-56 rounded-l-lg border border-gray-300 px-3 text-sm outline-none focus:border-gray-400"
                />
                <button
                  aria-label="タグを追加"
                  onClick={onAddTag}
                  className="-ml-px h-9 rounded-r-lg border border-green-300 bg-green-400 px-4 text-lg font-semibold text-gray-800 hover:bg-green-500"
                >
                  +
                </button>
              </div>
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
            <div className="flex items-center gap-2">
              {[1, 2, 3, 4, 5].map((n) => (
                <label
                  key={n}
                  className="cursor-pointer select-none"
                >
                  <input
                    type="radio"
                    name="effort"
                    className="sr-only"
                    aria-label={`${n} / ${maxEffort}`}
                    checked={effort === n}
                    onChange={() => setEffort(n)}
                  />
                  <span className={`${starStyle} ${n <= effort ? "text-yellow-400" : "text-gray-300"}`}>
                    {n <= effort ? "★" : "☆"}
                  </span>
                </label>
              ))}
            </div>

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

      <section className="rounded-xl border border-gray-300 p-6">
        <div className="mb-4 text-base font-semibold text-gray-800">最近のできごと</div>

        {recentRecords.length === 0 ? (
          <div className="text-sm text-gray-600">まだ記録がありません</div>
        ) : (
          <div className="space-y-6">
            {recentRecords.map((r) => (
              <div key={r.id} className="space-y-3">
                <div className="text-lg font-semibold text-gray-900">{r.createdAt}</div>

                {/* TODO: 星をコンポーネント化して置き換え */}
                <div className="flex gap-2" aria-label="頑張り度（星）">
                  {[1, 2, 3, 4, 5].map((n) => (
                    <span key={n} className={`${starStyle} ${n <= r.effort ? "text-yellow-400" : "text-gray-300"}`}>
                      {n <= r.effort ? "★" : "☆"}
                    </span>
                  ))}
                </div>

                {/* タグ */}
                <div className="flex flex-wrap gap-2">
                  {r.tags.map((t) => (
                    <span
                      key={`${r.id}_${t}`}
                      className="rounded-full border border-gray-300 bg-white px-3 py-1 text-sm"
                    >
                      {t}
                    </span>
                  ))}
                </div>

                <p className="text-sm leading-6 text-gray-800">{r.content}</p>
              </div>
            ))}
          </div>
        )}
      </section >
    </div >
  );
}

export default Home;
