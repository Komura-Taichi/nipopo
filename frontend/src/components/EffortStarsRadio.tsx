import { starStyle } from "../styles";

type Props = {
    effort: number;
    maxEffort: number;
    onChange: (clickStar: number) => void;
}

function EffortStarsRadio({ effort, maxEffort, onChange }: Props) {
    return (
        <div className="flex items-center gap-2">
            {Array.from({ length: maxEffort }, (_, i) => i + 1).map((n) => (
                <label
                    key={`er_{n}`}
                    className="cursor-pointer select-none"
                >
                    <input
                        type="radio"
                        name="effort"
                        className="sr-only"
                        aria-label={`${n} / ${maxEffort}`}
                        checked={effort === n}
                        onChange={() => onChange(n)}
                    />
                    <span className={`${starStyle} ${n <= effort ? "text-yellow-400" : "text-gray-300"}`}>
                        {n <= effort ? "★" : "☆"}
                    </span>
                </label>
            ))}
        </div>
    );
}

export default EffortStarsRadio;