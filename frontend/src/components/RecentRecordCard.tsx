import { primaryBtnStyle } from "../styles";
import type { RecentRecord } from "../types";
import EffortStarsDisplay from "./EffortStarsDisplay";
import TagChip from "./TagChip";

type Props = {
    recentRecord: RecentRecord;
    maxEffort: number;
    onClickDetail: (recordId: string) => void;
};

function RecentRecordCard({ recentRecord, maxEffort, onClickDetail }: Props) {
    return (
        <div className="flex items-start justify-between gap-4">
            {/* 左側: 内容 */}
            <div className="min-w-0 flex-1 space-y-3">
                <div className="text-lg font-semibold text-gray-900">{recentRecord.createdAt}</div>

                <EffortStarsDisplay effort={recentRecord.effort} maxEffort={maxEffort} />

                {/* タグ */}
                <div className="flex flex-wrap gap-2">
                    {recentRecord.tags.map((t) => (
                        <TagChip key={t} tagName={t} />
                    ))}
                </div>

                <p className="text-sm leading-6 text-gray-800">{recentRecord.content}</p>
            </div>

            {/* 右側: 詳細ボタン */}
            <button
                type="button"
                onClick={() => onClickDetail(recentRecord.id)}
                className={`${primaryBtnStyle} rounded-lg`}
            >
                詳細
            </button>
        </div>
    );
}

export default RecentRecordCard;