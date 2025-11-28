"use client";

import { useState } from "react";


interface ProductTotalAndAddToCartProps {
    stock: number;
}

function ProductTotalAndAddToCartProps({ stock }: ProductTotalAndAddToCartProps) {
    function submitForm(e: React.FormEvent) {
        e.preventDefault()

        // add product into cart
    }

    const [total, setTotal] = useState(stock);

    function handleTotal(e: React.ChangeEvent<HTMLInputElement>) {
        var val = Number(e.target.value) || 1;
        setTotal(val)
    }

    return (
        <form onSubmit={submitForm}>
            <input type="number"
                disabled={stock === 0}
                value={total}
                onChange={handleTotal} />
        </form>
    )
}

export default ProductTotalAndAddToCartProps     