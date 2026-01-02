import { useState } from 'react';
import './App.css';

function App() {
  const [tags, setTags] = useState<string[]>(["研究", "勉強"]);

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
            <button aria-label={`タグ ${t} を削除`}>×</button>
          </>
        ))}
        <input type="text" placeholder="タグを入力..." aria-label="タグを入力" />
        <button aria-label="タグを追加">+</button>
      </div>
      <textarea id="todayContent" placeholder="今日一番心に残っていること" aria-label="今日一番心に残っていることを入力" rows={10} />
      <label>頑張り度</label>
      <input type="radio" name="effort" aria-label="1 / 5" />
      <input type="radio" name="effort" aria-label="2 / 5" />
      <input type="radio" name="effort" aria-label="3 / 5" />
      <input type="radio" name="effort" aria-label="4 / 5" />
      <input type="radio" name="effort" aria-label="5 / 5" />
      <span aria-label="頑張り度の数字表記">0 / 5</span>
      <button>記録</button>
      <label>最近のできごと</label>
      <div aria-label="最近のできごと">
        <h2>2025/12/11</h2>
        { /* 星 */ }
        <img />
        <img />
        <img />
        <img />
        <img />
        <label>研究</label>
        <label>勉強</label>
        <p>
          今日は、分析データのフィルタリング実装をした。しかし、条件が複雑で実装に苦労した。
          そこで、先行事例を調査したり、pandasの機能を調べてトライアンドエラーを重ねた。
        </p>
      </div>
    </>
  );
}

export default App;
