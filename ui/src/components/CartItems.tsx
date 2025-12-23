"use client"

import { useCartStore } from "@/store/cart"
import { GetDetailsProducts } from "@/types/product"
import { DataResponse } from "@/types/response"
import Image from "next/image"
import { useEffect, useState } from "react"
import CartProductAndSubTotal from "./CartProductAndSubTotal"

function CartItems() {

    const cartItems = useCartStore((state) => state.Items)
    const [carts, setCarts] = useState<GetDetailsProducts>()
    const totalPrice = useCartStore((state) => state.totalPrice)
    const calculateTotalPrice = useCartStore((state) => state.calculateTotalPrice)

    useEffect(() => {
        if (cartItems.length === 0) {
            setCarts(undefined)
            return
        }

        const params = new URLSearchParams()
        params.append("with_stock", "true")
        cartItems.forEach((item) => {
            params.append("id", String(item.ProductID))
        })

        async function getData() {
            try {
                const res = await fetch(`http://localhost:8080/product?${params.toString()}`)

                var data: DataResponse = await res.json()
                setCarts(data.data as GetDetailsProducts)
                calculateTotalPrice(data.data.products)
            } catch (err) {
                console.error("failed to fetch products", err)
            }
        }
        getData()
    }, [cartItems])


    const idFormatter = Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
    })

    return (
        <div>
            <div className="flex items-center justify-between my-2 text-xl">
                <div>
                    total cart: {cartItems.length}
                </div>
                <div>
                    total price: {idFormatter.format(totalPrice)}
                </div>
            </div>
            <div>
                {carts?.products.map((product, index) => (
                    <div key={index} className="w-full max-w-5xl my-2 py-2">
                        <div className="grid grid-cols-5 gap-4 w-full">
                            <div className="col-span-2">
                                <div className="relative aspect-square w-full overflow-hidden rounded-lg bg-white">
                                    <Image
                                        src={product.image_url}
                                        alt={product.name}
                                        fill
                                        priority
                                        className="object-contain p-2"
                                        sizes="20vw"
                                    />
                                </div>
                            </div>
                            <div className="col-span-3">
                                <h1 className="text-xl font-bold mb-4">{product.name}</h1>
                                <CartProductAndSubTotal
                                    stock={product.stock || 0}
                                    price={product.price}
                                    productID={product.id}
                                />
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div >
    )
}

export default CartItems
