"use client";

interface ButtonCartAndStockProps {
    action: (event: React.MouseEvent<HTMLButtonElement>) => void;
    children: React.ReactNode;
    disabled: boolean;
}

function ButtonCartAndStock({
    action, children, disabled = false
}: ButtonCartAndStockProps) {
    console.log(disabled)
    return (
        <button
            className={
                "max-w-10 w-full h-10 bg-white border-2 rounded border-white text-black disabled:bg-gray-300 disabled:border-gray-300 disabled:cursor-not-allowed"
            }
            onClick={action}
            disabled={disabled}
        >
            {children}
        </button>
    )
}

export default ButtonCartAndStock   