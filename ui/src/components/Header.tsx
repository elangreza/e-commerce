"use client";

import { useCartStore } from "@/store/cart";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";
import HeaderText from "./HeaderText";

function Header() {
    const [search, setSearch] = useState("")
    const searchParams = useSearchParams();
    const router = useRouter()
    const totalCartItems = useCartStore((state) => state.Items.reduce((total, item) => total + item.Quantity, 0))
    const isLoading = useCartStore((state) => state.isLoading)

    function handleSubmit(e: React.FormEvent) {
        e.preventDefault()

        const params = new URLSearchParams(searchParams.toString())
        if (search.trim()) {
            params.set("search", search)
            params.delete("page")
        } else {
            params.delete("search")
            params.delete("page")
        }

        router.push(`/?${params.toString()}`)
    }

    return (
        <div className="flex justify-between items-center w-full max-w-5xl m-2 p-2">
            <Link href="/">
                <HeaderText name="Toko saya" />
            </Link>
            <form onSubmit={handleSubmit}>
                <input
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    type="text"
                    placeholder="Search for products..."
                    className="border border-white rounded-l-lg px-4 py-2 bg-gray-500 focus:bg-gray-600 text-white focus:outline-none"
                />
                <button className="px-4 py-2 border border-white rounded-r-lg hover:bg-white hover:text-gray-500 cursor-pointer" type="submit">
                    Search
                </button>
            </form>
            <Link href="/cart">
                <HeaderText name={isLoading ? `My Cart (loading...)` : `My Cart (${totalCartItems})`} />
            </Link>
        </div>
    )
}


export default Header



