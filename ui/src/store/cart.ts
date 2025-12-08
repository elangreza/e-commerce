import { create } from "zustand";
import { persist } from "zustand/middleware";
import { immer } from 'zustand/middleware/immer';

type Cart = {
    UserID: string;
    Items: CartItem[];
}

export type CartItem = {
    ProductID: string;
    Quantity: number;
}

type Actions = {
    addCartItem: (item: CartItem) => void
    removeCartItem: (productID: CartItem['ProductID']) => void
    updateUserID: (s: string) => void
}

export const useCartStore = create<Cart & Actions>()(
    persist(
        immer((set,) => ({
            UserID: "",
            Items: [],
            addCartItem: (item: CartItem) => set((state) => {
                const fItem = state.Items.find((i) => i.ProductID === item.ProductID);
                if (fItem) {
                    fItem.Quantity = item.Quantity;
                    return
                }

                state.Items.push(item)
            }),
            removeCartItem: (productID: string) => set((state) => {
                state.Items = state.Items.filter((item) => item.ProductID !== productID)
            }),
            updateUserID: (s: string) => set((state) => {
                state.UserID = s
            }),

        })), {
        name: "cart",
    })
)