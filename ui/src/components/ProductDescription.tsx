"use client";

import clsx from "clsx";
import { useState } from "react";

interface ProductDescriptionProps {
    text: string;
}

function ProductDescription({ text }: ProductDescriptionProps) {
    const [isOpen, setIsOpen] = useState(false)
    console.log(text.length)
    return (
        <div>
            <p className={
                clsx("mt-2 mb-1 text-gray-200", isOpen ? "" : "line-clamp-2")
            }>{text}</p>
            {
                text.length > 100 &&
                <button
                    onClick={() => setIsOpen(old => !old)}
                    className="text-green-300"
                >
                    {isOpen ? "see less" : "see more"}
                </button>
            }
        </div>
    )
}

export default ProductDescription