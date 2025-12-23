"use client";

interface ButtonCartAndStockProps {
    action: (event: React.MouseEvent<HTMLButtonElement>) => void;
    children: React.ReactNode;
    disabled: boolean;
}

function ButtonCartAndStock({
    action, children, disabled = false
}: ButtonCartAndStockProps) {
    return (
        <button
            className={
                "w-auto px-4 h-10 bg-white border-2 rounded border-white text-black disabled:bg-gray-300 disabled:border-gray-300 disabled:cursor-not-allowed"
            }
            onClick={action}
            disabled={disabled}
            type="button"
        >
            {children}
        </button>
    )
}

export default ButtonCartAndStock   