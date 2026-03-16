export const baseBtnStyle = "h-10 border px-4 whitespace-nowrap text-sm font-semibold focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400";
export const primaryBtnStyle = `
${baseBtnStyle}
border-sky-300 bg-sky-500 text-white
enabled:hover:bg-sky-400
disabled:border-gray-300 disabled:bg-gray-200 disabled:text-gray-500
disabled:cursor-not-allowed
`;
export const secondaryBtnStyle = `${baseBtnStyle} border-gray-300 bg-white text-gray-800 hover:bg-gray-50`;

export const starStyle = "text-3xl leading-none";

export const borderStyle = "rounded-xl border border-gray-300 p-6";

export const errorTextStyle = "text-sm text-red-600";