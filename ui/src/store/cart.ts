import { Product } from "@/types/product";
import { create } from "zustand";
import { persist } from "zustand/middleware";
import { immer } from 'zustand/middleware/immer';

type Cart = {
    UserID: string;
    Items: CartItem[];
    totalPrice: number;
}

export type CartItem = {
    ProductID: string;
    Quantity: number;
}

type Actions = {
    addCartItem: (productID: CartItem['ProductID'], quantity: number) => void
    removeCartItem: (productID: CartItem['ProductID']) => void
    updateUserID: (s: string) => void
    calculateTotalPrice: (products: Product[]) => void
}

export const useCartStore = create<Cart & Actions>()(
    persist(
        immer((set, get) => ({
            UserID: "",
            Items: [],
            totalPrice: 0,
            addCartItem: (productID: CartItem['ProductID'], quantity: number) => {
                if (quantity === 0) {
                    return
                }

                set((state) => {
                    var item = state.Items.find((item) => item.ProductID === productID)
                    if (item) {
                        item.Quantity = quantity;
                        return
                    }

                    state.Items.push({
                        ProductID: productID,
                        Quantity: quantity,
                    })
                })
            },
            removeCartItem: (productID: string) => set((state) => {
                state.Items = state.Items.filter((item) => item.ProductID !== productID);
            }),
            updateUserID: (s: string) => set((state) => {
                state.UserID = s
            }),
            calculateTotalPrice: (products: Product[]) => {
                var total = 0
                for (const item of get().Items) {
                    const product = products.find((p) => p.id === item.ProductID)
                    if (product) {
                        total += product.price.units * item.Quantity
                    }
                }
                set((state) => {
                    state.totalPrice = total
                })
            }
        })), {
        name: "cart",
    })
)