import CartItems from "@/components/CartItems";



export default async function CartPage() {

    return (
        <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-gray-600">
            <main className="min-h-screen w-full max-w-5xl flex-col items-center justify-between p-10 bg-white dark:bg-gray-500 sm:items-start">
                {/* <div className="grid grid-rows-1 grid-cols-2 gap-12"> */}
                <CartItems />
                {/* </div> */}
            </main>
        </div>
    );
}
