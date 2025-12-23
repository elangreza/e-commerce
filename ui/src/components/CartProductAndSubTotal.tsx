"use client";

import { useCartStore } from "@/store/cart";
import { Money } from "@/types/product";
import { useState } from "react";
import ButtonCartAndStock from "./ButtonCartAndStock";

interface CartProductAndSubTotalProps {
    stock: number;
    price: Money;
    productID: string;
}

function CartProductAndSubTotal({ stock, price, productID }: CartProductAndSubTotalProps) {
    const removeCartItem = useCartStore((state) => state.removeCartItem)
    const addCartItem = useCartStore((state) => state.addCartItem)
    const currentQuantityInCart = useCartStore((state) => state.Items.find((item) => item.ProductID == productID)?.Quantity || 1)
    const [total, setTotal] = useState(currentQuantityInCart);

    function handleTotal(e: React.ChangeEvent<HTMLInputElement>) {
        let val = Number(e.target.value) || 1;
        if (val < 1) {
            val = 1;
        } else if (val > stock) {
            val = stock;
        }

        setTotal(val)
        addCartItem(productID, val)
    }

    function handleDecTotal() {
        if (total === 1) {
            return
        }

        const newTotal = total - 1
        setTotal(newTotal)
        addCartItem(productID, newTotal)
    }

    function handleAddTotal() {
        if (total === stock) {
            return
        }

        const newTotal = total + 1
        setTotal(newTotal)
        addCartItem(productID, newTotal)
    }

    function removeProductInCart() {
        removeCartItem(productID)
    }

    const idFormatter = Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: price.currency_code,
    })

    return (
        <div>
            <p className="text-center my-2">current stock: {stock}</p>
            <div className="my-2">
                <div className="flex items-center justify-between ">
                    <div className="gap-2 flex items-center">
                        <ButtonCartAndStock disabled={1 === total} action={handleDecTotal} >
                            -
                        </ButtonCartAndStock>
                        <input type="number"
                            disabled={stock === 0}
                            value={total}
                            className="text-center bg-white text-2xl text-black max-w-20 h-10 border-2 rounded border-white"
                            onChange={handleTotal} />
                        <ButtonCartAndStock disabled={stock === total} action={handleAddTotal} >
                            +
                        </ButtonCartAndStock>
                    </div>
                    <ButtonCartAndStock disabled={stock === total} action={removeProductInCart} >
                        remove from cart
                    </ButtonCartAndStock>
                </div>
                <div className="my-5 flex items-center justify-between">
                    <div>subtotal: </div>
                    <div>{idFormatter.format(total * price.units)}</div>
                </div>
            </div>
        </div>
    )
}

export default CartProductAndSubTotal     