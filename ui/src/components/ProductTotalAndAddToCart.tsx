"use client";

import { CartItem, useCartStore } from "@/store/cart";
import { Money } from "@/types/product";
import { useState } from "react";
import ButtonCartAndStock from "./ButtonCartAndStock";


interface ProductTotalAndAddToCartProps {
    stock: number;
    price: Money;
    productID: string;
}

function ProductTotalAndAddToCartProps({ stock, price, productID }: ProductTotalAndAddToCartProps) {

    const addCartItem = useCartStore((state) => state.addCartItem)
    const currentQuantityInCart = useCartStore((state) => state.Items.find((item) => item.ProductID == productID)?.Quantity || 1)
    const [total, setTotal] = useState(currentQuantityInCart);

    function submitForm(e: React.FormEvent) {
        e.preventDefault()

        var a: CartItem = {
            ProductID: productID,
            Quantity: total,
        }

        addCartItem(a)

        // TODO add product into cart
    }


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

    const idFormatter = Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: price.currency_code,
    })

    return (
        <div>
            <p className="text-center my-2">current stock: {stock}</p>
            <form onSubmit={submitForm} className="my-2">
                <div className="flex items-center justify-center">
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
                </div>
                <div className="my-5 flex items-center justify-between">
                    <div>subtotal: </div>
                    <div>{idFormatter.format(total * price.units)}</div>
                </div>

                <button type="submit" className="w-full mt-2 border rounded-xl cursor-pointer border-gray-300 bg-gray-200 text-gray-700 h-10">
                    Add To Cart
                </button>
            </form>
        </div>
    )
}

export default ProductTotalAndAddToCartProps     