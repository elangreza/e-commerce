"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";
import HeaderText from "./HeaderText";

function Header() {
    const [search, setSearch] = useState("")
    const searchParams = useSearchParams();
    const router = useRouter()

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
        <div className="flex justify-between items-center mb-10">
            <HeaderText name="Toko saya" />
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
            <HeaderText name="My Cart" />
        </div>
    )
}


export default Header



