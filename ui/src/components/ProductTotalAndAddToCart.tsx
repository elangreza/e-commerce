"use client";

import { useState } from "react";
import ButtonCartAndStock from "./ButtonCartAndStock";


interface ProductTotalAndAddToCartProps {
    stock: number;
}

function ProductTotalAndAddToCartProps({ stock }: ProductTotalAndAddToCartProps) {
    function submitForm(e: React.FormEvent) {
        e.preventDefault()

        // TODO add product into cart
    }

    const [total, setTotal] = useState(1);

    function handleTotal(e: React.ChangeEvent<HTMLInputElement>) {
        var val = Number(e.target.value) || 1;
        if (val < 1) {
            setTotal(1)
        }
        if (val > stock) {
            setTotal(stock)
        }
    }

    function handleDecTotal() {
        if (total === 1) {
            return
        }

        setTotal(old => old - 1)
    }


    function handleAddTotal() {
        if (total === stock) {
            return
        }

        setTotal(old => old + 1)
    }



    return (
        <div>
            <p className="text-center">current stock: {stock}</p>

            <form onSubmit={submitForm}
                className="flex items-center justify-center"
            >
                <ButtonCartAndStock disabled={1 === total} action={handleDecTotal} >
                    -
                </ButtonCartAndStock>
                <input type="number"
                    disabled={stock === 0}
                    value={total}
                    className="text-center bg-white text-2xl m-2 text-black max-w-20 h-10 border-2 rounded border-white"
                    onChange={handleTotal} />
                <ButtonCartAndStock disabled={stock === total} action={handleAddTotal} >
                    +
                </ButtonCartAndStock>
            </form>
        </div>
    )
}

export default ProductTotalAndAddToCartProps     