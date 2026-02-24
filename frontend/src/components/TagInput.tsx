type Props = {
    inputTag: string;
    onChange: (v: string) => void;
    onAddTag: () => void;
}

function TagInput({ inputTag, onChange, onAddTag }: Props) {
    return (
        <div className="flex">
            <input
                type="text"
                value={inputTag}
                placeholder="タグを入力..."
                aria-label="タグを入力"
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                    onChange(e.target.value);
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
    );
}

export default TagInput;