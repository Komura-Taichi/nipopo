import React, { useEffect, useRef, useState } from 'react';
import './App.css';

function App() {
  const maxEffort = 5;
  const [inputTag, setInputTag] = useState<string>("");
  const [tags, setTags] = useState<string[]>(["研究", "勉強"]);

  const [content, setContent] = useState<string>("");

  const [effort, setEffort] = useState<number>(0);

  const [isSaving, setIsSaving] = useState<boolean>(false);

  const [showDialog, setShowDialog] = useState<boolean>(false);

  const dialogRef = useRef<HTMLDialogElement>(null);

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
  };

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    if (showDialog) dialog.showModal();
    else if (dialog.open) dialog.close();
  }, [showDialog]);

  return (
    <>
      <h1>nipopo</h1>
      <input type="text" placeholder="検索内容を入力..." aria-label="検索内容を入力" />
      <button aria-label="詳細検索">詳細検索</button>

      <label htmlFor="todayContent">今日のできごと</label>
      <div aria-label="今日のタグ一覧">
        {tags.map((t) => (
          <>
            <span key={t}>{t}</span>
            <button
              key={t}
              aria-label={`タグ 「${t}」 を削除`}
              onClick={() => {
                setTags(prev => prev.filter((cur_t) => cur_t !== t))
              }}
            >
              ×
            </button>
          </>
        ))}
        <input 
          type="text"
          value={inputTag}
          placeholder="タグを入力..."
          aria-label="タグを入力"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setInputTag(e.target.value);
          }}
        />
        <button
          aria-label="タグを追加"
          onClick={onAddTag}
        >+</button>
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
      />
      <label>頑張り度</label>
      {[1, 2, 3, 4, 5].map((n) => (
        <input
          type="radio"
          key={`${n} / ${maxEffort}`}
          name="effort"
          aria-label={`${n} / ${maxEffort}`}
          checked={effort === n}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setEffort(n);
          }}
        />
      ))}
      <span aria-label="頑張り度の数字表記">{`${effort} / ${maxEffort}`}</span>
      <button
        disabled={isSaving}
        onClick={onSave}
      >
        {isSaving ? "記録中..." : "記録"}
      </button>

      <dialog ref={dialogRef} aria-label="入力エラー">
        <h2>タグ または 内容を入力してください</h2>
        <button onClick={() => setShowDialog(false)}>閉じる</button>
      </dialog>

      <label>最近のできごと</label>
      <div aria-label="最近のできごと">
        <h2>2025/12/11</h2>
        { /* 星 */ }
        <img />
        <img />
        <img />
        <img />
        <img />
        <span>研究</span>
        <span>勉強</span>
        <p>
          今日は、分析データのフィルタリング実装をした。しかし、条件が複雑で実装に苦労した。
          そこで、先行事例を調査したり、pandasの機能を調べてトライアンドエラーを重ねた。
        </p>
      </div>
    </>
  );
}

export default App;
