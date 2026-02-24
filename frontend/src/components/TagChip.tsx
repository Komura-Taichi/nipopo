type Props = {
    tagName: string;
    onRemove?: () => void;
}

function TagChip({ tagName, onRemove }: Props) {
    return (
        <span
            className="inline-flex items-center gap-2 rounded-full border border-gray-300 bg-white px-3 py-1 text-sm text-gray-800"
        >
            {tagName}
            {onRemove && (
                <button
                    type="button"
                    aria-label={`タグ 「${tagName}」 を削除`}
                    onClick={onRemove}
                    className="rounded-full px-1 text-gray-500 hover:bg-gray-100"
                >
                    ×
                </button>
            )}
        </span>
    );
}

export default TagChip;