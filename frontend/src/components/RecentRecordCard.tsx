import type { RecentRecord } from "../types";
import EffortStarsDisplay from "./EffortStarsDisplay";
import TagChip from "./TagChip";

type Props = {
    recentRecord: RecentRecord;
    maxEffort: number;
};

function RecentRecordCard({ recentRecord, maxEffort }: Props) {
    return (
        <div key={recentRecord.id} className="space-y-3">
            <div className="text-lg font-semibold text-gray-900">{recentRecord.createdAt}</div>

            <EffortStarsDisplay effort={recentRecord.effort} maxEffort={maxEffort} />

            {/* タグ */}
            <div className="flex flex-wrap gap-2">
                {recentRecord.tags.map((t) => (
                    <TagChip tagName={t} />
                ))}
            </div>

            <p className="text-sm leading-6 text-gray-800">{recentRecord.content}</p>
        </div>
    );
}

export default RecentRecordCard;