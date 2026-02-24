import { starStyle } from "../styles";

type Props = {
    effort: number;
    maxEffort: number;
};

function EffortStarsDisplay({ effort, maxEffort }: Props) {
    return (
        <div className="flex gap-2" aria-label="頑張り度（星）">
            {Array.from({ length: maxEffort }, (_, i) => i + 1).map((n) => (
                <span key={`ed_${n}`} className={`${starStyle} ${n <= effort ? "text-yellow-400" : "text-gray-300"}`}>
                    {n <= effort ? "★" : "☆"}
                </span>
            ))}
        </div>
    )
}

export default EffortStarsDisplay;