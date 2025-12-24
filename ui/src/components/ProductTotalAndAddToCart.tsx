"use client";

import { useCartStore } from "@/store/cart";
import { Money } from "@/types/product";
import clsx from "clsx";
import { useState } from "react";
import ButtonCartAndStock from "./ButtonCartAndStock";

interface ProductTotalAndAddToCartProps {
    stock: number;
    price: Money;
    productID: string;
}

function ProductTotalAndAddToCartProps({ stock, price, productID }: ProductTotalAndAddToCartProps) {
    const addCartItem = useCartStore((state) => state.addCartItem)
    const addQuantityInCart = useCartStore((state) => state.addQuantityInCart)
    const currentQuantityInCart = useCartStore((state) => state.Items.find((item) => item.ProductID == productID)?.Quantity || 0)
    const isProductInCart = useCartStore((state) => state.Items.findIndex((item) => item.ProductID == productID) !== -1)
    const isLoading = useCartStore((state) => state.isLoading)
    const [total, setTotal] = useState(1);
    const setErrorMessage = useCartStore((state) => state.setErrorMessage)
    const errorMessage = useCartStore((state) => state.errorMessage)


    function submitForm(e: React.FormEvent) {
        e.preventDefault()

        if (isProductInCart === true) {
            if (total + currentQuantityInCart > stock) {
                setErrorMessage("quantity cannot exceed stock")
                setTimeout(() => {
                    setTotal(1)
                }, 1000)
                return
            }

            addQuantityInCart(productID, total)
        } else {
            addCartItem(productID, total)
        }

        setTotal(1)
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
                    <ButtonCartAndStock disabled={isLoading || 1 === total} action={handleDecTotal} >
                        -
                    </ButtonCartAndStock>
                    <input type="number"
                        disabled={isLoading || stock === 0}
                        value={total}
                        className="text-center bg-white text-2xl m-2 text-black max-w-20 h-10 border-2 rounded border-white"
                        onChange={handleTotal} />
                    {/* <ButtonCartAndStock disabled={isLoading || stock === total} action={handleAddTotal} > */}
                    <ButtonCartAndStock disabled={isLoading || total + currentQuantityInCart > stock - 1 || stock === total} action={handleAddTotal} >
                        +
                    </ButtonCartAndStock>
                </div>
                <div className="my-5 flex items-center justify-between">
                    <div>subtotal: </div>
                    <div>{idFormatter.format(total * price.units)}</div>
                </div>

                {errorMessage && <p className="text-red-500 p-2 text-center border rounded-xl bg-slate-200">{errorMessage}</p>}

                <button type="submit" disabled={isLoading || stock === 0}
                    className={
                        clsx("w-full mt-2 border rounded-xl border-gray-300 bg-gray-200 text-gray-700 h-10",
                            isLoading || stock === 0 ? "cursor-not-allowed bg-gray-400 text-gray-200" : "cursor-pointer")
                    }
                >
                    {isLoading ? "loading..." : "Add To Cart"}
                </button>
            </form>
        </div>
    )
}

export default ProductTotalAndAddToCartProps     