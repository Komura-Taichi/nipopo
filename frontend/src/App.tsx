import { useState } from 'react'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  return (
    <>
      <h1>nipopo</h1>
      <input type="text" placeholder="検索内容を入力..." aria-label="検索内容を入力" />
      <button aria-label="詳細検索">詳細検索</button>
      <label>今日のできごと</label>
      <input type="text" placeholder="タグを入力..." aria-label="タグを入力" />
      <input type="text" placeholder="今日一番心に残っていること" aria-label="今日一番心に残っていることを入力" />
      <label>頑張り度</label>
      <input type="radio" aria-label="1 / 5" />
      <input type="radio" aria-label="2 / 5" />
      <input type="radio" aria-label="3 / 5" />
      <input type="radio" aria-label="4 / 5" />
      <input type="radio" aria-label="5 / 5" />
      <label aria-label="頑張り度の数字表記">0 / 5</label>
      <button aria-label="記録">記録</button>
      <label>最近のできごと</label>
      <div>
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
  )
}

export default App
